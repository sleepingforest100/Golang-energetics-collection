package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/Golang-energetics-collection/models"
	"github.com/Golang-energetics-collection/utils"
	"github.com/google/uuid"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var emailAuth smtp.Auth

var smtpHost string
var email string
var password string

func InitEmailAuth() {

	if err := godotenv.Load(); err != nil {
		logrus.Print("No .env file found")
	}

	smtpHost = os.Getenv("SMTP_HOST")
	email = os.Getenv("SMTP_EMAIL")
	password = os.Getenv("EMAIL_PASSWORD")

	emailAuth = smtp.PlainAuth("", email, password, smtpHost)

	msg := []byte("To: 222316@astanait.edu.kz\r\n" +
		"Subject: Test" +
		"\r\n" +
		"Test test test")

	err := smtp.SendMail(smtpHost+":587", emailAuth, "mspolinko@gmail.com", []string{"222316@astanait.edu.kz"}, msg)
	if err != nil {
		log.Fatal(err)
	}

}

type Token struct {
	Token string `json:"token"`
}

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Limit   string `json:"limit"`
}

var jwtKey = []byte("my_secret_key")

func Login(w http.ResponseWriter, r *http.Request) {

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "Incorrect input"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var existingUser models.User

	models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "There is no such user"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "invalid password"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	if existingUser.Confirmed == false {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "Prohibited. Please, confirm your email!"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &models.Claims{
		Role:  existingUser.Role,
		Email: existingUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		answer := Message{Status: "500", Message: "could not generate token"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	answer := Token{Token: tokenString}
	json.NewEncoder(w).Encode(answer)
}

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "Incorrect input"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var existingUser models.User

	models.DB.Where("email = ?", user.Email).First(&existingUser)

	if existingUser.ID != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		answer := Message{Status: "400", Message: "Such user already exists"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		w.WriteHeader(http.StatusInternalServerError)
		answer := Message{Status: "500", Message: "could not generate password hash"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	user.Confirmed = false
	user.Role = "user"

	token := generateConfirmationToken()
	user.ConfirmToken = token

	models.DB.Create(&user)

	sendConfirmationToken(token, user.Email)

	w.WriteHeader(http.StatusOK)
	answer := Message{Status: "200", Message: "User created. Please, confirm your email"}
	json.NewEncoder(w).Encode(answer)
}

func Home(w http.ResponseWriter, r *http.Request) {

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

	if claims.Role != "user" && claims.Role != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "unauthorized"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	w.WriteHeader(http.StatusOK)
	answer := Message{Status: "200", Message: "success"}
	json.NewEncoder(w).Encode(answer)
	return
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
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

	if claims.Role != "user" && claims.Role != "admin" {
		w.WriteHeader(http.StatusUnauthorized)
		answer := Message{Status: "401", Message: "unauthorized"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "Incorrect input"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var existingUser models.User

	models.DB.Where("email = ?", claims.Email).First(&existingUser)

	if existingUser.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "There is no such user"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		w.WriteHeader(http.StatusInternalServerError)
		answer := Message{Status: "500", Message: "could not generate password hash"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	models.DB.Model(&existingUser).Update("password", user.Password)

	w.WriteHeader(http.StatusOK)
	answer := Message{Status: "200", Message: "Password updated"}
	json.NewEncoder(w).Encode(answer)
}

func ConfirmEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	email := r.URL.Query().Get("email")

	if token == "" {
		http.Error(w, "Missing token parameter", http.StatusBadRequest)
		return
	}
	if email == "" {
		http.Error(w, "Missing email parameter", http.StatusBadRequest)
		return
	}

	var user models.User
	models.DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "Invalid token"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	if user.ConfirmToken != token {
		w.WriteHeader(http.StatusBadRequest)
		answer := Message{Status: "400", Message: "Invalid token"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	user.Confirmed = true
	models.DB.Save(&user)

	w.WriteHeader(http.StatusOK)
	answer := Message{Status: "200", Message: "Thank you for confirming your email!"}
	json.NewEncoder(w).Encode(answer)

}

func generateConfirmationToken() string {
	uuidObj := uuid.New()
	timestamp := time.Now().Unix()
	token := fmt.Sprintf("%s-%d", uuidObj.String(), timestamp)
	return token
}

func sendConfirmationToken(token string, email string) {

	body := "Click the following link to confirm your email: http://localhost:8080/confirm?token=" + token + "&email=" + email

	msg := []byte("To: " + email + "\r\n" +
		"Subject: Confirm your email \r\n" +
		"\r\n" +
		body)

	err := smtp.SendMail(smtpHost+":587", emailAuth, "mspolinko@gmail.com", []string{email}, msg)
	if err != nil {
		log.Fatal(err)
	}

}
