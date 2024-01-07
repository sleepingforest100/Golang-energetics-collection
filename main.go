package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Energetic struct {
	EnergeticsID       uint   `gorm:"primaryKey"`
	Name               string `gorm:"not null"`
	Taste              string
	Description        string `gorm:"size:128"`
	ManufacturerName   string `gorm:"size:35"`
	ManufactureCountry string `gorm:"size:35"`
	PictureURL         string
	Composition        Composition `gorm:"foreignKey:EnergeticsID"`
}

type Composition struct {
	CompositionID uint `gorm:"primaryKey"`
	EnergeticsID  uint `gorm:"index"`
	Caffeine      uint
	Taurine       uint
}

type CompositionUniqueConstraint struct {
	CompositionID uint `gorm:"uniqueIndex:idx_composition_energetics"`
}

var energeticsList []Energetic

func main() {

	dsn := "host=localhost user=postgres password=222316pb dbname=energetix port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Energetic{}, &Composition{})

	db.Preload("Composition").Find(&energeticsList)

	// createEnergetic(db, "Monster2", "Orange", "Wow", "Ho", "Hungry", "dfsfsfs", 32, 400)

	for _, energetics := range energeticsList {
		log.Printf("EnergeticsID: %d, Name: %s, Taste: %s, Description: %s, ManufacturerName: %s, ManufactureCountry: %s, Composition: %+v\n",
			energetics.EnergeticsID, energetics.Name, energetics.Taste, energetics.Description, energetics.ManufacturerName, energetics.ManufactureCountry, energetics.Composition)
	}

	var energetic1 Energetic
	db.Preload("Composition").First(&energetic1, 1)
	log.Println(energetic1)

	mux := http.NewServeMux()
	mux.HandleFunc("/energetics", myHandler)
	fmt.Print("server starts... port 8080")
	http.ListenAndServe(":8080", mux)

}

func createEnergetic(db *gorm.DB,
	name, taste, description, manufacturerName, manufactureCountry, pictureURL string, caffeine, taurine uint) error {
	newEnergetic := Energetic{
		Name:               name,
		Taste:              taste,
		Description:        description,
		ManufacturerName:   manufacturerName,
		ManufactureCountry: manufactureCountry,
		PictureURL:         pictureURL,
		Composition: Composition{
			Caffeine: caffeine,
			Taurine:  taurine,
		},
	}
	result := db.Create(&newEnergetic)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getEnergetics(w, r)
	case http.MethodPost:
		postEnergetic(w, r)
	default:
		http.Error(w, "Invalid http method", http.StatusMethodNotAllowed)
	}
}

func getEnergetics(w http.ResponseWriter, r *http.Request) {
	responseJSON, err := json.Marshal(energeticsList)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Write(responseJSON)
}

func postEnergetic(w http.ResponseWriter, r *http.Request) {

	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil || message.Message == "" {
		answer := Message{Status: "400", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	fmt.Println(message.Message)

	fmt.Printf("Message: %v", message.Message)
	answer := Message{"Success", "Data successfully received"}
	json.NewEncoder(w).Encode(answer)

}
