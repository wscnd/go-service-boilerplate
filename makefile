# ============================================
# DEPS
APP_VERSION := 0.1
SHELL := /bin/bash

K8S_KIND_CLUSTER_NAME := service-start-cluster
# https://github.com/kubernetes-sigs/kind/releases/
K8S_KIND_VERSION := kindest/node:v1.29.2@sha256:51a1434a5397193442f0be2a297b488b6c919ce8a3931be0ce822606ea5ca245
K8S_NAMESPACE := service-system

# ============================================
# basic run & build

run:
	go run main.go

build:
# main.build is a var that is located in the main file that can be configurable via flags
	go build -ldflags "-X main.build=local" -o service main.go

# ============================================
# BUILD CONTAINERS

all: serviceapi

serviceapi:
	docker build \
		-f zarf/docker/dockerfile \
		-t ghcr.io/wscnd/service:$(APP_VERSION) \
		--build-arg BUILD_REF=$(APP_VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ============================================
# k8s/kind

dev-up:
	kind create cluster \
		--image $(K8S_KIND_VERSION) \
		--name $(K8S_KIND_CLUSTER_NAME) \
		--config zarf/k8s/kind/kind-config.yaml

dev-load:
	kind load docker-image ghcr.io/wscnd/service:$(APP_VERSION) --name $(K8S_KIND_CLUSTER_NAME)

dev-down:
	kind delete cluster --name $(K8S_KIND_CLUSTER_NAME)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

