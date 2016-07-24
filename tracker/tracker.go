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
	queryTrackingListByUser              string
	queryTrackingEventByID               string
	insertTrackingCatalog                string
	updateTrackingCatalogWithLatestEvent string
	insertTrackingEvent                  string

	// TODO: Hardcoded timezone for now (Pacific time).
	pacific *time.Location
)

// Prepare queries.
func init() {
	fields := []string{
		"id", "name", "unit", "latest_event",
	}
	queryTrackingListByUser = fmt.Sprintf(
		"SELECT %s FROM %s WHERE disabled IS FALSE AND username = $1 AND app = $2",
		strings.Join(fields, ","), trackerCatalogTableName)

	insertTrackingCatalog = fmt.Sprintf(
		"INSERT INTO %s (USERNAME, APP, NAME, UNIT) VALUES ($1, $2, $3, $4) RETURNING id",
		trackerCatalogTableName)

	updateTrackingCatalogWithLatestEvent = fmt.Sprintf(
		"UPDATE %s SET latest_event = $1 WHERE id = $2", trackerCatalogTableName)

	queryTrackingEventByID = fmt.Sprintf(
		"SELECT value, marked_at FROM %s WHERE id = $1", trackerEventTableName)

	insertTrackingEvent = fmt.Sprintf(
		"INSERT INTO %s (catalog_id, value) VALUES ($1, $2) RETURNING id",
		trackerEventTableName)

	// Load timezone as Pacific time.
	var err error
	pacific, err = time.LoadLocation("US/Pacific")
	if err != nil {
		panic(err)
	}
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
		var latestEventID sql.NullInt64

		err = rows.Scan(&catalog.ID, &catalog.Name, &catalog.Unit, &latestEventID)
		if err != nil {
			return nil, err
		}

		// Fetch latest event to see whether this catalog has been completed.
		if latestEventID.Valid {
			var value float32
			var markedAt time.Time
			// TODO: join the table before to avoid this extra SQL query.
			err = db.QueryRow(
				queryTrackingEventByID, latestEventID.Int64).Scan(&value, &markedAt)
			if err != nil {
				return nil, err
			}

			// Convert to pacific time.
			now := time.Now().In(pacific)
			markedAt = markedAt.In(pacific)
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
func MarkDone(db *sql.DB, catalogID int, value float64) error {
	var eventID int64
	err := db.QueryRow(insertTrackingEvent, catalogID, value).Scan(&eventID)
	if err != nil {
		return err
	}

	_, err = db.Exec(updateTrackingCatalogWithLatestEvent, eventID, catalogID)
	if err != nil {
		return err
	}

	return nil
}

// Add a tracking item for the particular user.
func AddTracking(db *sql.DB, username string, app string, name string, unit string) (int64, error) {
	var id int64
	err := db.QueryRow(insertTrackingCatalog, username, app, name, unit).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Modify the tracking catalog with new name / unit.
func UpdateTracking(db *sql.DB, catalogID int, newName string, newUnit string) (int, error) {
	return 0, errors.New("not implemented")
}

// Delete the tracking item.
func RemoveTracking(db *sql.DB, catalogID int) error {
	return errors.New("not implemented")
}
