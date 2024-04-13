package models

type SubscriptionCategory string

const (
	OTT       SubscriptionCategory = "ott"
	Music     SubscriptionCategory = "music"
	Gaming    SubscriptionCategory = "gaming"
	Delivery  SubscriptionCategory = "delivery"
	Fittness  SubscriptionCategory = "fittness"
	Education SubscriptionCategory = "education"
	Magzine   SubscriptionCategory = "magzine"
	Software  SubscriptionCategory = "software"
	Finance   SubscriptionCategory = "finance"
	Fashion   SubscriptionCategory = "fashion"
	Other     SubscriptionCategory = "other"
)

type SubscriptionDynamodb struct {
	UserName        string               `json:"username"`
	UUID            string               `json:"uuid"`
	Name            string               `json:"name"`
	Url             string               `json:"url"`
	SettingsUrl     string               `json:"settings_url"`
	Plan            string               `json:"plan"`
	StartDate       string               `json:"start_date"`
	Cost            float32              `json:"cost"`
	Icon            string               `json:"icon"`
	LastPaymentDate string               `json:"last_payment_date"`
	Category        SubscriptionCategory `json:"category"`
}

type SubscriptionUpdate struct {
	Name            string  `json:"name"`
	Plan            string  `json:"plan"`
	StartDate       string  `json:"start_date"`
	Cost            float32 `json:"cost"`
	LastPaymentDate string  `json:"last_payment_date"`
	Category        string  `json:"category"`
}
