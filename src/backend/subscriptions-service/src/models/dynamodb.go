package models

import "github.com/aws/aws-sdk-go/service/dynamodb"

type DynamoAttr struct {
	DynamoCli *dynamodb.DynamoDB
	AwsRegion string
	TableName string
}
