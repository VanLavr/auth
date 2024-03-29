package models

type RefreshToken struct {
	GUID        string `json:"guid"`
	IsUsed      bool   `json:"-"`
	TokenString string `json:"refresh_token"`
}
