/*
	Package to alert users when their

subscription is ending its near so that they can
either cancel or renew it
*/
package main

import (
	dynamoSub "Notifier/src/dynamo"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"os"
)

func main() {
	/*
		Call the alerting service
	*/
	region := os.Getenv("region")
	snsArn := os.Getenv("sns_arn")

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(region),
		},
	}))

	dynamoCli := dynamodb.New(sess)
	snsCli := sns.New(sess)

	dynamoSub.SendAlert(dynamoCli, snsCli, snsArn)

}
