package main

import (
	"log"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/Golang-energetics-collection/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func TestGetEnergetics(t *testing.T) {

	setupTestDatabase()

	request, _ := http.NewRequest("GET", "/energetix", nil)
	response := httptest.NewRecorder()
	getEnergetics(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Incorrect status code. Expected: %d, Got: %d", http.StatusOK, response.Code)
	}

}

func setupTestDatabase() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	err = models.DB.AutoMigrate()
	if err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}
}
