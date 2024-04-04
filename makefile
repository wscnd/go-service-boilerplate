SHELL := /bin/bash

run:
	go run main.go

build:
	go build -ldflags "-X main.build=local" -o service main.go

