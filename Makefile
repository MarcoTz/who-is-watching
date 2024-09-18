bot: 
	go run ./telegram_bot

build: 
	export CGO_ENABLED=1
	go build
