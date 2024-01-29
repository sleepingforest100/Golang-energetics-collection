package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=222316pb dbname=energetix port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&Energetic{}, &Composition{})
}

func main() {

	initDB()
	db.Preload("Composition").Find(&energeticsList)

	router := mux.NewRouter()
	router.HandleFunc("/energetix", getEnergetics).Methods("GET")
	router.HandleFunc("/energetix", postEnergetic).Methods("POST")
	router.HandleFunc("/energetix/{id}", getEnergeticsById).Methods("GET")
	router.HandleFunc("/energetix/{id}", updateEnergeticsById).Methods("PUT")
	router.HandleFunc("/energetix/{id}", deleteEnergeticById).Methods("DELETE")
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "gofront/index-go.html")
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:63342", "http://127.0.0.1:5500"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	credentials := handlers.AllowCredentials()
	http.Handle("/", handlers.CORS(headers, origins, methods, credentials)(router))
	erro := http.ListenAndServe(":8080", nil)
	if erro != nil {
		panic(erro)
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

func getEnergetics(w http.ResponseWriter, r *http.Request) {
	db.AutoMigrate(&Energetic{}, &Composition{})

	sort := r.FormValue("sort")
	order := r.FormValue("order")

	// if len(sort) > 0 && len(order) > 0 {
	if len(sort) < 1 && len(order) < 1 {
		sort = "energetics_id"
		order = "desc"
	}

	err := db.
		Model(&Energetic{}).
		Joins("Composition").
		Order(sort + " " + order).
		Find(&energeticsList).
		Error

	if err != nil {
		http.Error(w, "Failed to marshal JSON with sorting", http.StatusInternalServerError)
		return
	}
	// }

	// db.Table("compositions").
	// 	Select("compositions.*, energetics.*").
	// 	Joins("INNER JOIN energetics ON energetics.composition_id = compositions.energetics_id"). // After this you can order by any field
	// 	Order("taurine asc").Find(&energeticsList)

	// err := db.
	// 	Model(&Energetic{}).
	// 	Joins("Composition").
	// 	Order("taurine DESC").
	// 	Find(&energeticsList).
	// 	Error

	// db.Preload("Composition").Find(&energeticsList)

	responseJSON, err := json.Marshal(energeticsList)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Write(responseJSON)
}

func getEnergeticsById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var energetic1 Energetic

	db.AutoMigrate(&Energetic{}, &Composition{})
	err2 := db.Preload("Composition").First(&energetic1, params["id"]).Error

	if err2 != nil {
		answer := Message{Status: "404", Message: "Energy drink with such ID does not exist"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	json.NewEncoder(w).Encode(energetic1)
	return
}

func updateEnergeticsById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	targetID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		answer := Message{Status: "400", Message: "Incorrect id"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	var updatedEnergetic Energetic
	if err := json.NewDecoder(r.Body).Decode(&updatedEnergetic); err != nil {
		answer := Message{Status: "404", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	db.AutoMigrate(&Energetic{}, &Composition{})

	var existingEnergetic Energetic
	if err := db.First(&existingEnergetic, targetID).Error; err != nil {
		answer := Message{Status: "404", Message: "Energy drink with such ID does not exist"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	existingEnergetic.Composition = updatedEnergetic.Composition

	db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&existingEnergetic).Updates(&updatedEnergetic)

	db.Preload("Composition").Find(&energeticsList)

	w.WriteHeader(http.StatusOK)
	answer := Message{Status: "200", Message: "Energy drink was updated"}
	json.NewEncoder(w).Encode(answer)
	return
}

func postEnergetic(w http.ResponseWriter, r *http.Request) {
	var newEnergetic Energetic
	if err := json.NewDecoder(r.Body).Decode(&newEnergetic); err != nil {
		answer := Message{Status: "404", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	db.AutoMigrate(&Energetic{}, &Composition{})

	if err := db.Create(&newEnergetic).Error; err != nil {
		answer := Message{Status: "404", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	db.Preload("Composition").Find(&energeticsList)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEnergetic)
}

func deleteEnergeticById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	targetID, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		answer := Message{Status: "400", Message: "Incorrect id"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	db.AutoMigrate(&Energetic{}, &Composition{})

	if err := db.Delete(&Energetic{}, targetID).Error; err != nil {
		answer := Message{Status: "404", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		return
	}
	db.Preload("Composition").Find(&energeticsList)

	answer := Message{Status: "410", Message: "Energy drink was deleted successfully"}
	json.NewEncoder(w).Encode(answer)
	w.WriteHeader(http.StatusOK)
}
