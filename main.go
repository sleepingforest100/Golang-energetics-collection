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
	"golang.org/x/time/rate"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Limit   string `json:"limit"`
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

var limiter = rate.NewLimiter(10, 10)

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
	dsn := "host=localhost user=postgres password=pgadmin4 dbname=energetix port=5432 sslmode=disable TimeZone=Asia/Shanghai"
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
	logrus.Info("Automigration for energetics and compositions")
}

func main() {
	initLog()
	initDB()
	db.Preload("Composition").Find(&energeticsList)
	logrus.Info("Preload energetics collection")

	router := mux.NewRouter()

	router.Handle("/energetix", RateLimitMiddleware(http.HandlerFunc(getEnergetics))).Methods("GET")
	router.Handle("/energetix", RateLimitMiddleware(http.HandlerFunc(postEnergetic))).Methods("POST")
	router.Handle("/energetix/{id}", RateLimitMiddleware(http.HandlerFunc(getEnergeticsById))).Methods("GET")
	router.Handle("/energetix/{id}", RateLimitMiddleware(http.HandlerFunc(updateEnergeticsById))).Methods("PUT")
	router.Handle("/energetix/{id}", RateLimitMiddleware(http.HandlerFunc(deleteEnergeticById))).Methods("DELETE")
	router.Handle("/pages", RateLimitMiddleware(http.HandlerFunc(getNumberOfPages))).Methods("GET")

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "index-go.html", http.StatusSeeOther)
	})
	router.HandleFunc("/index-go.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index-go.html")
	})
	router.HandleFunc("/form-go.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "form-go.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	headers := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:63342", "http://127.0.0.1:5500"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	credentials := handlers.AllowCredentials()
	http.Handle("/", handlers.CORS(headers, origins, methods, credentials)(router))
	erro := http.ListenAndServe(":8080", nil)
	logrus.WithFields(logrus.Fields{
		"module":   "main",
		"function": "main",
		"action":   "servers starts",
		"port":     "8080",
	}).Info("Server launches")
	if erro != nil {
		logrus.Panic("Server did not run")
		panic(erro)
	}
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			message := Message{
				Status:  "Request Failed",
				Message: "Too many requests in one time. Try again a bit later",
				Limit:   "2 request per second with a burst of 3 requests",
			}
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(&message)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getEnergetics(w http.ResponseWriter, r *http.Request) {

	if !limiter.Allow() {
		message := Message{
			Status:  "Request Failed",
			Message: "Too many requests in one time. Try again a bit later",
			Limit:   "1 request per second with a burst of 3 requests",
		}
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(&message)
		return
	}

	db.AutoMigrate(&Energetic{}, &Composition{})
	logrus.Info("Automigration for energetics and compositions")

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

	logrus.WithFields(logrus.Fields{
		"module":   "main",
		"function": "getEnergetics",
		"action":   "query reading",
		"queryParams": logrus.Fields{
			"sort":                sort,
			"order":               order,
			"taurine_gte":         taurine_gte,
			"taurine_lte":         taurine_lte,
			"caffeine_gte":        caffeine_gte,
			"caffeine_lte":        caffeine_lte,
			"taste":               taste,
			"energeticName":       nameEn,
			"manufacturerName":    manufacturerName,
			"manufacturerCountry": manufacturerCountry,
		},
	}).Info("Initializating of parametrs from the query")

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

	logrus.WithFields(logrus.Fields{
		"module":   "main",
		"function": "getEnergetics",
		"action":   "query params updating",
		"queryParams": logrus.Fields{
			"sort":                sort,
			"order":               order,
			"taurine_gte":         taurine_gte,
			"taurine_lte":         taurine_lte,
			"caffeine_gte":        caffeine_gte,
			"caffeine_lte":        caffeine_lte,
			"taste":               taste,
			"energeticName":       nameEn,
			"manufacturerName":    manufacturerName,
			"manufacturerCountry": manufacturerCountry,
		},
	}).Info("Fixing parametrs from the query")

	pageInt, errConv := strconv.Atoi(page)

	logrus.WithFields(logrus.Fields{
		"function": "getEnergetics",
		"action":   "converting str  to int",
		"page":     page,
		"pageInt":  pageInt,
	}).Info("Attempt to convert param 'page' to int")

	if errConv != nil {
		http.Error(w, "Failed to parse page number", http.StatusInternalServerError)
		logrus.Error("Failed to convert page to int")
		return
	}

	offset := (pageInt - 1) * limit
	logrus.WithFields(logrus.Fields{
		"function": "getEnergetics",
		"offset":   offset,
	}).Info("Computing offset")

	var totalCount int64

	erro := db.
		Model(&Energetic{}).
		Joins("Composition").
		Order(sort+" "+order).
		Where("name ILIKE ? AND taste ILIKE ? AND taurine >= ? AND taurine <= ? AND caffeine >= ? AND caffeine <= ? AND manufacturer_name ILIKE ? AND manufacture_country ILIKE ?",
			"%"+nameEn+"%", "%"+taste+"%", taurine_gte, taurine_lte, caffeine_gte, caffeine_lte,
			"%"+manufacturerName+"%", "%"+manufacturerCountry+"%").
		Count(&totalCount).
		Error

	if erro != nil {
		http.Error(w, "Failed to marshal JSON with sorting (counting)", http.StatusInternalServerError)
		logrus.Error("Couldn't execute query")
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

	logrus.WithFields(logrus.Fields{
		"function": "getEnergetics",
		"action":   "search energetics in db using filters, sorting and pages",
		"sorting": logrus.Fields{
			"sort":  sort,
			"order": order,
		},
		"pagination": logrus.Fields{
			"page":   pageInt,
			"limit":  limit,
			"offset": offset,
		},
		"filters": logrus.Fields{
			"taurine_gte":         taurine_gte,
			"taurine_lte":         taurine_lte,
			"caffeine_gte":        caffeine_gte,
			"caffeine_lte":        caffeine_lte,
			"taste":               taste,
			"energeticName":       nameEn,
			"manufacturerName":    manufacturerName,
			"manufacturerCountry": manufacturerCountry,
		},
	}).Info("DB.FIND()")

	if err != nil {
		http.Error(w, "Failed to marshal JSON with sorting", http.StatusInternalServerError)
		logrus.Error("Couldn't execute query")
		return

	}

	totalCount = int64(math.Ceil(float64(totalCount) / float64(limit)))
	response := struct {
		TotalCount int64       `json:"total_count" `
		Data       []Energetic `json:"data"`
	}{
		TotalCount: totalCount,
		Data:       energeticsList,
	}

	responseJSON, err := json.Marshal(response)
	logrus.Info("Sending energetics list json as a response")
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		logrus.Error("Couldn't marshal JSON")
		return
	}

	w.Write(responseJSON)
}

func getEnergeticsById(w http.ResponseWriter, r *http.Request) {

	if !limiter.Allow() {
		message := Message{
			Status:  "Request Failed",
			Message: "Too many requests in one time. Try again a bit later",
			Limit:   "1 request per second with a burst of 3 requests",
		}
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(&message)
		return
	}

	params := mux.Vars(r)
	var energetic1 Energetic
	db.AutoMigrate(&Energetic{}, &Composition{})
	err2 := db.Preload("Composition").First(&energetic1, params["id"]).Error

	logrus.WithFields(logrus.Fields{
		"function":      "getEnergeticsById",
		"action":        "getting 1 energetic by id",
		"energetics_id": params["id"],
	}).Info("Attemt to get energetics by ID " + params["id"] + " .db.First()")

	if err2 != nil {
		answer := Message{Status: "404", Message: "Energy drink with such ID does not exist"}
		json.NewEncoder(w).Encode(answer)
		logrus.Info("Energy drink with " + params["id"] + " was not found in DB")
		return
	}
	json.NewEncoder(w).Encode(energetic1)
	return
}

func updateEnergeticsById(w http.ResponseWriter, r *http.Request) {

	if !limiter.Allow() {
		message := Message{
			Status:  "Request Failed",
			Message: "Too many requests in one time. Try again a bit later",
			Limit:   "1 request per second with a burst of 3 requests",
		}
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(&message)
		return
	}

	logrus.WithFields(logrus.Fields{
		"function": "updateEnergeticsById",
		"action":   "launching of updating 1 energetic by id",
	}).Info("updateEnergeticsById() starts")

	params := mux.Vars(r)
	targetID, err := strconv.ParseUint(params["id"], 10, 64)
	logrus.WithFields(logrus.Fields{
		"function":      "updateEnergeticsById",
		"action":        "convert id to int",
		"energetics_id": targetID,
	}).Info("parsed ID for updating")

	if err != nil {
		answer := Message{Status: "400", Message: "Incorrect id"}
		json.NewEncoder(w).Encode(answer)
		logrus.Error("Error with parsing id to int")
		return
	}

	var updatedEnergetic Energetic

	if err := json.NewDecoder(r.Body).Decode(&updatedEnergetic); err != nil {
		answer := Message{Status: "404", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		logrus.WithFields(logrus.Fields{
			"function":    "updateEnergeticsById",
			"action":      "reading attributes for update",
			"requestBody": r.Body,
		}).Error("Incorrect energetic fields")
		return
	}

	logrus.WithFields(logrus.Fields{
		"function":    "updateEnergeticsById",
		"action":      "reading attributes for update",
		"requestBody": r.Body,
	}).Info("Correct fields for updating were accepted")

	db.AutoMigrate(&Energetic{}, &Composition{})
	var existingEnergetic Energetic

	if err := db.Preload("Composition").First(&existingEnergetic, targetID).Error; err != nil {
		answer := Message{Status: "404", Message: "Energy drink with such ID does not exist"}
		json.NewEncoder(w).Encode(answer)
		logrus.WithFields(logrus.Fields{
			"function": "updateEnergeticsById",
			"action":   "check if such energetic exists",
			"targetID": targetID,
		}).Error("Energetic with such ID cannot be found")
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

	logrus.WithFields(logrus.Fields{
		"function": "updateEnergeticsById",
		"action":   "updating all fields",
		"targetID": targetID,
		"changes": logrus.Fields{
			"name":                existingEnergetic.Name,
			"taste":               existingEnergetic.Taste,
			"description":         existingEnergetic.Description,
			"manufacturerName":    existingEnergetic.ManufactureCountry,
			"manufactureCountry ": existingEnergetic.ManufactureCountry,
			"pictureURL":          existingEnergetic.PictureURL,
			"caffeine":            existingEnergetic.Composition.Caffeine,
			"taurine":             existingEnergetic.Composition.Taurine,
		},
	}).Error("The energetic was updated")

	w.WriteHeader(http.StatusOK)
	answer := Message{Status: "200", Message: "Energy drink was updated"}
	json.NewEncoder(w).Encode(answer)
	return
}

func postEnergetic(w http.ResponseWriter, r *http.Request) {

	if !limiter.Allow() {
		message := Message{
			Status:  "Request Failed",
			Message: "Too many requests in one time. Try again a bit later",
			Limit:   "1 request per second with a burst of 3 requests",
		}
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(&message)
		return
	}

	var newEnergetic Energetic

	logrus.WithFields(logrus.Fields{
		"function": "postEnergetic",
	}).Info("Post energetic method runs")

	if err := json.NewDecoder(r.Body).Decode(&newEnergetic); err != nil {
		answer := Message{Status: "404", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		logrus.WithFields(logrus.Fields{
			"function": "postEnergetic",
			"action":   "attempt to decoding json",
		}).Error("Invalid json message was recieved")
		return
	}

	db.AutoMigrate(&Energetic{}, &Composition{})

	if err := db.Create(&newEnergetic).Error; err != nil {
		answer := Message{Status: "404", Message: "Invalid JSON message"}
		json.NewEncoder(w).Encode(answer)
		logrus.WithFields(logrus.Fields{
			"function": "postEnergetic",
			"action":   "creating new energetic",
		}).Error("Json message fields are not appropriate for energetics collection")
		return
	}

	logrus.WithFields(logrus.Fields{
		"function": "postEnergetic",
		"action":   "creating new energetic",
	}).Info("Sucessful create of energetic")

	db.Preload("Composition").Find(&energeticsList)
	logrus.Info("Updating energetics List by preloading and find")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEnergetic)
}

func deleteEnergeticById(w http.ResponseWriter, r *http.Request) {

	if !limiter.Allow() {
		message := Message{
			Status:  "Request Failed",
			Message: "Too many requests in one time. Try again a bit later",
			Limit:   "1 request per second with a burst of 3 requests",
		}
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(&message)
		return
	}

	logrus.WithFields(logrus.Fields{
		"function": "deleteEnergeticById",
	}).Info("Delete energetic method runs")

	params := mux.Vars(r)
	targetID, err := strconv.ParseUint(params["id"], 10, 64)

	logrus.WithFields(logrus.Fields{
		"function":      "deleteEnergeticById",
		"action":        "convert id to int",
		"energetics_id": targetID,
	}).Info("parsed ID for deleting")

	if err != nil {
		answer := Message{Status: "400", Message: "Incorrect id"}
		json.NewEncoder(w).Encode(answer)
		logrus.Error("Error with parsing id to int")
		return
	}
	db.AutoMigrate(&Energetic{}, &Composition{})

	logrus.WithFields(logrus.Fields{
		"function":      "deleteEnergeticById",
		"action":        "attempt to delete energetic",
		"energetics_id": targetID,
	}).Info("Deleting energetic by ID")

	if err := db.Delete(&Energetic{}, targetID).Error; err != nil {
		answer := Message{Status: "404", Message: "Invalid id"}
		json.NewEncoder(w).Encode(answer)
		logrus.Error("Deleting unexisting energetic")
		return
	}
	db.Preload("Composition").Find(&energeticsList)

	answer := Message{Status: "410", Message: "Energy drink was deleted successfully"}
	logrus.Info("Energetic was deleted")
	json.NewEncoder(w).Encode(answer)
	w.WriteHeader(http.StatusOK)
}

func getNumberOfPages(w http.ResponseWriter, r *http.Request) {

	if !limiter.Allow() {
		message := Message{
			Status:  "Request Failed",
			Message: "Too many requests in one time. Try again a bit later",
			Limit:   "1 request per second with a burst of 3 requests",
		}
		w.WriteHeader(http.StatusTooManyRequests)
		json.NewEncoder(w).Encode(&message)
		return
	}

	db.Find(&energeticsList)
	count := int(math.Ceil(float64(len(energeticsList)) / float64(limit)))
	number := pagesCount{Pages: count}

	logrus.WithFields(logrus.Fields{
		"function":      "getNumberOfPages",
		"action":        "counting number of pages",
		"numberOfPages": count,
	}).Info("Count number of pages")

	json.NewEncoder(w).Encode(number)
}
