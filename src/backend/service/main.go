/*
Package for handling the requests sent to the API gateway and call
the correct functions to handle the request
TODO: Simply the initialize process by removing re-use of initialize function
*/
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"subHandler/src/handlers"
)

type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func callHandler(hfunc HandlerFunc, ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	/*
		Handles the routing of the specific API call to the respective function.
		Also, handles the response and error handling for the API calls.
		Params: hfunc A function of type HandlerFunc which is the handler function to be called.
				ctx context.Context
				request
		Returns: events.APIGatewayProxyResponse
				 error
	*/
	response, err := hfunc(ctx, request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
	}

	responseJSON, err := json.Marshal(response.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       string(responseJSON),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Methods":     "GET, POST, OPTIONS, DELETE",
			"Access-Control-Allow-Headers":     "Content-Type, Authorization",
			"Access-Control-Allow-Credentials": "true",
		},
	}, nil
}

func getHandlerFunc(path string) (HandlerFunc, error) {
	subscriptionRegex, err := regexp.Compile(`^\/subscriptions$`)
	if err != nil {
		return nil, err
	}
	subscriptionIDRegex, err := regexp.Compile(`^\/subscriptions\/([0-9a-zA-Z-]+)\/user\/([0-9a-zA-Z-]+)$`)
	if err != nil {
		return nil, err
	}
	subscriptionListRegex, err := regexp.Compile(`^\/subscriptions\/list\/([0-9a-zA-Z-]+)$`)
	if err != nil {
		return nil, err
	}
	updateSubscriptionIDRegex, err := regexp.Compile(`^\/subscriptions\/update\/([0-9a-zA-Z-]+)$`)
	if err != nil {
		return nil, err
	}
	switch true {
	case subscriptionListRegex.MatchString(path) || subscriptionRegex.MatchString(path):
		return handlers.HandlerSubscription, nil
	case subscriptionIDRegex.MatchString(path):
		return handlers.HandlerSubscriptionID, nil
	case updateSubscriptionIDRegex.MatchString(path):
		return handlers.HandlerUpdate, nil
	default:
		return nil, errors.New("no handler found")
	}
}

func handlerPath(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	/*
		Handles the routing of the specific API call to the respective function
		Params: ctx context.Context
				request events.APIGatewayProxyRequest
		Returns: events.APIGatewayProxyResponse
		         error
	*/
	path := request.Path
	var handlerFunc HandlerFunc
	fmt.Println("Entering Handler path: ", path)

	// get handler function based on the path
	handlerFunc, err := getHandlerFunc(path)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
	}

	res, err := callHandler(handlerFunc, ctx, request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
	}
	return res, nil
}

func main() {
	lambda.Start(handlerPath)
}
