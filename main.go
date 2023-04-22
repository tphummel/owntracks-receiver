package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type LocationUpdate struct {
	Type      string  `json:"_type"`
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
	Timestamp int64   `json:"tst"`
}

func initDB(filename string) *sql.DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS location_updates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		longitude REAL,
		latitude REAL,
		timestamp INTEGER
	)`)

	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	return db
}

func saveLocationUpdate(db *sql.DB, locationUpdate *LocationUpdate) error {
	stmt, err := db.Prepare("INSERT INTO location_updates (type, longitude, latitude, timestamp) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(locationUpdate.Type, locationUpdate.Longitude, locationUpdate.Latitude, locationUpdate.Timestamp)
	return err
}

func handleLocationUpdate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
		return
	}

	var locationUpdate LocationUpdate
	err := json.NewDecoder(r.Body).Decode(&locationUpdate)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}

	err = saveLocationUpdate(db, &locationUpdate)
	if err != nil {
		http.Error(w, "Failed to save location update", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Location update received and saved: %+v\n", locationUpdate)
	w.WriteHeader(http.StatusOK)
}

func main() {
	db := initDB("/app/data/owntracks.db")
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleLocationUpdate(db, w, r)
	})

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
