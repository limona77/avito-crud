deps:
	go mod tidy
run:
	go run cmd/main.go

dcu:
	docker-compose up
