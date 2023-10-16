package models

// AccessToken is a struct for the access token
// It is used to return the access token to the user
type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ExpiresAt   int64  `json:"expires_at"`
	TokenType   string `json:"token_type"`
}
