package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"testing"
)

func TestLogin(t *testing.T) {

	requestBody, err := json.Marshal(map[string]string{
		"email":    "letun_igor007@mail.ru",
		"password": "321",
	})

	if err != nil {
		t.Errorf("Didn't manage to marshal json")
	}

	request, err2 := http.NewRequest("POST", "http://localhost:8080/auth/login", bytes.NewBuffer(requestBody))

	// fmt.Print(requestBody)

	if err2 != nil {
		t.Errorf("Didn't manage to send request")
	}

	response := httptest.NewRecorder()
	Login(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Incorrect status code. Expected: %d, Got: %d", http.StatusOK, response.Code)
	}

	var token Token

	if err := json.NewDecoder(response.Body).Decode(&token); err != nil {
		t.Errorf("Incorrect response body. Expected:token, Got: %s", response.Body.String())
	}

}
