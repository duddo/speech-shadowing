.PHONY: run build test clean help

run:
	go run ./cmd

build:
	go build -o speech-shadowing ./cmd

test:
	go test ./...

clean:
	rm -f speech-shadowing

help:
	@echo "Available commands:"
	@echo "  make run    - Run the application"
	@echo "  make build  - Build the binary"
	@echo "  make test   - Run tests"
	@echo "  make clean  - Remove built binary"
