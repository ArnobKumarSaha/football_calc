REGISTRY   ?= ghcr.io/arnobkumarsaha
IMAGE      ?= football-calc
TAG        ?= latest
PLATFORMS  ?= linux/amd64,linux/arm64

.PHONY: build fmt lint container push

build:
	CGO_ENABLED=0 go build -mod=vendor -o bin/football-calc ./cmd/server

fmt:
	gofmt -w .
	goimports -w .

lint:
	golangci-lint run ./...

container:
	docker buildx build --platform $(PLATFORMS) -t $(REGISTRY)/$(IMAGE):$(TAG) .

push:
	docker buildx build --platform $(PLATFORMS) -t $(REGISTRY)/$(IMAGE):$(TAG) --push .
