package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PortalUser struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var jwtSecret = []byte("My secret")

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func getUser(users []PortalUser, email string) (PortalUser, error) {
	for _, user := range users {
		if user.Email == email {
			return user, nil
		}
	}

	return PortalUser{}, errors.New("user not found")
}

type TokenResponse struct {
	Token string `json:"token"`
}

func Signin(users []PortalUser) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		email, password, err := decodeBasicAuth(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for _, u := range users {
			if u.Email == email && u.PasswordHash == hashPassword(password) {

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

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(TokenResponse{
					Token: tokenString,
				})

				slog.Info(fmt.Sprintf("User %s logged in", email))

				return
			}
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func AuthenticateUser(users []PortalUser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
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

			user, err := getUser(users, subject)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user", &user))

			next.ServeHTTP(w, r)
		})
	}
}
