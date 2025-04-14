package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
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

func (user LoginUser) Login(users []PortalUser) (string, error) {
	for _, u := range users {
		if u.Email == user.Email && u.PasswordHash == hashPassword(user.Password) {

			claims := jwt.MapClaims{
				"email": user.Email,
				"exp":   time.Now().Add(120 * time.Minute).Unix(),
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			tokenString, err := token.SignedString(jwtSecret)
			if err != nil {
				return "", err
			}

			return tokenString, nil
		}
	}

	return "", errors.New("could not login")
}

func AuthenticateUser(users []PortalUser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil || cookie.Value == "" {
				http.Redirect(w, r, "/signin", http.StatusSeeOther)
				return
			}

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			email := claims["email"]

			user, err := getUser(users, email.(string))
			if err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user", &user))

			next.ServeHTTP(w, r)
		})
	}
}
