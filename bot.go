package main

import ( 
    "fmt"
    "log"
    "net/http"
    "os"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
   "database/sql"
)

func main() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", Index)
    router.HandleFunc("/steam_discounts", Steam_discount)
    
    log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome! I am a bot.")
}

func Steam_discount(w http.ResponseWriter, r *http.Request) {
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    rows, err := db.Query("select * from steam_discounts where timestamp=2016-06-11")
    fmt.Fprintln(w, rows)
}