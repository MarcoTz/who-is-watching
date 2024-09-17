package main

import (
  "rooxo/whoiswatching/telegram_bot"
	"fmt"
//	"rooxo/whoiswatching/types"
//  "rooxo/whoiswatching/database"
)

func main() {
  err := telegram_bot.RunBot()
  if err != nil {
    fmt.Printf("%s",err)
  }
/*  db,err := database.ConnectDB("./watchers.db")
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
*/

}
