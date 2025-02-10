
build:
	@go build -o bin/backend cmd/main.go

test: build
	@./bin/backend

migration: 
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migration-up:
	@go run cmd/migrate/main.go up

migration-down:
	@go run cmd/migrate/main.go down
