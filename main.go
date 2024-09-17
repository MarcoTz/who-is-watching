package main

import (
	"fmt"
	"rooxo/whoiswatching/types"
  "rooxo/whoiswatching/database"
)

func main() {
  db,err := database.ConnectDB("./watchers.db")
  if err != nil {
    fmt.Printf("Could not initialize database: %s",err)
    return 
  }
  shows, err := database.GetAllShows(db)
  if err != nil { 
    fmt.Printf("Could not get shows: %s",err)
    return 
  }

  rand_show := types.RandomShow(shows)
  fmt.Printf("random show: %s",rand_show.Name)


}
