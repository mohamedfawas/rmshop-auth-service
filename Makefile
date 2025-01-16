.PHONY: build run test docker

build:
	go build -o auth-service ./cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test ./...

docker:
	docker build -t rmshop-auth-service . 