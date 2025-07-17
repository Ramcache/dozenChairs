run:
	go run ./app/main.go

build:
	go build -o dozenChairs ./app

tidy:
	go mod tidy
