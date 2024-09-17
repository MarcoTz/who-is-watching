package main

import (
	"fmt"
	"rooxo/whoiswatching/types"
)

func main() {
	show := types.NewShow(1, "Hunter x Hunter")
	fmt.Printf("%s", show)

}
