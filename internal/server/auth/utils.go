package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"
)

func getUsers() []PortalUser {
	file, err := os.Open("users.json")
	if err != nil {
		panic(err)
	}

	var users []PortalUser

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&users); err != nil {
		panic(err)
	}
	file.Close()

	return users
}

func checkUser(email string, password string) (loginSuccessful bool) {
	for _, user := range getUsers() {
		if user.Email == email && user.PasswordHash == hashPassword(password) {
			return true
		}
	}

	return false
}

func getUser(email string) (PortalUser, error) {
	for _, user := range getUsers() {
		if user.Email == email {
			return user, nil
		}
	}

	return PortalUser{}, errors.New("user not found")
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

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}
