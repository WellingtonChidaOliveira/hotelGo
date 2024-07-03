build:
	@go build -o bin/api

run:build
	@./bin/api

test: 
	@go test -v ./...

db-up:
	@docker-compose up -d

db-down:
	@docker-compose down