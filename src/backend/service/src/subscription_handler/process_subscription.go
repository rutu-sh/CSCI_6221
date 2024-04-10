package subscriptions

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const workerCount = 10

// GetSubscriptionList retrieves the list of subscriptions for the given username
func GetSubscriptionList(dynamoClient *dynamodb.DynamoDB, tableName, username string) SubResponse {
	/*
		Retrieves the list of subscriptions for the given username
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				username
		Return: SubResponse
	*/
	var response SubResponse

	// Create the expression to get the list of subscriptions for the given username
	filt := expression.Name("username").Equal(expression.Value(username))
	proj := expression.NamesList(expression.Name("uuid"), expression.Name("username"), expression.Name("remind_time"), expression.Name("vendor_name"), expression.Name("vendor_url"), expression.Name("duration"))
	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	if err != nil {
		log.Println("Got error building expression:")
		log.Println(err.Error())
		response.StatusCode = 500
		response.Message = "Internal Server Error"
		return response
	}

	// Create the input for the query
	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the query to the table
	result, err := dynamoClient.Scan(input)
	if err != nil {
		log.Println("Got error calling Scan:")
		log.Println(err.Error())
		response.StatusCode = 500
		response.Message = "Internal Server Error"
		return response
	}

	// Parse the result and return the response
	response.StatusCode = 200
	response.Message = "Success"
	response.Subscriptions = parseSubscriptionList(result.Items)
	return response
}
