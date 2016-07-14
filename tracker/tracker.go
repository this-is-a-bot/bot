package tracker

import (
	"database/sql"
)

type Catalog struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Unit  string  `json:"unit,omitempty"`
	Done  bool    `json:"done"`
	Value float32 `json:"value,omitempty"`
}

// Get a list of tracking catalogs, specifying the current status.
func GetTrackingCatalogs(db *sql.DB, username string, app string) ([]Catalog, error) {
	// TODO:
	// 1. First fetch from `tracker_catalog` table,
	// 2. For each catalog, check whether latest event is for that specific day.
	// 3. If yes, fetch the latset event along with the value.
	return nil, errors.New("not implemented")
}

// Mark done for a given catalog (add an event to the catalog with timestamp).
func MarkDone(db *sql.DB, catalogID int, value float32) error {
	// TODO:
	// 1. Insert into `tracker_events` table.
	// 2. Update `tracker_catalog` table to let the corresponding catalog point to
	//    latest event.
	return errors.New("not implemented")
}

// Add a tracking item for the particular user.
func AddTracking(db *sql.DB, username string, app string, name string, unit string) (int, error) {
	// TODO: Insert into `tracker_catalog` and return the ID.
	return 0, errors.New("not implemented")
}

// Modify the tracking catalog with new name / unit.
func UpdateTracking(db *sql.DB, catalogID int, newName string, newUnit string) (int, error) {
	// TODO: Update `tracker_catalog` with new name / unit if not empty.
	return 0, errors.New("not implemented")
}

// Delete the tracking item.
func RemoveTracking(db *sql.DB, catalogID int) error {
	return errors.New("not implemented")
}
