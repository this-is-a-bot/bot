package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/this-is-a-bot/bot/redis"
	"github.com/this-is-a-bot/bot/steam"
	"github.com/this-is-a-bot/bot/tracker"
)

// Global instance.
var (
	db *sql.DB
	rs redis.RedisStore
)

func setup() {
	var err error
	if os.Getenv("ENV") == "HEROKU" {
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
		rs = redis.NewStore(os.Getenv("REDIS_URL"))
	} else {
		// Local dev.
		db, err = sql.Open("postgres", "dbname=bot sslmode=disable")
		rs = redis.NewStore("redis://127.0.0.1:6379")
	}

	if err != nil {
		// Fatal error, stop.
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/steam/discounts", handleSteamDiscounts)
	http.HandleFunc("/steam/featured", handleSteamFeatured)
	http.HandleFunc("/tracker/listing/text", handleTrackerListingText)
	http.HandleFunc("/tracker/marking/text", handleTrackerMarkingText)

	// Init database & redis.
	setup()
	defer db.Close()

	hostport := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Server running on %s\n", hostport)
	log.Fatal(http.ListenAndServe(hostport, nil))
}

// Dummy index handler.
func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome! I am a bot.")
}

/* Steam. */

// Return a list of discounted steam games in JSON format.
func handleSteamDiscounts(w http.ResponseWriter, r *http.Request) {
	games, err := steam.GetDiscounts(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(games)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Return a list of featured steam games in JSON format
func handleSteamFeatured(w http.ResponseWriter, r *http.Request) {
	feature := r.FormValue("feature")
	games, err := steam.GetFeatured(db, feature)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(games)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/* Tracker. */

type ByID []tracker.Catalog

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Return plain texts of tracking list.
func handleTrackerListingText(w http.ResponseWriter, r *http.Request) {
	username, app := r.FormValue("username"), r.FormValue("app")
	switch r.Method {
	case "POST":
		// POST method for adding new tracking.
		name, unit := r.PostFormValue("name"), r.PostFormValue("unit")
		if name == "" {
			http.Error(w, "'name' field is required", http.StatusBadRequest)
			return
		}
		_, err := tracker.AddTracking(db, username, app, name, unit)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	catalogs, err := tracker.GetTrackingCatalogs(db, username, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(catalogs) == 0 {
		http.Error(w, "no catalogs for such user", http.StatusNotFound)
		return
	}

	sort.Sort(ByID(catalogs))
	res := make([]string, 0)
	for _, catalog := range catalogs {
		s := fmt.Sprintf("%d. %s", catalog.ID, catalog.Name)
		if catalog.Done {
			if catalog.Value > 0 {
				s += fmt.Sprintf(": %v", catalog.Value)
				if catalog.Unit != "" {
					s += " " + catalog.Unit
				}
			} else {
				s += ": done"
			}
		} else {
			s += ": x"
		}
		res = append(res, s)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strings.Join(res, "\n")))
}

// Mark event done, then return plain texts of tracking list.
func handleTrackerMarkingText(w http.ResponseWriter, r *http.Request) {
	// Only allow POST.
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	catalogIDText, valueText := r.PostFormValue("catalogID"), r.PostFormValue("value")
	catalogID, err := strconv.Atoi(catalogIDText)
	if err != nil {
		http.Error(w, "'catalogID' must be an integer", http.StatusBadRequest)
		return
	}

	var value float64 = 0
	if valueText != "" {
		value, err = strconv.ParseFloat(valueText, 64)
		if err != nil {
			http.Error(w, "'value' must be a float", http.StatusBadRequest)
			return
		}
	}

	err = tracker.MarkDone(db, catalogID, value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// FIXME: A hack to reuse the code by changing the method to GET to avoid creating tracking.
	r.Method = "GET"
	handleTrackerListingText(w, r)
}
