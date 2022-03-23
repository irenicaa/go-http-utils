package auth

import (
	"encoding/base64"
	"os"
)

// MakeBasicAuthHeader ...
func MakeBasicAuthHeader(usernameEnv string, passwordEnv string) string {
	username, password := os.Getenv(usernameEnv), os.Getenv(passwordEnv)
	if username == "" || password == "" {
		return ""
	}

	credentials := username + ":" + password
	credentials = base64.StdEncoding.EncodeToString([]byte(credentials))

	return "Basic " + credentials
}

// MakeBearerAuthHeader ...
func MakeBearerAuthHeader(tokenEnv string) string {
	token := os.Getenv(tokenEnv)
	if token == "" {
		return ""
	}

	return "Bearer " + token
}
