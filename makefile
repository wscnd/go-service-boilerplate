SHELL := /bin/bash

run:
	go run main.go

build:
	go build -ldflags "-X main.build=local" -o service main.go

# ============================================
# BUILD CONTAINERS

VERSION:= 0.1

all: serviceapi

serviceapi:
	docker build \
		-f zarf/docker/dockerfile \
		-t ghcr.io/wscnd/service:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.