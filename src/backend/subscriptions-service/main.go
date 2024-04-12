package main

import (
	"context"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog/log"

	"subHandler/src/handlers"
)

func getCORSHeaders() map[string]string {
	/*
		Returns the CORS headers for the API response
		Params: None
		Return: map[string]string
	*/
	return map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Methods":     "GET, POST, PUT, OPTIONS, DELETE",
		"Access-Control-Allow-Headers":     "Content-Type, Authorization",
		"Access-Control-Allow-Credentials": "true",
	}
}

type HandlerFunc func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func getHandlerFunc(path string) (HandlerFunc, error) {
	subscriptionsRegex, err := regexp.Compile(`^\/v2\/subscriptions$`)
	if err != nil {
		return nil, err
	}
	if subscriptionsRegex.MatchString(path) {
		return handlers.SubscriptionsHandler, nil
	}

	subscriptionByIdRegex, err := regexp.Compile(`^\/v2\/subscriptions\/[a-zA-Z0-9-]+$`)
	if err != nil {
		return nil, err
	}
	if subscriptionByIdRegex.MatchString(path) {
		return handlers.SubscriptionByIDHandler, nil
	}

	paymentsRegex, err := regexp.Compile(`^\/v2\/payments$`)
	if err != nil {
		return nil, err
	}
	if paymentsRegex.MatchString(path) {
		return handlers.PaymentsHandler, nil
	}
	return nil, nil
}

func pathHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Info().Str("path", request.Path).Msg("Received request")
	handler, err := getHandlerFunc(request.Path)
	headers := getCORSHeaders()
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error", Headers: headers}, err
	}
	if handler == nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Not Found", Headers: headers}, nil
	}
	response, err := handler(ctx, request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error", Headers: headers}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       response.Body,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(pathHandler)
}
