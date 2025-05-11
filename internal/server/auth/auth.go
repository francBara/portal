package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PortalUser struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

var jwtSecret = []byte("My secret")

type tokenResponse struct {
	User  PortalUser `json:"user"`
	Token string     `json:"token"`
}

// Signin accepts Basic authentication header and returns a JWT id token.
func Signin() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, password, err := decodeBasicAuth(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		slog.Info("Login attempt", "email", email)

		if !checkUser(email, password) {
			http.Error(w, "Email or password are not correct", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{
			"sub": email,
			"exp": time.Now().Add(time.Hour * 2).Unix(),
			"iat": time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := getUser(email)
		if err != nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenResponse{
			User: PortalUser{
				Email: email,
				Name:  user.Name,
			},
			Token: tokenString,
		})

		slog.Info(fmt.Sprintf("User %s logged in", email))
	}
}

// AuthenticateUser is a middleware that verifies Bearer token authorization to protect API routes.
func AuthenticateUser() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bearerToken, err := decodeBearerAuth(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			subject, err := token.Claims.GetSubject()
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}

			user, err := getUser(subject)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user", &user))

			next.ServeHTTP(w, r)
		})
	}
}
