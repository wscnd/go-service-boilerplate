
SHELL := /bin/bash

# ============================================
# basic run & build
run:
	go run main.go

build:
# main.build is a var that is located in the main file that can be configurable via flags
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

# ============================================
# k8s/kind
KIND_CLUSTER_NAME := service-start-cluster

# https://github.com/kubernetes-sigs/kind/releases/
KIND_VERSION := kindest/node:v1.29.2@sha256:51a1434a5397193442f0be2a297b488b6c919ce8a3931be0ce822606ea5ca245

dev-up:
	kind create cluster \
		--image $(KIND_VERSION) \
		--name $(KIND_CLUSTER_NAME) \
		--config zarf/k8s/kind/kind-config.yaml

dev-down:
	kind delete cluster --name $(KIND_CLUSTER_NAME)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces