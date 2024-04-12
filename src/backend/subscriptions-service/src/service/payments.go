package service

import (
	"strconv"
	"subHandler/src/models"
	"subHandler/src/repository"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func AddPayment(item models.PaymentCreateInput) (models.PaymentDynamodb, error) {
	/*
		Adds a given Item to the DynamoDB table.
		Params: dynamoClient *dynamodb.DynamoDB
			    tableName
				item models.PaymentCreateInput
		Return: None
	*/
	uuid := uuid.New().String()
	amountFloat, convErr := strconv.ParseFloat(item.Amount, 32)
	if convErr != nil {
		log.Error().Err(convErr).Str("SubscriptionId", uuid).Str("UserName", item.UserName).Msg("Error converting cost to float")
		return models.PaymentDynamodb{}, convErr
	}
	paymentNew := models.PaymentDynamodb{
		UUID:           uuid,
		SubscriptionId: item.SubscriptionId,
		UserName:       item.UserName,
		Amount:         float32(amountFloat),
		PaymentDate:    item.PaymentDate,
	}

	log.Info().Str("UUID", uuid).Str("SubscriptionId", item.SubscriptionId).Msg("Adding payment")
	res, err := repository.AddSubscriptionPayment(paymentNew)
	if err != nil {
		log.Error().Err(err).Str("UUID", uuid).Str("SubscriptionId", item.SubscriptionId).Msg("Error adding payment")
		return models.PaymentDynamodb{}, err
	}
	log.Info().Str("UUID", uuid).Str("SubscriptionId", item.SubscriptionId).Msg("Payment added")
	return res, nil
}

func GetPayments(partitionKey string) ([]models.PaymentDynamodb, error) {
	/*
		Returns all the payments for a given subscription.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
		Return: []models.PaymentDynamodb, error
	*/
	log.Info().Str("SubscriptionId", partitionKey).Msg("Getting payments")
	res, err := repository.GetSubscriptionPayments(partitionKey)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", partitionKey).Msg("Error getting payments")
		return []models.PaymentDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", partitionKey).Msg("Payments retrieved")
	return res, nil
}
