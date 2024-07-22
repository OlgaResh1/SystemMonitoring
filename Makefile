BIN := "./bin/SystemStatsDaemon"
BIN_CLIENT := "./bin/SystemStatsClient"
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./...

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/service

build-client:
	go build -v -o $(BIN_CLIENT) -ldflags "$(LDFLAGS)" ./cmd/client

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run  --config=.golangci.yml ./...

test: 
	go test -race ./internal/... 

run: build
	$(BIN) --config ./configs/configd.toml

run-client: build-client
	$(BIN_CLIENT) --config ./configs/configc.toml

.PHONY: build run run-scheduler run-sender build-img run-img version test lint integration-tests build-client run-client
