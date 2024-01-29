package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type pagesCount struct {
	Pages int `json:"pages`
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

var limit = 3

var logFile = "log.json"

func initLog() {
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile" + logFile)
		panic(err)
	}
	logrus.SetOutput(f)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.WithFields(logrus.Fields{
		"module":          "main",
		"function":        "initLog",
		"action":          "log file opening",
		"logrusFormatter": "JSONformatter",
		"logFile":         logFile,
	}).Info("Log file was opened")
}

func initDB() {
	dsn := "host=localhost user=postgres password=222316pb dbname=energetix port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var errConn error
	db, errConn = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	logrus.WithFields(logrus.Fields{
		"module":           "main",
		"function":         "initDB",
		"action":           "db connection",
		"connectionString": dsn,
	}).Info("Attempt to connect to db")
	if errConn != nil {
		logrus.Fatal("Failed to connect to the db")
	}

	db.AutoMigrate(&Energetic{}, &Composition{})
}

func main() {
	initLog()
	initDB()
	db.Preload("Composition").Find(&energeticsList)

	router := mux.NewRouter()
	router.HandleFunc("/energetix", getEnergetics).Methods("GET")
	router.HandleFunc("/energetix", postEnergetic).Methods("POST")
	router.HandleFunc("/energetix/{id}", getEnergeticsById).Methods("GET")
	router.HandleFunc("/energetix/{id}", updateEnergeticsById).Methods("PUT")
	router.HandleFunc("/energetix/{id}", deleteEnergeticById).Methods("DELETE")
	router.HandleFunc("/pages", getNumberOfPages).Methods("GET")
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

	taurine_gte := r.FormValue("taurine_gte")
	taurine_lte := r.FormValue("taurine_lte")
	caffeine_gte := r.FormValue("caffeine_gte")
	caffeine_lte := r.FormValue("caffeine_lte")
	taste := r.FormValue("taste")
	nameEn := r.FormValue("name")
	manufacturerName := r.FormValue("manufacturer")
	manufacturerCountry := r.FormValue("country")

	page := r.FormValue("page")

	if len(taurine_gte) < 1 {
		taurine_gte = "0"
	}
	if len(taurine_lte) < 1 {
		taurine_lte = "500000"
	}
	if len(caffeine_gte) < 1 {
		caffeine_gte = "0"
	}
	if len(caffeine_lte) < 1 {
		caffeine_lte = "500000"
	}
	if len(taste) < 1 {
		taste = ""
	}
	if len(nameEn) < 1 {
		nameEn = ""
	}
	if len(manufacturerName) < 1 {
		manufacturerName = ""
	}
	if len(manufacturerCountry) < 1 {
		manufacturerCountry = ""
	}

	if len(sort) < 1 && len(order) < 1 {
		sort = "energetics_id"
		order = "desc"
	}

	if len(page) < 1 {
		page = "1"
	}

	pageInt, errConv := strconv.Atoi(page)
	offset := (pageInt - 1) * limit
	print(offset)

	if errConv != nil {
		http.Error(w, "Failed to parse page number", http.StatusInternalServerError)
		return
	}

	err := db.
		Model(&Energetic{}).
		Joins("Composition").
		Limit(limit).
		Offset(offset).
		Order(sort+" "+order).
		Find(&energeticsList, "name ILIKE ? AND taste ILIKE ? AND taurine >= ? AND taurine <= ? AND caffeine >= ? AND caffeine <= ? AND manufacturer_name ILIKE ? AND manufacture_country ILIKE ?",
			"%"+nameEn+"%", "%"+taste+"%", taurine_gte, taurine_lte, caffeine_gte, caffeine_lte,
			"%"+manufacturerName+"%", "%"+manufacturerCountry+"%").
		Error

	if err != nil {
		http.Error(w, "Failed to marshal JSON with sorting", http.StatusInternalServerError)
		return
	}

	logrus.Info("Info message")

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
	if err := db.Preload("Composition").First(&existingEnergetic, targetID).Error; err != nil {
		answer := Message{Status: "404", Message: "Energy drink with such ID does not exist"}
		json.NewEncoder(w).Encode(answer)
		return
	}

	existingEnergetic.Name = updatedEnergetic.Name
	existingEnergetic.Taste = updatedEnergetic.Taste
	existingEnergetic.Description = updatedEnergetic.Description
	existingEnergetic.ManufacturerName = updatedEnergetic.ManufacturerName
	existingEnergetic.ManufactureCountry = updatedEnergetic.ManufactureCountry
	existingEnergetic.PictureURL = updatedEnergetic.PictureURL

	existingEnergetic.Composition.Caffeine = updatedEnergetic.Composition.Caffeine
	existingEnergetic.Composition.Taurine = updatedEnergetic.Composition.Taurine

	db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existingEnergetic)

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

func getNumberOfPages(w http.ResponseWriter, r *http.Request) {
	count := int(math.Ceil(float64(len(energeticsList)) / float64(limit)))
	number := pagesCount{Pages: count}
	json.NewEncoder(w).Encode(number)
}
