package models

type SubscriptionCreateInput struct {
	UserName    string `json:"username"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	SettingsUrl string `json:"settings_url"`
	Plan        string `json:"plan"`
	Cost        string `json:"cost"`
	StartDate   string `json:"start_date"`
}
