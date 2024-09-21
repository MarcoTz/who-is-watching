bot: 
	go run ./telegram_bot

build: 
	CGO_ENABLED=1 go build -o whoiswatching_bot ./telegram_bot
