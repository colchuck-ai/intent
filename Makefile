.PHONY: build test tidy fmt

build:
	go build -o bin/intent ./cmd/intent

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...
