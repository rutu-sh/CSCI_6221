package cognito_auth

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"os"
)

type ReturnResults struct {
	Status  int    `json:"Status"`
	Message string `json:"Message"`
}

type TableItem struct {
	UserName string
	Email    string
}

type GetItem struct {
	UserName string
	Email    string
}

func CheckIfUserExists(dynamoClient *dynamodb.DynamoDB, tableName, userName, email string) bool {

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"UserName": {
				S: aws.String(userName),
			},
			"Email": {
				S: aws.String(email),
			},
		},
	}

	result, derr := dynamoClient.GetItem(input)
	log.Println(result)
	if derr != nil {
		if aerr, ok := derr.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Printf("Throughput exceeded for table %s please try again after sometime", tableName)
				os.Exit(1)
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Printf("Table %s not found, or item not found please check tablename/item", tableName)
				os.Exit(1)
			}
		} else {
			log.Printf("Got error calling GetItem: %s", derr)
			os.Exit(1)
		}
	}

	item := GetItem{}
	err := dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Printf("Failed to unmarshal Record, %v", err)
		os.Exit(1)
	}

	log.Printf("Successfully retreived Item from table %s", tableName)
	if (GetItem{}) == item {
		return false
	}
	return true
}

func SignUp(cognitoClient *cognitoidentityprovider.CognitoIdentityProvider, clientId, username, password, name, email string) ReturnResults {

	var res ReturnResults
	fmt.Println(clientId, username, password, name, email)
	input := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(clientId),
		Username: aws.String(username),
		Password: aws.String(password),
		UserAttributes: []*cognitoidentityprovider.AttributeType{
			{
				Name:  aws.String("name"),
				Value: aws.String(name),
			},
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
		},
	}
	// TODO: check if email already exists
	_, err := cognitoClient.SignUp(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cognitoidentityprovider.ErrCodeResourceNotFoundException:
				res.Message = fmt.Sprintf("ClientId %s does not exist", clientId)
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeInvalidPasswordException:
				res.Message = fmt.Sprintf("Password doesnt conform to rules specified")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeUsernameExistsException:
				res.Message = fmt.Sprintf("User %s already exists in user pool", username)
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeCodeDeliveryFailureException:
				res.Message = fmt.Sprintf("Verification code failed to be sent to user")
				res.Status = 400
				return res
			case cognitoidentityprovider.ErrCodeInvalidParameterException:
				res.Message = fmt.Sprintf("Please check if the password satisfies all contraints")

			}
			fmt.Println(aerr.Error())
		} else {
			fmt.Println(err.Error())
		}
	}
	res.Message = fmt.Sprintf("Succesfully signed up user %s", username)
	res.Status = 200
	return res
}

func AddUserToTable(dynamoClient *dynamodb.DynamoDB, tableName string, tableItems TableItem) {
	mappedItem, _ := dynamodbattribute.MarshalMap(tableItems)
	tableInput := &dynamodb.PutItemInput{
		Item:      mappedItem,
		TableName: aws.String(tableName),
	}
	_, err := dynamoClient.PutItem(tableInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				log.Printf("Conditional Check failed for item addition")
				os.Exit(1)
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				log.Printf("Throughput exceeded for table %s please try again after sometime", tableName)
				os.Exit(1)
			case dynamodb.ErrCodeResourceNotFoundException:
				log.Printf("Table %s not found, please check tablename", tableName)
				os.Exit(1)
			case dynamodb.ErrCodeTransactionConflictException:
				log.Printf("Transaction already in progress for itme, please try again after sometime")
				os.Exit(1)
			}
		} else {
			log.Printf("Unable to add item to table: %s\n", err)
			os.Exit(1)
		}
	}
	log.Printf("Successfully added the items into the table %s", tableName)
}
