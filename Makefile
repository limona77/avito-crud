deps:
	go mod tidy
run:
	go run cmd/main.go

docker-compose-run:
	docker-compose up
