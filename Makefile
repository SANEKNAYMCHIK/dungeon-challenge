APP_NAME=dungeon-challenge

run:
	go run ./cmd/dungeon-challenge/main.go

build:
	go build -o bin/$(APP_NAME) ./cmd/dungeon-challenge/main.go

test:
	go test ./... -v

cover:
	go test ./... -cover

fmt:
	go fmt ./...

lint:
	golangci-lint run