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
	Type       string    `json:"_type"`
	Acc        *int      `json:"acc,omitempty"`
	Alt        *int      `json:"alt,omitempty"`
	Batt       *int      `json:"batt,omitempty"`
	BS         *int      `json:"bs,omitempty"`
	Cog        *int      `json:"cog,omitempty"`
	Latitude   float64   `json:"lat"`
	Longitude  float64   `json:"lon"`
	Rad        *int      `json:"rad,omitempty"`
	T          string    `json:"t,omitempty"`
	TID        string    `json:"tid,omitempty"`
	Timestamp  int64     `json:"tst"`
	Vac        *int      `json:"vac,omitempty"`
	Vel        *int      `json:"vel,omitempty"`
	P          *float64  `json:"p,omitempty"`
	POI        *string   `json:"poi,omitempty"`
	Conn       *string   `json:"conn,omitempty"`
	Tag        *string   `json:"tag,omitempty"`
	Topic      string    `json:"topic"`
	InRegions  *[]string `json:"inregions,omitempty"`
	InRIDs     *[]string `json:"inrids,omitempty"`
	SSID       *string   `json:"SSID,omitempty"`
	BSSID      *string   `json:"BSSID,omitempty"`
	CreatedAt  *string   `json:"created_at,omitempty"`
	Monitoring *int      `json:"m,omitempty"`
}

func initDB(filename string) *sql.DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS location_updates (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT,
		acc INTEGER,
		alt INTEGER,
		batt INTEGER,
		bs INTEGER,
		cog INTEGER,
		latitude REAL,
		longitude REAL,
		rad INTEGER,
		t TEXT,
		tid TEXT,
		timestamp INTEGER,
		vac INTEGER,
		vel INTEGER,
		p REAL,
		poi TEXT,
		conn TEXT,
		tag TEXT,
		topic TEXT,
		inregions TEXT,
		inrids TEXT,
		ssid TEXT,
		bssid TEXT,
		created_at TEXT,
		monitoring INTEGER
	)`)

	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	return db
}

func saveLocationUpdate(db *sql.DB, locationUpdate *LocationUpdate) error {
	stmt, err := db.Prepare(`INSERT INTO location_updates (
		type, acc, alt, batt, bs, cog, latitude, longitude, rad, t, tid, timestamp, vac, vel, p, poi, conn, tag, topic, inregions, inrids, ssid, bssid, created_at, monitoring
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		locationUpdate.Type,
		locationUpdate.Acc,
		locationUpdate.Alt,
		locationUpdate.Batt,
		locationUpdate.BS,
		locationUpdate.Cog,
		locationUpdate.Latitude,
		locationUpdate.Longitude,
		locationUpdate.Rad,
		locationUpdate.T,
		locationUpdate.TID,
		locationUpdate.Timestamp,
		locationUpdate.Vac,
		locationUpdate.Vel,
		locationUpdate.P,
		locationUpdate.POI,
		locationUpdate.Conn,
		locationUpdate.Tag,
		locationUpdate.Topic,
		locationUpdate.InRegions,
		locationUpdate.InRIDs,
		locationUpdate.SSID,
		locationUpdate.BSSID,
		locationUpdate.CreatedAt,
		locationUpdate.Monitoring,
	)
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
	db := initDB("/home/opc/data/owntracks.db")
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
