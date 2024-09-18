run: 
	go run ./

build: 
	export CGO_ENABLED=1
	go build
