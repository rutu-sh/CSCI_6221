// Package subscriptions
/*
Used to update the fetch a specific item from DynamoDB or
fetch all subscriptions for a user from DynamoDB
*/
package subscriptions

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const workerCount = 10

type GetItem struct {
	UUID                 string `json:"uuid"`
	VendorName           string `json:"vendor_name"`
	VendorUrl            string `json:"vendor_url"`
	SubscriptionDuration string `json:"duration"`
	RemindTime           string `json:"remind_time"`
}

type GetResponse struct {
	UUID                 string `json:"uuid"`
	VendorName           string `json:"vendor_name"`
	VendorUrl            string `json:"vendor_url"`
	SubscriptionDuration string `json:"duration"`
	Status               int    `json:"status"`
	Message              string `json:"message"`
}

type SubList struct {
	GetItem
	UserName string `json:"username"`
}

type StatusResponse struct {
	Message string `json:"messsage"`
	Status  string `json:"status"`
}

type SubResponse struct {
	StatusCode    int
	Subscriptions []GetResponse `json:"subscriptions"`
}

type Result struct {
	subscription GetResponse
	err          error
}

type Job struct {
	item map[string]*dynamodb.AttributeValue
}

func worker(jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		item := GetItem{}
		err := dynamodbattribute.UnmarshalMap(job.item, &item)
		if err != nil {
			results <- Result{err: err}
			continue
		}
		results <- Result{
			subscription: GetResponse{
				UUID:                 item.UUID,
				VendorName:           item.VendorName,
				VendorUrl:            item.VendorUrl,
				SubscriptionDuration: item.SubscriptionDuration,
				Status:               200,
			},
		}
	}
}

func GetSubscription(dynamoClient *dynamodb.DynamoDB, tableName, uuid, username string) GetResponse {
	/*
		Retrieves the given Item in the DynamoDB table based on the UUID
		and the username provided.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName string
				uuid string
				username string
		Return: GetResponse
	*/

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UUID": {
				S: aws.String(uuid),
			},
			"UserName": {
				S: aws.String(username),
			},
		},
	}

	result, derr := dynamoClient.GetItem(input)
	log.Println(result)
	if derr != nil {
		if aerr, ok := derr.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Printf("Throughput exceeded for table %s please try again after sometime", tableName)
				return GetResponse{
					Message: "Throughput exceeded for table",
					Status:  400,
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Printf("Table %s not found, or item not found please check tablename/item", tableName)
				return GetResponse{
					Message: "Table not found, please check the table name",
					Status:  400,
				}
			}
		} else {
			log.Printf("Got error calling GetItem: %s", derr)
			return GetResponse{
				Message: "Unexpected error occurred",
				Status:  400,
			}
		}
	}

	item := GetItem{}
	err := dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Printf("Failed to unmarshal Record, %v", err)
		return GetResponse{Status: 500}
	}

	log.Printf("Successfully retreived Item from table %s", tableName)
	return GetResponse{
		UUID:                 item.UUID,
		VendorName:           item.VendorName,
		VendorUrl:            item.VendorUrl,
		SubscriptionDuration: item.SubscriptionDuration,
		Status:               200,
	}
}

func GetSubscriptions(dynamoClient *dynamodb.DynamoDB, tableName, userName string) SubResponse {
	/*
		Retrieves all Items in the DynamoDB table that matches the
		userName.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName string
				username string
		Return: SubResponse
	*/

	filter := expression.Name("UserName").Equal(expression.Value(userName))
	expr, _ := expression.NewBuilder().WithFilter(filter).Build()

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		TableName:                 aws.String(tableName),
	}

	result, err := dynamoClient.Scan(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Printf("Throughput exceeded for table %s please try again after sometime", tableName)
				return SubResponse{
					StatusCode: 400,
					Subscriptions: []GetResponse{
						{
							Message: "Throughput exceeded for table",
						},
					},
				}
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Printf("Table %s not found, or item not found please check tablename/item", tableName)
				return SubResponse{
					StatusCode: 400,
					Subscriptions: []GetResponse{
						{
							Message: "Table not found, please check the table name",
						},
					},
				}
			}
		} else {
			log.Printf("Got error calling DeleteItem: %s", err)
			return SubResponse{
				StatusCode: 400,
				Subscriptions: []GetResponse{
					{
						Message: "Unexpected error occurred",
					},
				},
			}
		}
	}
	/*
		Using the worker pool pattern we process all the subscriptions for each user
		Two channels are created
			1. Job - to send work to the worker goroutines
			2. Results - to receive the results from the worker goroutines
		10 worker goroutines are created. Each worker goroutine receives Job instances
		fom the jobs channel, processes them, and stores the results back to Results
		channel.
		All subscriptions queried from the Table are sent to the Jobs channel. Then
		10 workers process them in parallel.
	*/
	jobs := make(chan Job, len(result.Items))
	results := make(chan Result, len(result.Items))
	subscriptions := []GetResponse{}

	for w := 0; w < workerCount; w++ {
		go worker(jobs, results)
	}

	for _, i := range result.Items {
		jobs <- Job{item: i}
	}
	close(jobs)

	for range result.Items {
		res := <-results
		if res.err != nil {
			log.Printf("Failed to unmarshal Record, %v", res.err)
			return SubResponse{
				StatusCode: 400,
				Subscriptions: []GetResponse{
					{
						Message: "Failed to unmarshal Record",
					},
				},
			}
		}
		subscriptions = append(subscriptions, res.subscription)
	}

	finalResponse := SubResponse{
		StatusCode:    200,
		Subscriptions: subscriptions,
	}

	return finalResponse
}
