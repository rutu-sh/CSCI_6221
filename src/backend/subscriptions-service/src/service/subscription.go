package service

import (
	"strconv"
	"subHandler/src/models"
	"subHandler/src/repository"

	"github.com/google/uuid"

	"github.com/rs/zerolog/log"
)

func AddSubscription(item models.SubscriptionCreateInput) (models.SubscriptionDynamodb, error) {
	/*
		Adds a given Item to the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				item models.SubscriptionCreateInput
		Return: None
	*/

	uuid := uuid.New().String()

	costFloat, convErr := strconv.ParseFloat(item.Cost, 32)
	if convErr != nil {
		log.Error().Err(convErr).Str("SubscriptionId", uuid).Str("UserName", item.UserName).Str("Name", item.Name).Str("Url", item.Url).Msg("Error converting cost to float")
		return models.SubscriptionDynamodb{}, convErr
	}

	subNew := models.SubscriptionDynamodb{
		UUID:            uuid,
		UserName:        item.UserName,
		Name:            item.Name,
		Url:             item.Url,
		SettingsUrl:     item.SettingsUrl,
		Plan:            item.Plan,
		Cost:            float32(costFloat),
		StartDate:       item.StartDate,
		Icon:            "https://via.placeholder.com/150",
		LastPaymentDate: item.StartDate,
	}
	log.Info().Str("SubscriptionId", uuid).Str("UserName", item.UserName).Str("Name", item.Name).Str("Url", item.Url).Msg("Adding subscription")
	res, err := repository.AddSubscription(subNew)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", uuid).Str("UserName", item.UserName).Str("Name", item.Name).Str("Url", item.Url).Msg("Error adding subscription")
		return models.SubscriptionDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", uuid).Str("UserName", item.UserName).Str("Name", item.Name).Str("Url", item.Url).Msg("Subscription added")
	return res, nil
}

func GetSubscription(subscriptionId string, userName string) (models.SubscriptionDynamodb, error) {
	/*
		Gets a given Item from the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				subscriptionId
				userName
		Return: models.SubscriptionDynamodb, error
	*/

	log.Info().Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Getting subscription")
	item, err := repository.GetSubscription(userName, subscriptionId)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Error getting subscription")
		return models.SubscriptionDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Subscription retrieved")
	return item, nil
}

func DeleteSubscription(subscriptionId string, userName string) error {
	/*
		Deletes a given Subscription from the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName string
				subscriptionId string
				userName string
		Return: error
	*/
	log.Info().Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Deleting subscription")
	err := repository.DeleteSubscription(userName, subscriptionId)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Error deleting subscription")
		return err
	}
	log.Info().Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Subscription deleted")
	return nil
}

func UpdateSubscription(subscriptionId string, userName string, updateItem models.SubscriptionUpdate) (models.SubscriptionDynamodb, error) {
	/*
		Updates a given Subscription in the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				subscriptionId
				userName
				updateItem models.SubscriptionUpdateInput
		Return: models.SubscriptionDynamodb, error
	*/
	log.Info().Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Updating subscription")
	updatedSubscription, err := repository.UpdateSubscription(userName, subscriptionId, updateItem)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Error updating subscription")
		return models.SubscriptionDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", subscriptionId).Str("UserName", userName).Msg("Subscription updated")
	return updatedSubscription, nil
}

func GetUserSubscriptions(userName string) ([]models.SubscriptionDynamodb, error) {
	/*
		Gets all Subscriptions from the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				userName
		Return: []models.SubscriptionDynamodb, error
	*/
	log.Info().Str("UserName", userName).Msg("Getting all subscriptions")
	items, err := repository.GetUserSubscriptions(userName)
	if err != nil {
		log.Error().Err(err).Str("UserName", userName).Msg("Error getting all subscriptions")
		return nil, err
	}
	log.Info().Str("UserName", userName).Msg("All subscriptions retrieved")
	return items, nil
}
