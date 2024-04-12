package repository

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"subHandler/src/config"
	"subHandler/src/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/rs/zerolog/log"
)

func initialize() models.DynamoAttr {
	/*
		Used to initialize the table attributes and sdk clients
		Params: None
		Return: None
	*/
	awsRegion := config.AWS_REGION
	dynamodbTable := config.DYNAMODB_TABLE_NAME

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		fmt.Println("Error creating sess", err)
		os.Exit(0)
	}

	dynamoClient := dynamodb.New(sess)
	return models.DynamoAttr{
		DynamoCli: dynamoClient,
		AwsRegion: awsRegion,
		TableName: dynamodbTable,
	}
}

func AddSubscription(item models.SubscriptionDynamodb) (models.SubscriptionDynamodb, error) {
	/*
		Adds a given Item to the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				item models.SubscriptionDynamodb
		Return: models.SubscriptionDynamodb, error
	*/

	da := initialize()
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Msg("Adding subscription")
	mappedItem, _ := dynamodbattribute.MarshalMap(item)
	tableInput := &dynamodb.PutItemInput{
		Item:      mappedItem,
		TableName: aws.String(tableName),
	}
	_, err := dynamoClient.PutItem(tableInput)
	if err != nil {
		log.Error().Err(err).Msg("Error adding subscription")
		return models.SubscriptionDynamodb{}, err
	}
	log.Info().Msg("Subscription added")
	return item, nil
}

func GetSubscription(partitionKey string, sortKey string) (models.SubscriptionDynamodb, error) {
	/*
		Gets a given Item from the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				partitionKey
				sortKey
		Return: models.SubscriptionDynamodb, error
	*/
	da := initialize()
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Msg("Getting subscription")
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(partitionKey),
			},
			"uuid": {
				S: aws.String(sortKey),
			},
		},
	}

	result, err := dynamoClient.GetItem(input)
	if err != nil {
		log.Error().Err(err).Msg("Error getting subscription")
		return models.SubscriptionDynamodb{}, err
	}

	if len(result.Item) == 0 {
		log.Error().Msg("Error getting subscription. No item found.")
		return models.SubscriptionDynamodb{}, errors.New("404")
	}

	item := models.SubscriptionDynamodb{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		log.Error().Err(err).Msg("Error getting subscription")
		return models.SubscriptionDynamodb{}, err
	}

	log.Info().Msg("Subscription retrieved")
	return item, nil
}

func DeleteSubscription(partitionKey string, sortKey string) error {
	/*
		Deletes a given Item from the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				partitionKey
				sortKey
		Return: error
	*/

	da := initialize()
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Msg("Deleting subscription")
	_, sub_get_err := GetSubscription(partitionKey, sortKey)
	if sub_get_err != nil {
		log.Error().Err(sub_get_err).Msg("Error deleting subscription. Subscription Not found")
		return sub_get_err
	}

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(partitionKey),
			},
			"uuid": {
				S: aws.String(sortKey),
			},
		},
	}

	_, err := dynamoClient.DeleteItem(input)
	if err != nil {
		log.Error().Err(err).Msg("Error deleting subscription")
		return err
	}

	log.Info().Msg("Subscription deleted")
	return nil
}

func UpdateSubscription(partitionKey string, sortKey string, updateItem models.SubscriptionUpdate) (models.SubscriptionDynamodb, error) {
	/*
		Updates a given Item in the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				partitionKey
				sortKey
				updateItem models.SubscriptionUpdate
		Return: models.SubscriptionDynamodb, error
	*/
	da := initialize()
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Msg("Updating subscription")
	subscription, err := GetSubscription(partitionKey, sortKey)
	if err != nil {
		return models.SubscriptionDynamodb{}, err
	}

	newSubscription := models.SubscriptionDynamodb{
		UserName:        partitionKey,
		UUID:            sortKey,
		Name:            updateItem.Name,
		Url:             subscription.Url,
		SettingsUrl:     subscription.SettingsUrl,
		Plan:            updateItem.Plan,
		StartDate:       updateItem.StartDate,
		Cost:            updateItem.Cost,
		LastPaymentDate: updateItem.LastPaymentDate,
		Icon:            subscription.Icon,
	}
	tableInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(partitionKey),
			},
			"uuid": {
				S: aws.String(sortKey),
			},
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(updateItem.Name),
			},
			":plan": {
				S: aws.String(updateItem.Plan),
			},
			":start_date": {
				S: aws.String(updateItem.StartDate),
			},
			":cost": {
				N: aws.String(strconv.FormatFloat(float64(updateItem.Cost), 'f', -1, 32)),
			},
			":last_payment_date": {
				S: aws.String(updateItem.LastPaymentDate),
			},
		},
		ReturnValues:     aws.String("ALL_NEW"),
		UpdateExpression: aws.String("SET #name = :name, #plan = :plan, #start_date = :start_date, #cost = :cost, #last_payment_date = :last_payment_date"),
		ExpressionAttributeNames: map[string]*string{
			"#name":              aws.String("name"),
			"#plan":              aws.String("plan"),
			"#start_date":        aws.String("start_date"),
			"#cost":              aws.String("cost"),
			"#last_payment_date": aws.String("last_payment_date"),
		},
	}
	_, err = dynamoClient.UpdateItem(tableInput)
	if err != nil {
		log.Error().Err(err).Msg("Error updating subscription")
		return models.SubscriptionDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", sortKey).Str("UserName", partitionKey).Msg("Subscription updated")
	return newSubscription, nil
}

func GetUserSubscriptions(partitionKey string) ([]models.SubscriptionDynamodb, error) {
	/*
		Gets all the Items for a given User from the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				partitionKey
		Return: []models.SubscriptionDynamodb, error
	*/
	da := initialize()
	dynamoClient := da.DynamoCli
	tableName := da.TableName

	log.Info().Str("UserName", partitionKey).Msg("Getting user subscriptions")

	// query the dynamodb table using the partition key
	input := &dynamodb.QueryInput{
		TableName: aws.String(tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"username": {
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
		log.Error().Err(err).Msg("Error getting user subscriptions")
		return nil, err
	}

	items := []models.SubscriptionDynamodb{}
	for _, i := range result.Items {
		item := models.SubscriptionDynamodb{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			log.Error().Err(err).Msg("Error getting user subscriptions")
			return nil, err
		}
		items = append(items, item)
	}

	log.Info().Str("UserName", partitionKey).Int("SubscriptionCount", len(items)).Msg("User subscriptions retrieved successfully")
	return items, nil
}
