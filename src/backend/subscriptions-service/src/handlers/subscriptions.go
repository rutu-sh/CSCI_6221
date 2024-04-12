package handlers

import (
	"context"
	"encoding/json"
	"subHandler/src/models"
	"subHandler/src/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func SubscriptionsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	/*
		Handles the routing of the specific API call to the respective function.
		Also, handles the response and error handling for the API calls.
		Params: ctx context.Context
				request events.APIGatewayProxyRequest
		Returns: events.APIGatewayProxyResponse
				 error
	*/
	reqMethod := request.HTTPMethod
	if reqMethod == "POST" {
		reqBody := request.Body
		if reqBody == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		var sub models.SubscriptionCreateInput
		err := json.Unmarshal([]byte(reqBody), &sub)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		res, err := service.AddSubscription(sub)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		resBody, err := json.Marshal(res)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 201,
			Body:       string(resBody),
		}, nil
	}
	if reqMethod == "GET" {
		userName := request.QueryStringParameters["username"]
		if userName == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		res, err := service.GetUserSubscriptions(userName)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		resBody, err := json.Marshal(res)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(resBody),
		}, nil
	}
	if reqMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}

func SubscriptionByIDHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	/*
		Handles the routing of the specific API call to the respective function.
		Also, handles the response and error handling for the API calls.
		Params: ctx context.Context
				request events.APIGatewayProxyRequest
		Returns: events.APIGatewayProxyResponse
				 error
	*/
	reqMethod := request.HTTPMethod
	if reqMethod == "GET" {
		subID := request.PathParameters["subscription-id"]
		userName := request.QueryStringParameters["username"]
		log.Info().Str("subID", subID).Str("userName", userName).Msg("Received request with parameters")
		if subID == "" || userName == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		res, err := service.GetSubscription(subID, userName)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		resBody, err := json.Marshal(res)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}

		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(resBody),
		}, nil
	}

	if reqMethod == "DELETE" {
		subID := request.PathParameters["subscription-id"]
		userName := request.QueryStringParameters["username"]
		if subID == "" || userName == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		err := service.DeleteSubscription(subID, userName)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Body:       `{"message": "Subscription deleted"}`,
		}, nil
	}

	if reqMethod == "PATCH" {
		subID := request.PathParameters["subscription-id"]
		userName := request.QueryStringParameters["username"]
		if subID == "" || userName == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		reqBody := request.Body
		if reqBody == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		var sub models.SubscriptionUpdate
		err := json.Unmarshal([]byte(reqBody), &sub)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		res, err := service.UpdateSubscription(subID, userName, sub)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		resBody, err := json.Marshal(res)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(resBody),
		}, nil
	}

	if reqMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}
