package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Golang-energetics-collection/models"
	"github.com/Golang-energetics-collection/utils"
)

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
	if newUser.Email != "" {
		var user1 models.User
		findError := models.DB.Where(fmt.Sprintf("email = '%s'", newUser.Email)).Find(&user1).Error

		if findError == nil {
			answer := Message{Status: "404", Message: "User with such email exists"}
			json.NewEncoder(w).Encode(answer)
			return
		} else {
			existingUser.Email = newUser.Email
		}
	}

	models.DB.Save(existingUser)
	json.NewEncoder(w).Encode(existingUser)
}
