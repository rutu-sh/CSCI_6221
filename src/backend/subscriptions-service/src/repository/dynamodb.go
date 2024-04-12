package repository

import (
	"fmt"
	"os"
	"subHandler/src/config"
	"subHandler/src/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func initialize(service string) models.DynamoAttr {
	/*
		Used to initialize the table attributes and sdk clients
		Params: None
		Return: None
	*/
	awsRegion := config.AWS_REGION
	var dynamodbTable string
	if service == "payments" {
		dynamodbTable = config.PAYMENTS_DYNAMODB_TABLE
	} else {
		dynamodbTable = config.SUBSCRIPTIONS_DYNAMODB_TABLE
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		fmt.Println("Error creating sess", err)
		os.Exit(0)
	}

	dynamoClient := dynamodb.New(sess)
	return models.DynamoAttr{
		DynamoCli: dynamoClient,
		AwsRegion: awsRegion,
		TableName: dynamodbTable,
	}
}
