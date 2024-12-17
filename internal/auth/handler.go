package auth

import (
	"encoding/json"
	"net/http"
	"server/internal/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := CreateUser(req.Username, req.Password); err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User registered successfully"))
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if !AuthenticateUser(req.Username, req.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(req.Username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: token,
		Path:  "/",
	})
	w.Write([]byte("Login successful"))
}
