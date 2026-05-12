package models

type SetUrlResponse struct {
	BaseURL  string   `json:"base_url"`
	ShortURL string   `json:"short_url"`
	Warnings []string `json:"warnings,omitempty"`
}
