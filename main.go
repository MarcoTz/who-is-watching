package main

import (
  "fmt"
	"rooxo/whoiswatching/database"
  "rooxo/whoiswatching/types"
)

func main() {
  show := types.NewShow(1,"Hunter x Hunter")
	database.Load()
  fmt.Printf("%s",show)

}
