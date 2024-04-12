package models

type PaymentsDynamodb struct {
	SubscriptionId string  `json:"subscription_id"`
	UUID           string  `json:"uuid"`
	UserName       string  `json:"username"`
	Amount         float32 `json:"amount"`
}
