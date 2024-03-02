package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/Golang-energetics-collection/models"
	"github.com/Golang-energetics-collection/utils"
)

type UserRole struct {
	Role string `json:"role"`
}

type emailMessage struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "No Authorization header provided"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "unauthorized"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	if claims.Role != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "not allowed"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var users []models.User

	models.DB.Find(&users)

	json.NewEncoder(w).Encode(users)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "No Authorization header provided"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "unauthorized"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var user models.User

	findError := models.DB.Where(fmt.Sprintf("email = '%s'", claims.Email)).Find(&user).Error

	if findError != nil {
		answer := Message{Status: "404", Message: "User with such email does not exist"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "No Authorization header provided"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "unauthorized"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var existingUser models.User

	findError := models.DB.Where(fmt.Sprintf("email = '%s'", claims.Email)).Find(&existingUser).Error

	if findError != nil {
		answer := Message{Status: "404", Message: "User with such email does not exist"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "Incorrect input"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	if newUser.Name != "" {
		existingUser.Name = newUser.Name
	}
	if newUser.Address != "" {
		existingUser.Address = newUser.Address
	}

	// if newUser.Email != "" {
	// 	var user1 models.User
	// 	findError := models.DB.Where(fmt.Sprintf("email = '%s'", newUser.Email)).Find(&user1).Error

	// 	if findError == nil {
	// 		answer := Message{Status: "404", Message: "User with such email exists"}
	// 		json.NewEncoder(w).Encode(answer)
	// 		return
	// 	} else {
	// 		existingUser.Email = newUser.Email
	// 	}
	// }

	models.DB.Save(existingUser)
	json.NewEncoder(w).Encode(existingUser)
}

func ChangeRole(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "No Authorization header provided"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "unauthorized"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	if claims.Role != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "not allowed"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	params := mux.Vars(r)
	var user models.User
	err = models.DB.First(&user, params["id"]).Error

	if err != nil {
		answer := Message{Status: "404", Message: "User with such ID does not exist"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var userRole UserRole

	err = json.NewDecoder(r.Body).Decode(&userRole)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	role := userRole.Role
	if role == "admin" || role == "user" {
		user.Role = role
		models.DB.Save(&user)
		answer := Message{Status: "200", Message: "User role was changed"}
		json.NewEncoder(w).Encode(answer)
	} else {
		answer := Message{Status: "400", Message: "There is no such role"}
		json.NewEncoder(w).Encode(answer)
	}
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "No Authorization header provided"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ParseToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "unauthorized"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	if claims.Role != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "not allowed"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var emailMessage emailMessage
	err = json.NewDecoder(r.Body).Decode(&emailMessage)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	var users []models.User

	models.DB.Find(&users)

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	emailAuth = smtp.PlainAuth("", email, password, smtpHost)

	emailBody := "Subject: " + emailMessage.Subject +
		"\r\n" + emailMessage.Body

	for _, recipient := range users {
		msg := []byte("To: " + recipient.Email + "\r\n" + emailBody)

		err := smtp.SendMail(smtpHost+":587", emailAuth, "mspolinko@gmail.com", []string{recipient.Email}, msg)
		if err != nil {
			log.Fatal(err)
		}
	}

	answer := Message{Status: "200", Message: "Email was successfully sent"}
	json.NewEncoder(w).Encode(answer)
}
