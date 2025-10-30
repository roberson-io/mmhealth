.PHONY: build
GO_PACKAGES=$(shell go list ./...)
GO ?= $(shell command -v go 2> /dev/null)
BUILD_HASH ?= $(shell git rev-parse HEAD)
BUILD_VERSION ?= $(shell git ls-remote --tags --refs https://github.com/coltoneshaw/mmhealth.git | tail -n1 | sed 's/.*\///')

DOCKER_IMAGE_PROD ?= ghcr.io/coltoneshaw/mmhealth
DOCKER_IMAGE_DEV ?= mmhealth

BUILD_ENV ?= dev

ifeq ($(BUILD_ENV),prod)  
	LDFLAGS += -X "github.com/coltoneshaw/mmhealth/mmhealth.GitCommit=$(BUILD_HASH)"
	LDFLAGS += -X "github.com/coltoneshaw/mmhealth/mmhealth.GitVersion=$(BUILD_VERSION)"
	else 
endif


BUILD_COMMAND ?= go build -ldflags '$(LDFLAGS)' -o ./bin/mmhealth 

build: test
	mkdir -p bin
	$(BUILD_COMMAND)

buildDocker: build
	docker build --platform=linux/amd64 -f ./docker/Dockerfile -t $(DOCKER_IMAGE_DEV) .

run:
	go run ./main.go

package: test
	mkdir -p build bin 

	@echo Build Linux amd64
	env GOOS=linux GOARCH=amd64 $(BUILD_COMMAND)
	tar cf - -C bin mmhealth | gzip -9 > build/linux_amd64.tar.gz


	@echo Build OSX amd64
	env GOOS=darwin GOARCH=amd64 $(BUILD_COMMAND)
	tar cf - -C bin mmhealth | gzip -9 > build/darwin_amd64.tar.gz

	@echo Build OSX arm64
	env GOOS=darwin GOARCH=arm64 $(BUILD_COMMAND)
	tar cf - -C bin mmhealth | gzip -9 > build/darwin_arm64.tar.gz

	@echo Build Windows amd64
	env GOOS=windows GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o ./bin/mmhealth.exe 
	zip -9 build/windows_amd64.zip ./bin/mmhealth.exe

	rm ./bin/mmhealth ./bin/mmhealth.exe

check-style:
# https://stackoverflow.com/a/677212/1027058 (check if a command exists or not)
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed. Please see https://github.com/golangci/golangci-lint#install for installation instructions."; \
		exit 1; \
	fi; \

	@echo Checking gofmt
	@test -z "$$(gofmt -l . | tee /dev/stderr)" || (echo "Files not formatted with gofmt"; exit 1)
	@echo Running golangci-lint
	golangci-lint run --timeout 5m ./...

test: check-style
	@echo Running tests
	$(GO) test -race -cover -v $(GO_PACKAGES)

verify-gomod:
	$(GO) mod download
	$(GO) mod verify

make plugins:
	bash ./scripts/update_plugins.sh