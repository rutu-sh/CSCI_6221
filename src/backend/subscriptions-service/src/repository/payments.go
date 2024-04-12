package repository

import (
	"errors"
	"subHandler/src/config"
	"subHandler/src/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)

func AddSubscriptionPayment(item models.PaymentDynamodb) (models.PaymentDynamodb, error) {
	/*
		Adds a given Item to the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				item models.PaymentsDynamodb
		Return: models.PaymentsDynamodb, error
	*/

	da := initialize("payments")
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Msg("Adding subscription payment")
	if !IsSubscriptionExists(dynamoClient, config.SUBSCRIPTIONS_DYNAMODB_TABLE, item.UserName, item.SubscriptionId) {
		log.Info().Msg("Subscription does not exists")
		return item, errors.New("subscription does not exist")
	}

	mappedItem, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Error().Err(err).Msg("Error adding payment")
		return item, err
	}

	tableInput := &dynamodb.PutItemInput{
		Item:      mappedItem,
		TableName: aws.String(tableName),
	}

	_, err = dynamoClient.PutItem(tableInput)
	if err != nil {
		log.Error().Err(err).Msg("Error adding payment")
		return item, err
	}

	log.Info().Msg("Payment added successfully")
	return item, nil
}
