# ![Go](./assets/go.png) CSCI_6221 - GoLang Group 

### Members

1. Namana Y Tarikere
2. Ruturaj Shitole
3. Srinivas Ravindranath
4. Tejaaswini Narendran



# SubHub - A Subscription Management Platform

## Idea

The basic idea of this project is to create a service for managing user subscriptions and providing alerts based on when the subscription is nearing it's end so that the user can make a decision on continuing/discontinuing the subscription. 

This project compromises of the following components (written in Go):

1. **Authentication**: This component is used for authenticating and managing the users. 
2. **Subscription Management**: This component is used for managing the subscriptions of the users.
3. **Alerting**: This component is used for sending alerts to the users when their subscriptions are nearing their end.

## Language of choice 

We have used the [Go Programming Language](https://go.dev/) for writing the backend for this project. 

## Architecture

### 1. Authentication

We have used AWS Cognito for authentication. The flow is as follows:

1. User sends an authentication request to the API Gateway
2. API Gateway forwards the request to the Authentication Lambda function
3. The Authentication Lambda function validates the user credentials with AWS Cognito
4. If the user is authenticated, the Authentication Lambda function generates a JWT token and sends it back to the user

![Authentication](./assets/golang-auth-flow.drawio.png)

The main component of this flow is the Authentication Lambda function. This function is responsible for validating the user credentials with AWS Cognito and generating a JWT token. This function is written in **Go**. It uses the `github.com/aws/aws-lambda-go/lambda` package to handle the request and response. It uses the `github.com/aws/aws-sdk-go` package to interact with AWS Cognito.

The source code for this function can be found in the [src/backend/auth](/src/backend/auth). 


### 2. Subscription Management

We have used AWS Lambda and DynamoDB for managing the subscriptions. The flow is as follows:

1. User sends a request to the API Gateway
2. API Gateway forwards the request to the Subscription Management Lambda function
3. The Subscription Management Lambda function interacts with DynamoDB to perform CRUD operations on the subscriptions
4. The Subscription Management Lambda function sends the response back to the user

![Subscription Management](./assets/subscription-management-flow.drawio.png)

The main component of this flow is the Subscription Management Lambda function. This function is responsible for interacting with DynamoDB to perform CRUD operations on the subscriptions. This function is written in **Go**. It uses the `github.com/aws/aws-lambda-go/lambda` package to handle the request and response. It uses the `github.com/aws/aws-sdk-go` package to interact with DynamoDB.

The source code for this function can be found in the [src/backend/service](/src/backend/service).

### 3. Alerting

We have used AWS Lambda and EventBridge for sending alerts to the users. The flow is as follows:

1. A scheduled event is triggered in EventBridge
2. EventBridge triggers the Alerting Lambda function
3. The Alerting Lambda function interacts with DynamoDB to get the list of subscriptions that are nearing their end
4. The Alerting Lambda function sends an email to the users with the list of subscriptions that are nearing their end

![Subscription Alerter](./assets/subscription-sns-alerter.drawio.png)

The main component of this flow is the Alerting Lambda function. This function is responsible for interacting with DynamoDB to get the list of subscriptions that are nearing their end and sending an email to the users. This function is written in **Go**. It uses the `github.com/aws/aws-lambda-go/lambda` package to handle the request and response. It uses the `github.com/aws/aws-sdk-go` package to interact with DynamoDB and SNS.