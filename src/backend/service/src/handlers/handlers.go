package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"

	subscriptions "subHandler/src/subscription_handler"
)

type dynamoAttr struct {
	dynamoCli *dynamodb.DynamoDB
	awsRegion string
	tableName string
}

type GetItem struct {
	UserName string
}

func initialize() dynamoAttr {
	/*
		Used to initialize the table attributes and sdk clients
		Params: None
		Return: None
	*/
	awsRegion := os.Getenv("aws_region") // Get from ENV vars of lambda
	dynamodbTable := os.Getenv("table_name")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		fmt.Println("Error creating sess", err)
		os.Exit(0)
	}

	dynamoClient := dynamodb.New(sess)
	return dynamoAttr{
		dynamoCli: dynamoClient,
		awsRegion: awsRegion,
		tableName: dynamodbTable,
	}
}

func HandlerSubscription(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	/*
		Handles the dynamodb Put and Get operation for the API gateway calls
		GET request returns all the subscriptions based on the username
		POST adds the given subscription to the table
		Params: ctx context.Context
				request events.APIGatewayProxyRequest
		Returns: events.APIGatewayProxyResponse
		         error
	*/

	httpMethod := request.HTTPMethod

	dynamoCli := initialize()
	if httpMethod == "GET" {

		userName := request.PathParameters["userName"]

		subscriptionsResponse := subscriptions.GetSubscriptions(dynamoCli.dynamoCli, dynamoCli.tableName, userName)
		responseBody, err := json.Marshal(subscriptionsResponse)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: subscriptionsResponse.StatusCode,
			Body:       string(responseBody),
		}, nil
	}
	if httpMethod == "POST" {
		var addItem subscriptions.AddItem
		err := json.Unmarshal([]byte(request.Body), &addItem)
		finalAddItem := subscriptions.AddItem{
			UUID:                 uuid.New().String(),
			UserName:             addItem.UserName,
			VendorName:           addItem.VendorName,
			VendorUrl:            addItem.VendorUrl,
			SubscriptionDuration: addItem.SubscriptionDuration,
			RemindTime:           subscriptions.CalcRemindTime(addItem.SubscriptionDuration), // need to use this to calc RemindTime
		}
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		subscriptionsPostResponse := subscriptions.AddItemToTable(dynamoCli.dynamoCli, dynamoCli.tableName, finalAddItem)
		responseBody, err := json.Marshal(subscriptionsPostResponse)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: subscriptionsPostResponse.StatusCode,
			Body:       string(responseBody),
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}

func HandlerSubscriptionID(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	/*
		Handles the dynamodb Get and Delete operation for the API gateway calls
		GET request returns information of a specific ID and username
		DELETE request is used to remove the specified Item from the table
		Params: ctx context.Context
				request events.APIGatewayProxyRequest
		Returns: events.APIGatewayProxyResponse
		         error
	*/
	httpMethod := request.HTTPMethod
	dynamoCli := initialize()
	if httpMethod == "GET" {
		subscriptionID := request.PathParameters["subscription_id"]
		userName := request.PathParameters["userName"]
		subscriptionResponse := subscriptions.GetSubscription(dynamoCli.dynamoCli, dynamoCli.tableName, subscriptionID, userName)
		responseBody, err := json.Marshal(subscriptionResponse)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: subscriptionResponse.Status,
			Body:       string(responseBody),
		}, nil
	}
	if httpMethod == "DELETE" {
		subscriptionID := request.PathParameters["subscription_id"]
		userName := request.PathParameters["userName"]
		deleteResponse := subscriptions.DeleteItemFromTable(dynamoCli.dynamoCli, dynamoCli.tableName, subscriptionID, userName)
		responseBody, err := json.Marshal(deleteResponse)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: deleteResponse.StatusCode,
			Body:       string(responseBody),
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}

func HandlerUpdate(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	/*
		Handles the dynamodb UPDATE operation for the API gateway calls
		POST updates the given subscription attributes of the table
		Params: ctx context.Context
				request events.APIGatewayProxyRequest
		Returns: events.APIGatewayProxyResponse
		         error
	*/
	httpMethod := request.HTTPMethod
	dynamoCli := initialize()
	if httpMethod == "POST" {
		var updateItem subscriptions.UpdateItem
		subscriptionID := request.PathParameters["subscription_id"]
		err := json.Unmarshal([]byte(request.Body), &updateItem)
		finalUpdateItem := subscriptions.UpdateItem{
			UserName:             updateItem.UserName,
			VendorName:           updateItem.VendorName,
			VendorUrl:            updateItem.VendorUrl,
			SubscriptionDuration: updateItem.SubscriptionDuration,
			RemindTime:           subscriptions.CalcRemindTime(updateItem.SubscriptionDuration),
		}
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		updateResponse := subscriptions.UpdateSubscription(dynamoCli.dynamoCli, dynamoCli.tableName, subscriptionID, finalUpdateItem)
		responseBody, err := json.Marshal(updateResponse)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: updateResponse.Status,
			Body:       string(responseBody),
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}
