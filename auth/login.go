package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"rouge/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fileData, err := os.ReadFile("users.json")
	if err != nil {
		http.Error(w, "Couldnt read user file", http.StatusInternalServerError)
		return
	}
	var users []models.User
	if err := json.Unmarshal(fileData, &users); err != nil {
		http.Error(w, "Error parsing users file", http.StatusInternalServerError)
		return
	}

	var foundUser models.User
	found := false

	for _, u := range users {
		if u.Username == req.Username {
			foundUser = u
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": foundUser.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "supersecretkey123"
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
