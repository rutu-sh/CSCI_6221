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

func GetSubscriptionPayments(partitionKey string) ([]models.PaymentDynamodb, error) {
	/*
		Returns all the payments for a given subscription.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
		Return: []models.PaymentDynamodb, error
	*/

	da := initialize("payments")
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Str("SubscriptionId", partitionKey).Msg("Getting subscription payments")

	// query the dynamodb table using the partition key
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"subscription_id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(partitionKey),
					},
				},
			},
		},
	}
	result, err := dynamoClient.Query(input)
	if err != nil {
		log.Error().Err(err).Msg("Error getting subscription payments")
		return nil, err
	}

	items := []models.PaymentDynamodb{}
	for _, i := range result.Items {
		item := models.PaymentDynamodb{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Error().Err(err).Msg("Error getting subscription payments")
			return nil, err
		}
		items = append(items, item)
	}

	log.Info().Str("SubscriptionId", partitionKey).Int("PaymentCount", len(items)).Msg("Subscription payments retrieved successfully")
	return items, nil
}
