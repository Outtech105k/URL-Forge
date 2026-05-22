package models

type SetUrlRequest struct {
	BaseURL      string    `json:"base_url" binding:"required,url"`
	CustomID     *string   `json:"custom_id" binding:"omitempty"`
	UseUppercase *bool     `json:"use_uppercase"`
	UseLowercase *bool     `json:"use_lowercase"`
	UseNumbers   *bool     `json:"use_numbers"`
	IDLength     *uint32   `json:"id_length"`
	ExpireIn     *Duration `json:"expire_in"`
	SandCushion  *bool     `json:"sand_cushion"`
	IsPublicCtrl *bool     `json:"public_ctrl"`
}
