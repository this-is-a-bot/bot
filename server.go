package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/this-is-a-bot/bot/steam"
    "github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// Global instance of database connection.
var db *sql.DB

func getDB() (db *sql.DB) {
	var err error
	if os.Getenv("ENV") == "HEROKU" {
		db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	} else {
		// Local dev.
		db, err = sql.Open("postgres", "dbname=bot sslmode=disable")
	}

	if err != nil {
		// Fatal error, stop.
		panic(err)
	}
	return
}

func main() {
    mux := mux.NewRouter()
	mux.HandleFunc("/steam/discounts", handleSteamDiscounts)
	mux.HandleFunc("/steam/featured/{feature:[a-z]+}", handleSteamFeatured)

	mux.HandleFunc("/", handleIndex)

	// Init database.
	db = getDB()
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
	params := mux.Vars(r)
	feature := params["feature"]
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
