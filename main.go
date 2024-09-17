package main

import (
	"fmt"
	"rooxo/whoiswatching/telegram_bot"
	//		"rooxo/whoiswatching/types"
	 "rooxo/whoiswatching/database"
)

func main() {
  db,err := database.ConnectDB("./watchers.db")
	    if err != nil {
	      fmt.Printf("Could not initialize database: %s",err)
	      return
	    }

	err = telegram_bot.RunBot(db)
	if err != nil {
		fmt.Printf("%s", err)
	}
	/*  	    shows, err := database.GetAllShows(db)
	    if err != nil {
	      fmt.Printf("Could not get shows: %s",err)
	      return
	    }

	    rand_show := types.RandomShow(shows)
	    fmt.Printf("random show: %s",rand_show.Name)
	*/

}
