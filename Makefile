build:
	@go build -o bin/api

run:build
	@./bin/api

seed:
	@go run scripts/seed.go

test: 
	@go test -v ./...

db-up:
	@docker-compose up -d

db-down:
	@docker-compose down