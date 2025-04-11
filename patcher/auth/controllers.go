package auth

import (
	"encoding/json"
	"net/http"
)

func GetSigninPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./patcher/static/login.html")
}

func Signin(users []PortalUser) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user LoginUser

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		token, err := user.Login(users)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   3600,
		})
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}
