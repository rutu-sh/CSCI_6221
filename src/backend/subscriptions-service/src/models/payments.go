package models

type PaymentDynamodb struct {
	SubscriptionId string  `json:"subscription_id"`
	UUID           string  `json:"uuid"`
	UserName       string  `json:"username"`
	Amount         float32 `json:"amount"`
	PaymentDate    string  `json:"payment_date"`
}

type PaymentUpdate struct {
	Amount      float32 `json:"amount"`
	PaymentDate string  `json:"payment_date"`
}
