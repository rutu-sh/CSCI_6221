package dynamoSub

import (
	"Notifier/src/sns_notifier"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
	"time"
)

type SubscriptionsToAlert struct {
	UserName   string
	VendorName string
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
	/*
		For each of the subscription send out
		an alert to notify them about the subscription.
		Params: dynamoCli *dynamodb.DynamoDB,
				snsCli *sns.SNS,
				snsArn string
		Returned: None
	*/

	for _, item := range GetAllExpiringSubscriptions(dynamoCli) {

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
			if aerr, ok := derr.(awserr.Error); ok {
				switch aerr.Code() {
				case dynamodb.ErrCodeProvisionedThroughputExceededException:
					log.Printf("Throughput exceeded for table %s please try again after sometime", "users")
					os.Exit(1)
				case dynamodb.ErrCodeResourceNotFoundException:
					log.Printf("Table %s not found, or item not found please check tablename/item", "users")
					os.Exit(1)
				}
			} else {
				log.Printf("Got error calling GetItem: %s", derr)
				os.Exit(1)
			}
		}
		if result.Item == nil {
			fmt.Printf("%s does not have any subscriptions \n", item.UserName)
			continue
		}
		fmt.Println(result)
		email := *result.Item["Email"].S
		emailValues := sns_notifier.SnsMessageFormat(item.UserName, item.VendorName)
		sns_notifier.PublishMessage(snsCli, snsArn, emailValues.Message, emailValues.Body, email)
	}
}
