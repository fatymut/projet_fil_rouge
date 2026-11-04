package auth

import (
	"encoding/json"
	"fmt"
	"os"

	"rouge/models"

	"golang.org/x/crypto/bcrypt"
)

func HandleCli(args []string) bool {
	if len(args) >= 5 && args[1] == "u" && args[3] == "-p" {
		username := args[2] // "admin"
		password := args[4] // "password"

		fmt.Printf("Creating user %s...\n", username)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error hashing password: %v\n", err)
			return true
		}

		newUser := models.User{
			Username:     username,
			PasswordHash: string(hashedPassword),
		}

		var users []models.User

		fileData, err := os.ReadFile("users.json")
		if err == nil {
			json.Unmarshal(fileData, &users)
		}

		users = append(users, newUser)

		jsonData, _ := json.MarshalIndent(users, "", "  ")

		os.WriteFile("users.json", jsonData, 0644)

		fmt.Printf("Hashed Password: %s\n", string(hashedPassword))

		return true
	}

	return false
}
