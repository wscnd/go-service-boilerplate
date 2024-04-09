SHELL_PATH = /bin/ash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/ash,/bin/bash)

# ==============================================================================
# DEPS
# ==============================================================================

BASE_IMAGE_NAME := localhost/wscnd/service
APP := sales
SERVICE_NAME := sales-api
APP_VERSION := 0.1
SERVICE_IMAGE_NAME := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(APP_VERSION)


# ------------------------------------------------------------------------------

K8S_KIND_CLUSTER_NAME := service-boilerplate-cluster
# https://github.com/kubernetes-sigs/kind/releases/
K8S_KIND_VERSION := kindest/node:v1.29.2@sha256:51a1434a5397193442f0be2a297b488b6c919ce8a3931be0ce822606ea5ca245
K8S_NAMESPACE := sales-system

# ==============================================================================
# BASIC RUN & BUILD
# ==============================================================================

run:
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

run-help:
	go run app/services/sales-api/main.go --help | go run app/tooling/logfmt/main.go

build:
# main.build is a var that is located in the main file that can be configurable via flags
	go build -ldflags "-X main.build=local" -o service app/services/sales-api/main.go

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
		-t $(SERVICE_IMAGE_NAME) \
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
		--config zarf/k8s/dev/kind-config.yaml

dev-down:
	kind delete cluster --name $(K8S_KIND_CLUSTER_NAME)

dev-watch:
	watch -n 0.3 kubectl get pods -o wide

dev-status:
	kubectl get svc,deploy,rs,nodes,pods --selector "app in ($(APP))" --all-namespaces

dev-logs:
	kubectl logs --namespace=$(K8S_NAMESPACE) -l app=$(APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)

port-forward:
	kubectl port-forward service/sales-api 3000:3000 4000:4000 --namespace sales-system

metrics-view:
	expvarmon -ports="localhost:4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

# ------------------------------------------------------------------------------

dev-load:
	kind load docker-image $(SERVICE_IMAGE_NAME) --name $(K8S_KIND_CLUSTER_NAME)

dev-apply:
	kustomize build zarf/k8s/base/sales | kubectl apply -f -

dev-restart:
	kubectl rollout restart deployment sales --namespace=$(K8S_NAMESPACE)
