package token

import (
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Create creates a new JWT token
// from: the time the token is created
// ttl: time to live in minutes
// issuer: the issuer of the token
// subject: the subject of the token
// audience: the audience of the token
// privateKey: the private key to sign the token with
func Create(from time.Time, ttl int, issuer string, subject string, audience string, privateKey string) (string, error) {
	// Private key used for signing the token
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   subject,
		Audience:  []string{audience},
		IssuedAt:  jwt.NewNumericDate(from),
		NotBefore: jwt.NewNumericDate(from),
		ExpiresAt: jwt.NewNumericDate(from.Add(time.Duration(ttl) * time.Minute)),
	}

	// Create the token using the claims and sign it using the key
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func Validate(tokenString string, publicKey string) (jwt.Claims, error) {
	// Public key used for validating the token
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error while validating: %w", err)
	}

	// Validate the token
	if token.Valid {
		// Get the claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("validate: invalid claims")
		}
		//
		return claims, nil
	} else {
		return nil, err
	}
}
