package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"movers/src/cognito_auth"
	"movers/src/notifier"
	"os"
)

type cognitoAttr struct {
	cognitoCli *cognitoidentityprovider.CognitoIdentityProvider
	awsRegion  string
	userpoolId string
	clientId   string
	dynamoCli  *dynamodb.DynamoDB
	tableName  string
	snsCli     *sns.SNS
	topicArn   string
}

type Signup struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

type Signin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Coggen struct {
	Username string `json:"username"`
}

type Resetgen struct {
	Username string `json:"username"`
	Confcode string `json:"confcode"`
}

type ResetPass struct {
	Username string `json:"username"`
	Confcode string `json:"confcode"`
	Password string `json:"password"`
}

func initialize() cognitoAttr {
	awsRegion := os.Getenv("aws_region")
	userpoolId := os.Getenv("userpool_id")
	tableName := os.Getenv("table_name")
	clientId := os.Getenv("client_id")
	snsArn := os.Getenv("sns_topic_arn")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		fmt.Println("Error creating sess", err)
		os.Exit(0)
	}

	cognitoClient := cognitoidentityprovider.New(sess)
	dynamoClient := dynamodb.New(sess)
	snsClient := sns.New(sess)
	return cognitoAttr{
		cognitoCli: cognitoClient,
		clientId:   clientId,
		awsRegion:  awsRegion,
		userpoolId: userpoolId,
		dynamoCli:  dynamoClient,
		tableName:  tableName,
		snsCli:     snsClient,
		topicArn:   snsArn,
	}
}

func handlerSignUp(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var signup Signup
	err := json.Unmarshal([]byte(request.Body), &signup)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
	}
	cog_cli := initialize()
	fmt.Println(request.Body)
	userExists := cognito_auth.CheckIfUserExists(cog_cli.dynamoCli, cog_cli.tableName, signup.Username, signup.Email)
	if !userExists {
		singupResponse := cognito_auth.SignUp(cog_cli.cognitoCli, cog_cli.clientId, signup.Username, signup.Password, signup.Name, signup.Email)
		cognito_auth.AddUserToTable(cog_cli.dynamoCli, cog_cli.tableName, cognito_auth.TableItem{
			UserName: signup.Username,
			Email:    signup.Email,
		})
		notifier.AddSNSSubscription(cog_cli.snsCli, cog_cli.topicArn, signup.Email)
		return events.APIGatewayProxyResponse{
			StatusCode: singupResponse.Status,
			Body:       singupResponse.Message,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Username already exists in the table",
	}, nil
}

func handlerSignIn(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.HTTPMethod == "POST" {
		var signin Signin
		err := json.Unmarshal([]byte(request.Body), &signin)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
		}
		cog_cli := initialize()
		fmt.Println(request.Body)
		singupResponse := cognito_auth.SignIn(cog_cli.cognitoCli, cog_cli.clientId, signin.Username, signin.Password)
		return events.APIGatewayProxyResponse{
			StatusCode: singupResponse.Status,
			Body:       singupResponse.Message,
		}, nil
	}
	return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
}

func handlerResendVerificationCode(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var cogGen Coggen
	err := json.Unmarshal([]byte(request.Body), &cogGen)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
	}
	cog_cli := initialize()
	fmt.Println(request.Body)
	singupResponse := cognito_auth.ResendVerificationCode(cog_cli.cognitoCli, cog_cli.clientId, cogGen.Username)
	return events.APIGatewayProxyResponse{
		StatusCode: singupResponse.Status,
		Body:       singupResponse.Message,
	}, nil
}

func handlerForgotPassword(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var resetGen Resetgen
	err := json.Unmarshal([]byte(request.Body), &resetGen)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
	}
	cog_cli := initialize()
	fmt.Println(request.Body)
	singupResponse := cognito_auth.ForgotPassword(cog_cli.cognitoCli, cog_cli.clientId, resetGen.Username)
	return events.APIGatewayProxyResponse{
		StatusCode: singupResponse.Status,
		Body:       singupResponse.Message,
	}, nil
}

func handlerConfirmSignup(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var signupConfirm Resetgen
	err := json.Unmarshal([]byte(request.Body), &signupConfirm)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
	}
	cog_cli := initialize()
	fmt.Println(request.Body)
	singupResponse := cognito_auth.ConfirmSignUp(cog_cli.cognitoCli, cog_cli.clientId, signupConfirm.Username, signupConfirm.Confcode)
	return events.APIGatewayProxyResponse{
		StatusCode: singupResponse.Status,
		Body:       singupResponse.Message,
	}, nil
}

func handlerConfirmForgotPassword(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var resetPass ResetPass
	err := json.Unmarshal([]byte(request.Body), &resetPass)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Bad Request"}, nil
	}
	cog_cli := initialize()
	fmt.Println(request.Body)
	singupResponse := cognito_auth.ConfirmForgetPassword(cog_cli.cognitoCli, cog_cli.clientId, resetPass.Username, resetPass.Confcode, resetPass.Password)
	return events.APIGatewayProxyResponse{
		StatusCode: singupResponse.Status,
		Body:       singupResponse.Message,
	}, nil
}

func handlerPath(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := request.Path

	var handlerFunc func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
	fmt.Println(request)
	fmt.Println("Entering Handler path")
	switch path {
	case "/auth-signin":
		handlerFunc = handlerSignIn
	case "/auth-signup":
		handlerFunc = handlerSignUp
	case "/auth-resend-verf-code":
		handlerFunc = handlerResendVerificationCode
	case "/auth-forgot-password":
		handlerFunc = handlerForgotPassword
	case "/auth-confirm-signup":
		handlerFunc = handlerConfirmSignup
	case "/auth-confirm-forgot-password":
		handlerFunc = handlerConfirmForgotPassword
	default:
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Not Found"}, nil

	}
	response, err := handlerFunc(ctx, request)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
	}

	responseJSON, err := json.Marshal(response.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Internal Server Error"}, err
	}

	// Return the response
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

func main() {
	initialize()

	lambda.Start(handlerPath)
}
