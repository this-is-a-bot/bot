package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/this-is-a-bot/bot/redis"
	"github.com/this-is-a-bot/bot/steam"
	"github.com/this-is-a-bot/bot/tracker"
)

// Global instance.
var (
	db *sql.DB
	rs redis.RedisStrore
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

// Return plain texts of tracking list.
func handleTrackerListingText(w http.ResponseWriter, r *http.Request) {
	username, app := r.FormValue("username"), r.FormValue("app")
	catalogs, err := tracker.GetTrackingCatalogs(db, username, app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else if len(catalogs) == 0 {
		http.Error(w, "no catalogs for such user", http.StatusNotFound)
		return
	}

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
				s += "done"
			}
		} else {
			s += ": x"
		}
		res = append(res, s)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(strings.Join(res, "\n")))
}
