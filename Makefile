.PHONY: gen-proto test benchmark docker-up

gen-proto:
	protoc --go_out=. --go-grpc_out=. proto/*.proto

test:
	go test ./...

benchmark:
	go run cmd/benchmark/main.go

docker-up:
	docker compose up --build