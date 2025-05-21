package auth

import (
	"encoding/base64"
	"errors"
	"os"
	"strings"
)

type Admin struct {
	Username string `json:"adminUsername"`
	Password string `json:"adminPassword"`
}

var admin Admin

func loadAdmin() error {
	admin.Username = os.Getenv("ADMIN_USERNAME")
	admin.Password = os.Getenv("ADMIN_PASSWORD")

	if admin.Username == "" {
		admin.Username = "admin"
	}
	if admin.Password == "" {
		admin.Password = "admin"
	}

	return nil
}

func checkUser(username string, password string) (loginSuccessful bool) {
	if admin.Username == "" {
		err := loadAdmin()
		if err != nil {
			panic(err)
		}
	}

	return admin.Username == username && admin.Password == password
}

func decodeBasicAuth(authHeader string) (string, string, error) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return "", "", errors.New("basic prefix not present")
	}

	payload, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(authHeader, "Basic "))
	if err != nil {
		return "", "", errors.New("invalid auth encoding")
	}

	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", errors.New("invalid header")
	}

	return pair[0], pair[1], nil
}

func decodeBearerAuth(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer") {
		return "", errors.New("invalid bearer token")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		return "", errors.New("invalid bearer token")
	}

	return parts[1], nil
}
