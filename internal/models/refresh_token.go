package models

type RefreshToken struct {
	GUID        string `json:"guid"`
	TokenString string `json:"refresh_token"`
}
