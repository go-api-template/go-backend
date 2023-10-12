package utils

import "github.com/thanhpk/randstr"

// GenerateRandomString generates a random string of the given length
func GenerateRandomString(length int) string {
	return randstr.String(length)
}
