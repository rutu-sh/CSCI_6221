package handlers

import (
	"context"
	"encoding/json"
	"subHandler/src/models"
	"subHandler/src/service"

	"github.com/aws/aws-lambda-go/events"
)

func PaymentsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
		var pay models.PaymentCreateInput
		err := json.Unmarshal([]byte(reqBody), &pay)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		res, err := service.AddPayment(pay)
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
		subscriptionId := request.QueryStringParameters["subscription_id"]
		if subscriptionId == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		res, err := service.GetPayments(subscriptionId)
		if len(res) == 0 {
			return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Not Found"}, nil
		}
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
	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}

func PaymentByIDHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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
		paymentId := request.PathParameters["payment_id"]
		subscriptionId := request.QueryStringParameters["subscription_id"]
		if paymentId == "" || subscriptionId == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		res, err := service.GetPayment(subscriptionId, paymentId)
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
	if reqMethod == "PATCH" {
		paymentId := request.PathParameters["payment_id"]
		subscriptionId := request.QueryStringParameters["subscription_id"]
		if paymentId == "" || subscriptionId == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		reqBody := request.Body
		if reqBody == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		var pay models.PaymentUpdate
		err := json.Unmarshal([]byte(reqBody), &pay)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		res, err := service.UpdatePayment(subscriptionId, paymentId, pay)
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
		paymentId := request.PathParameters["payment_id"]
		subscriptionId := request.QueryStringParameters["subscription_id"]
		if paymentId == "" || subscriptionId == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		err := service.DeletePayment(subscriptionId, paymentId)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
		}
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
		}, nil
	}
	if reqMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}
