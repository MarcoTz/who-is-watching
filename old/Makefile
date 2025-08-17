bot: 
	go run ./telegram_bot

server:
	- killall webserver
	go run ./webserver 

build: 
	CGO_ENABLED=1 go build -o whoiswatching_bot ./telegram_bot
