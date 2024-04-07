SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# DEPS
# ==============================================================================

APP_VERSION := 0.1
K8S_KIND_CLUSTER_NAME := service-start-cluster
# https://github.com/kubernetes-sigs/kind/releases/
K8S_KIND_VERSION := kindest/node:v1.29.2@sha256:51a1434a5397193442f0be2a297b488b6c919ce8a3931be0ce822606ea5ca245
K8S_NAMESPACE := service-system

# ==============================================================================
# BASIC RUN & BUILD
# ==============================================================================

run:
	go run main.go

build:
# main.build is a var that is located in the main file that can be configurable via flags
	go build -ldflags "-X main.build=local" -o service main.go

# ==============================================================================
# MODULES
# ==============================================================================

tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# BUILD CONTAINERS
# ==============================================================================

serviceapi:
	docker build \
		-f zarf/docker/dockerfile \
		-t ghcr.io/wscnd/service:$(APP_VERSION) \
		--build-arg BUILD_REF=$(APP_VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# K8S/KIND
# ==============================================================================

dev-init: serviceapi dev-up dev-load dev-apply
dev-update: serviceapi dev-load dev-restart
dev-update-apply: serviceapi dev-load dev-apply dev-restart

# ------------------------------------------------------------------------------

dev-up:
	kind create cluster \
		--image $(K8S_KIND_VERSION) \
		--name $(K8S_KIND_CLUSTER_NAME) \
		--config zarf/k8s/kind/kind-config.yaml

dev-down:
	kind delete cluster --name $(K8S_KIND_CLUSTER_NAME)

dev-watch:
	watch -n 0.3 kubectl get pods -o wide

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

dev-logs:
	watch -n 0.3 kubectl logs -l app=service --all-containers -f --tail=100 --namespace=$(K8S_NAMESPACE)

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image ghcr.io/wscnd/service:$(APP_VERSION) --name $(K8S_KIND_CLUSTER_NAME)

dev-apply:
	kustomize build zarf/k8s/kind/service-pod | kubectl apply -f -

dev-restart:
	kubectl rollout restart deployment service-pod --namespace=$(K8S_NAMESPACE)
