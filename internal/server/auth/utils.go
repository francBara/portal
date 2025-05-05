package auth

import (
	"encoding/base64"
	"errors"
	"strings"
)

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
