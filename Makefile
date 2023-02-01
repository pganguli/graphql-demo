.PHONY: all generate tidy test fmt lint

all:
	go build

generate:
	go generate ./...

tidy:
	go mod tidy

test:
	go test -race .

fmt:
	gofmt -l -w **/*.go

lint:
	golangci-lint run
