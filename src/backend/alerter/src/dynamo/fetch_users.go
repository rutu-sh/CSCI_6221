package dynamoSub

import (
	"Notifier/src/sns_notifier"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
)

const workerCount = 10

type SubscriptionsToAlert struct {
	UserName   string
	VendorName string
}

type Job struct {
	Subscription SubscriptionsToAlert
}

type Result struct {
	Subscription SubscriptionsToAlert
	Error        error
}

func worker(dynamoCli *dynamodb.DynamoDB, snsCli *sns.SNS, snsArn string, jobs <-chan Job, results chan<- Result) {
	for job := range jobs {
		item := job.Subscription

		input := &dynamodb.GetItemInput{
			TableName: aws.String("users"),
			Key: map[string]*dynamodb.AttributeValue{
				"UserName": {
					S: aws.String(item.UserName),
				},
			},
		}

		log.Printf("Getting email for username %s \n", item.UserName)
		result, derr := dynamoCli.GetItem(input)
		if derr != nil {
			results <- Result{Subscription: item, Error: derr}
			continue
		}
		if result.Item == nil {
			fmt.Printf("%s does not have any subscriptions \n", item.UserName)
			continue
		}
		fmt.Println(result)
		email := *result.Item["Email"].S
		emailValues := sns_notifier.SnsMessageFormat(item.UserName, item.VendorName)
		sns_notifier.PublishMessage(snsCli, snsArn, emailValues.Message, emailValues.Body, email)

		results <- Result{Subscription: item}
	}
}

func GetAllExpiringSubscriptions(dynamoCli *dynamodb.DynamoDB) []SubscriptionsToAlert {
	/*
		Gets all the subscriptions that are only one day
		from getting renewed.
		Params: dynamoCli *dynamodb.DynamoDB
		Returned: []SubscriptionsToAlert
	*/

	nextDay := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	scanExpr := &dynamodb.ScanInput{
		TableName:        aws.String("subscriptions"),
		FilterExpression: aws.String("RemindTime = :rt"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":rt": {S: aws.String(nextDay)},
		},
	}
	log.Println("Scanning the dynamoDB table to get expiring subscriptions")
	result, err := dynamoCli.Scan(scanExpr)
	if err != nil {
		fmt.Println("Error scanning table:", err)
		return nil
	}

	var subscriptions []SubscriptionsToAlert

	for _, item := range result.Items {
		subscription := SubscriptionsToAlert{
			UserName:   *item["UserName"].S,
			VendorName: *item["VendorName"].S,
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions
}

func SendAlert(dynamoCli *dynamodb.DynamoDB, snsCli *sns.SNS, snsArn string) {
	subscriptions := GetAllExpiringSubscriptions(dynamoCli)

	jobs := make(chan Job, len(subscriptions))
	results := make(chan Result, len(subscriptions))

	for w := 1; w <= workerCount; w++ {
		go worker(dynamoCli, snsCli, snsArn, jobs, results)
	}

	for _, subscription := range subscriptions {
		jobs <- Job{Subscription: subscription}
	}
	close(jobs)

	for range subscriptions {
		result := <-results
		if result.Error != nil {
			log.Printf("Error processing subscription for user %s: %v", result.Subscription.UserName, result.Error)
		}
	}
}
