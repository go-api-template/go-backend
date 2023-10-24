package token

// AccessToken is a struct for the access token
// It is used to return the access token to the user
type AccessToken struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	ExpiresAt        int64  `json:"expires_at"`
	TokenType        string `json:"token_type"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshExpiresAt int64  `json:"refresh_expires_at"`
}
