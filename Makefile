run:
	go run ./cmd/server/...

build:
	go build -o bin/server ./cmd/server/...

up:
	docker-compose up -d

down:
	docker-compose down

tidy:
	go mod tidy
