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

func GetPayments(subscriptionId string) ([]models.PaymentDynamodb, error) {
	/*
		Returns all the payments for a given subscription.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
		Return: []models.PaymentDynamodb, error
	*/
	log.Info().Str("SubscriptionId", subscriptionId).Msg("Getting payments")
	res, err := repository.GetSubscriptionPayments(subscriptionId)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", subscriptionId).Msg("Error getting payments")
		return []models.PaymentDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", subscriptionId).Msg("Payments retrieved")
	return res, nil
}

func GetPayment(subscriptionId string, paymentId string) (models.PaymentDynamodb, error) {
	/*
		Returns a payment for a given subscription.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
				sortKey
		Return: models.PaymentDynamodb, error
	*/
	log.Info().Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Getting payment")
	res, err := repository.GetSubscriptionPayment(subscriptionId, paymentId)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Error getting payment")
		return models.PaymentDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Payment retrieved")
	return res, nil
}

func UpdatePayment(subscriptionId string, paymentId string, item models.PaymentUpdate) (models.PaymentDynamodb, error) {
	/*
		Updates a payment for a given subscription.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
				sortKey
				item models.PaymentUpdate
		Return: models.PaymentDynamodb, error
	*/
	log.Info().Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Updating payment")
	res, err := repository.UpdateSubscriptionPayment(subscriptionId, paymentId, item)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Error updating payment")
		return models.PaymentDynamodb{}, err
	}
	log.Info().Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Payment updated")
	return res, nil
}

func DeletePayment(subscriptionId string, paymentId string) error {
	/*
		Deletes a payment for a given subscription.
		Params: dynamoClient *dynamodb.DynamoDB
				tableName
				partitionKey
				sortKey
		Return: error
	*/
	log.Info().Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Deleting payment")
	err := repository.DeleteSubscriptionPayment(subscriptionId, paymentId)
	if err != nil {
		log.Error().Err(err).Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Error deleting payment")
		return err
	}
	log.Info().Str("SubscriptionId", subscriptionId).Str("PaymentId", paymentId).Msg("Payment deleted")
	return nil
}
