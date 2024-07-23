BIN := "./bin/SystemStatsDaemon"
BIN_CLIENT := "./bin/SystemStatsClient"
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)
DOCKER_IMG="statsdaemon:develop"
DOCKER_IMG_CLIENT="statsdaemon:develop"
DOCKER_CONR_NAME="statsdaemon"
DOCKER_CONR_NAME="statsdaemon"

generate:
	go generate ./...

build-img-service:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/service/Dockerfile .

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/service/Dockerfile .

run-img: build-img
	docker run -it --name $(DOCKER_CONTAINER)  -p 50051:50051 $(DOCKER_IMG)

version: build
	$(BIN) version

build:
	CGO_ENABLED=0 go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/service

build-client:
	CGO_ENABLED=0  go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run  --config=.golangci.yml ./...

test: 
	go test -race ./internal/... 

test-integration:
	go test -tags integration ./tests/integration/... -count 1 -v

run: build
	$(BIN) --config ./configs/configd.toml

run-client: build-client
	$(BIN_CLIENT) --config ./configs/configc.toml

.PHONY: build run run-scheduler run-sender build-img run-img version test lint integration-tests build-client run-client
