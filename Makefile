.PHONY: build test clean install

build:
	@mkdir -p bin
	go build -o bin/shadow .

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

install:
	go install

clean:
	rm -rf bin/
	rm -f coverage.out

.DEFAULT_GOAL := build
