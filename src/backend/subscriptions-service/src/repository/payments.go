package repository

import (
	"errors"
	"strconv"
	"subHandler/src/config"
	"subHandler/src/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)

func IsPaymentExists(dynamoClient *dynamodb.DynamoDB, tableName string, partitionKey string, sortKey string) bool {
	/*
		Checks if a given Item exists in the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
				sortKey
		Return: bool
	*/
	log.Info().Msg("Checking if payment exists")
	res, err := dynamoClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"subscription_id": {
				S: aws.String(partitionKey),
			},
			"uuid": {
				S: aws.String(sortKey),
			},
		},
	})
	if err != nil {
		log.Error().Err(err).Msg("Error checking if payment exists")
		return false
	}
	if len(res.Item) == 0 {
		log.Info().Msg("Payment does not exist")
		return false
	}
	log.Info().Msg("Payment exists")
	return true
}

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

func GetSubscriptionPayment(partitionKey string, sortKey string) (models.PaymentDynamodb, error) {
	/*
		Returns a payment for a given subscription.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
				sortKey
		Return: models.PaymentDynamodb, error
	*/

	da := initialize("payments")
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Getting payment")

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"subscription_id": {
				S: aws.String(partitionKey),
			},
			"uuid": {
				S: aws.String(sortKey),
			},
		},
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		log.Error().Err(err).Msg("Error getting payment")
		return models.PaymentDynamodb{}, err
	}

	item := models.PaymentDynamodb{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Error().Err(err).Msg("Error getting payment")
		return models.PaymentDynamodb{}, err
	}

	log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Payment retrieved successfully")
	return item, nil
}

func UpdateSubscriptionPayment(partitionKey string, sortKey string, updateItem models.PaymentUpdate) (models.PaymentDynamodb, error) {
	/*
		Updates a given Item in the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
				sortKey
				updateItem
		Return: models.PaymentDynamodb, error
	*/
	da := initialize("payments")
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Updating payment")

	// subscriptionExists := IsSubscriptionExists(dynamoClient, config.SUBSCRIPTIONS_DYNAMODB_TABLE, updateItem.UserName, partitionKey)
	paymentExists := IsPaymentExists(dynamoClient, tableName, partitionKey, sortKey)

	if !paymentExists {
		log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Payment does not exist")
		return models.PaymentDynamodb{}, errors.New("payment does not exist")
	}

	payment, _ := GetSubscriptionPayment(partitionKey, sortKey)
	newPayment := models.PaymentDynamodb{
		SubscriptionId: partitionKey,
		UUID:           sortKey,
		UserName:       payment.UserName,
		Amount:         updateItem.Amount,
		PaymentDate:    updateItem.PaymentDate,
	}
	tableInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"subscription_id": {
				S: aws.String(partitionKey),
			},
			"uuid": {
				S: aws.String(sortKey),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":a": {
				N: aws.String(strconv.FormatFloat(float64(updateItem.Amount), 'f', -1, 32)),
			},
			":d": {
				S: aws.String(updateItem.PaymentDate),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET amount = :a, payment_date = :d"),
	}
	_, err := dynamoClient.UpdateItem(tableInput)
	if err != nil {
		log.Error().Err(err).Msg("Error updating payment")
		return models.PaymentDynamodb{}, err
	}

	log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Payment updated successfully")
	return newPayment, nil
}

func DeleteSubscriptionPayment(partitionKey string, sortKey string) error {
	/*
		Deletes a given Item from the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
				sortKey
		Return: error
	*/
	da := initialize("payments")
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Deleting payment")

	paymentExists := IsPaymentExists(dynamoClient, tableName, partitionKey, sortKey)
	if !paymentExists {
		log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Payment does not exist")
		return errors.New("payment does not exist")
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"subscription_id": {
				S: aws.String(partitionKey),
			},
			"uuid": {
				S: aws.String(sortKey),
			},
		},
	}

	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting payment")
		return err
	}

	log.Info().Str("SubscriptionId", partitionKey).Str("PaymentId", sortKey).Msg("Payment deleted successfully")
	return nil
}
