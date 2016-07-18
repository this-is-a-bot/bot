package tracker

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

const trackerCatalogTableName = "tracker_catalog"
const trackerEventTableName = "tracker_events"

var (
	queryTrackingListByUser string
	queryTrackingEventByID  string
)

// Prepare queries.
func init() {
	fields := []string{
		"id", "name", "unit", "latest_event",
	}
	queryTrackingListByUser = fmt.Sprintf(
		"SELECT %s FROM %s WHERE disabled IS FALSE AND username = $1 AND app = $2",
		strings.Join(fields, ","), trackerCatalogTableName)

	queryTrackingEventByID = fmt.Sprintf(
		"SELECT value, marked_at FROM %s WHERE id = $1", trackerEventTableName)
}

type Catalog struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Unit  string  `json:"unit,omitempty"`
	Done  bool    `json:"done"`
	Value float32 `json:"value,omitempty"`
}

// Get a list of tracking catalogs, specifying the current status for each one.
func GetTrackingCatalogs(db *sql.DB, username string, app string) ([]Catalog, error) {
	rows, err := db.Query(queryTrackingListByUser, username, app)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Catalog, 0)
	for rows.Next() {
		var catalog Catalog
		var latestEventID int

		err = rows.Scan(&catalog.ID, &catalog.Name, &catalog.Unit, &latestEventID)
		if err != nil {
			return nil, err
		}

		// Fetch latest event to see whether this catalog has been completed.
		if latestEventID > 0 {
			var value float32
			var markedAt time.Time
			err = db.QueryRow(queryTrackingEventByID, latestEventID).Scan(&value, &markedAt)
			if err != nil {
				return nil, err
			}

			now := time.Now()
			y1, m1, d1 := now.Date()
			y2, m2, d2 := markedAt.Date()
			if y1 == y2 && m1 == m2 && d1 == d2 {
				// Already finished for today.
				catalog.Done = true
				catalog.Value = value
			}
		}

		res = append(res, catalog)
	}
	return res, nil
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
