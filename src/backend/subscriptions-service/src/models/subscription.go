package models

type SubscriptionDynamodb struct {
	UserName        string  `json:"username"`
	UUID            string  `json:"uuid"`
	Name            string  `json:"name"`
	Url             string  `json:"url"`
	SettingsUrl     string  `json:"settings_url"`
	Plan            string  `json:"plan"`
	StartDate       string  `json:"start_date"`
	Cost            float32 `json:"cost"`
	Icon            string  `json:"icon"`
	LastPaymentDate string  `json:"last_payment_date"`
}

type PaymentsDynamodb struct {
	UUID           string  `json:"uuid"`
	SubscriptionID string  `json:"subscription_id"`
	Amount         float32 `json:"amount"`
	PaymentDate    string  `json:"payment_date"`
}

type SubscriptionUpdate struct {
	Name            string  `json:"name"`
	Plan            string  `json:"plan"`
	StartDate       string  `json:"start_date"`
	Cost            float32 `json:"cost"`
	LastPaymentDate string  `json:"last_payment_date"`
}
