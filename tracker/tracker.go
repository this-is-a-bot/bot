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
	return nil, errors.New("not implemented")
}

// Mark done for a given catalog.
func markDone(db *sql.DB, catalogID int, value float32) error {
	// TODO:
	// 1. Insert into `tracker_events` table.
	// 2. Update `tracker_catalog` table to let the corresponding catalog point to
	//    latest event.
	return errors.New("not implemented")
}
