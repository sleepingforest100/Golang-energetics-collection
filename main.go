package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Message struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	// mux := http.NewServeMux()
	// mux.HandleFunc("/", myHandler)
	// fmt.Print("server starts... port 8080")
	// http.ListenAndServe(":8080", mux)

	connStr := "host=localhost port=5432 user=postgres password=222316pb dbname=energetix sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database")
	defer db.Close()

}

func myHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getMessage(w, r)
	case http.MethodPost:
		postMessage(w, r)
	default:
		http.Error(w, "Invalid http method", http.StatusMethodNotAllowed)
	}
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello!")
}

func postMessage(w http.ResponseWriter, r *http.Request) {

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
