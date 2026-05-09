.PHONY: all run build test

all: run 

run:
	go run cmd/KitsuneERP/main.go

build:
	go build -o KitsuneERP cmd/KitsuneERP/main.go

test:
	go test -v ./tests/...
