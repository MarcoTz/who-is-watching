package main 

import (
	"fmt"
	 "rooxo/whoiswatching/database"
   "rooxo/whoiswatching/telegram_bot/bot"
)

func main() {
  db,err := database.ConnectDB("./watchers.db")
	    if err != nil {
	      fmt.Printf("Could not initialize database: %s",err)
	      return
	    }

	err = bot.RunBot(db)
	if err != nil {
		fmt.Printf("%s", err)
	}
}
