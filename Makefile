bot: 
	go run ./telegram_bot

build: 
	CGO_ENABLED=1 go build
