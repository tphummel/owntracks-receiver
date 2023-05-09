package main

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestInitDB tests the initDB function
func TestInitDB(t *testing.T) {
	os.Remove("test.db") // Remove test database file if it exists
	db := initDB("test.db")
	defer db.Close()

	_, err := db.Exec("SELECT * FROM location_updates")
	if err != nil {
		t.Errorf("Failed to query location_updates table: %s", err)
	}

	os.Remove("test.db") // Clean up the test database file
}

// TestSaveLocationUpdate tests the saveLocationUpdate function
func TestSaveLocationUpdate(t *testing.T) {
	os.Remove("test.db") // Remove test database file if it exists
	db := initDB("test.db")
	defer db.Close()

	locationUpdate := LocationUpdate{
		Latitude:  12.34,
		Longitude: 56.78,
		Type:      "location",
		Timestamp: 1618859345,
		TID:       "JJ",
		Topic:     "owntracks/user/device",
	}

	err := saveLocationUpdate(db, &locationUpdate)
	if err != nil {
		t.Errorf("Failed to save location update: %s", err)
	}

	rows, err := db.Query("SELECT * FROM location_updates")
	if err != nil {
		t.Errorf("Failed to query location_updates table: %s", err)
	}
	defer rows.Close()

	var id int
	var lat, lon float64
	var time int64
	var loc_type, tid, topic string
	var acc, alt, batt, bs, rad, vel, vac int
	var trigger, conn, ssid, bssid string
	var poi, tag string
	var createdAt, monitoringMode int64
	var inRegions, inRids, p float64

	found := false
	for rows.Next() {
		err := rows.Scan(&id, &loc_type, &acc, &alt, &batt, &bs, &lat, &lon, &rad, &trigger, &tid, &time, &vac, &vel, &p, &poi, &conn, &ssid, &bssid, &createdAt, &monitoringMode, &tag, &topic, &inRegions, &inRids)
		if err != nil {
			t.Errorf("Failed to scan row: %s", err)
		}

		if lat == locationUpdate.Latitude && lon == locationUpdate.Longitude && time == locationUpdate.Timestamp && tid == locationUpdate.TID && topic == locationUpdate.Topic {
			found = true
			break
		}
	}

	if !found {
		t.Error("Location update not found in the database")
	}

	os.Remove("test.db") // Clean up the test database file
}
