package main

import ( 
  "fmt"
  "net/http"
  "database/sql"
  "rooxo/whoiswatching/webserver/pages"
  "rooxo/whoiswatching/database"
)

func handle(w http.ResponseWriter, r *http.Request){
  fmt.Fprint(w,"Hello index")
}

func handleShows(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
  return func(w http.ResponseWriter, r*http.Request) { fmt.Fprint(w,pages.RenderShows(db)) }
}

func handleWatchers(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
  return func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w,pages.RenderWatchers(db)) }
}

func handleGroups(db *sql.DB) func(w http.ResponseWriter, r *http.Request){
  return func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w,pages.RenderGroups(db))  }
}

func main() {
  db, err := database.ConnectDB("watchers.db")
  if err != nil {
    fmt.Printf("Could not initialize database: %s\n",err)
    return 
  }

  http.HandleFunc("/", handle)
  http.HandleFunc("/shows",handleShows(db))
  http.HandleFunc("/watchers",handleWatchers(db))
  http.HandleFunc("/groups",handleGroups(db))
  http.ListenAndServe(":8080",nil)
}
