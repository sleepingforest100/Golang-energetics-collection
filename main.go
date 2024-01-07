package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	router.HandleFunc("/energetix", getEnergetics).Methods("GET")
	// myRouter.HandleFunc("/energetics", postEnergetics).Methods("POST")
	router.HandleFunc("/energetix/{id}", getEnergeticsById).Methods("GET")

	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:63342"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	credentials := handlers.AllowCredentials()
	http.Handle("/", handlers.CORS(headers, origins, methods, credentials)(router))
	erro := http.ListenAndServe(":8080", nil)
	if erro != nil {
		panic(err)
	}

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

func handleRequests(db *gorm.DB) {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/energetics", getEnergetics).Methods("GET")
	// myRouter.HandleFunc("/energetics", postEnergetics).Methods("POST")
	myRouter.HandleFunc("/energetics/{id}", getEnergeticsById).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func getEnergetics(w http.ResponseWriter, r *http.Request) {
	responseJSON, err := json.Marshal(energeticsList)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Write(responseJSON)
}

func getEnergeticsById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// log.Print(params)
	var energetic1 Energetic
	dsn := "host=localhost user=postgres password=222316pb dbname=energetix port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Energetic{}, &Composition{})
	err2 := db.Preload("Composition").First(&energetic1, params["id"]).Error
	sqlDB, err := db.DB()
	sqlDB.Close()
	if err2 != nil {
		answer := Message{Status: "404", Message: "Energy drink with such ID does not exist"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	json.NewEncoder(w).Encode(energetic1)
	return
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
