package main

import (
	"os"
	"reflect"
	"strings"
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
		Type:       "location",
		Acc:        intPtr(101),
		Alt:        intPtr(500),
		Batt:       intPtr(85),
		BS:         intPtr(1),
		COG:        intPtr(180),
		Latitude:   12.34,
		Longitude:  56.78,
		Rad:        intPtr(10),
		T:          "t",
		TID:        "JJ",
		Timestamp:  1618859345,
		Vac:        intPtr(5),
		Vel:        intPtr(60),
		P:          float64Ptr(0.5),
		POI:        stringPtr("MyPoint"),
		Conn:       stringPtr("w"),
		Tag:        stringPtr("tag"),
		Topic:      "owntracks/user/device",
		InRegions:  []string{"region1", "region2"},
		InRIDs:     []string{"1", "2"},
		SSID:       stringPtr("MySSID"),
		BSSID:      stringPtr("00:11:22:33:44:55"),
		CreatedAt:  stringPtr("2023-05-10T04:00:00Z"),
		Monitoring: intPtr(1),
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
	var lat, lon, p float64
	var time int64
	var createdAt, loc_type, tid, topic string
	var acc, alt, batt, bs, cog, rad, vel, vac int
	var trigger, conn, ssid, bssid string
	var poi, tag string
	var monitoringMode int64
	var inRegions, inRids string

	found := false
	for rows.Next() {
		err := rows.Scan(&id, &loc_type, &acc, &alt, &batt, &bs, &cog, &lat, &lon, &rad, &trigger, &tid, &time, &vac, &vel, &p, &poi, &conn, &tag, &topic, &inRegions, &inRids, &ssid, &bssid, &createdAt, &monitoringMode)
		if err != nil {
			t.Errorf("Failed to scan row: %s", err)
		}

		inRegionsList := strings.Split(inRegions, ",")
		inRidsList := strings.Split(inRids, ",")

		if lat == locationUpdate.Latitude && lon == locationUpdate.Longitude && time == locationUpdate.Timestamp && tid == locationUpdate.TID && topic == locationUpdate.Topic && reflect.DeepEqual(inRegionsList, locationUpdate.InRegions) && reflect.DeepEqual(inRidsList, locationUpdate.InRIDs) {
			found = true
			break
		}
	}

	if !found {
		t.Error("Location update not found in the database")
	}

	os.Remove("test.db") // Clean up the test database file
}

// Helper functions to create pointers for basic types
func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}
