.PHONY: tools check clean test build package package-snapshot docs docker-image

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
DATE := $(shell date +'%Y-%m-%d %H:%M:%S')

export GO111MODULE=on
export SHA
export VERSION
export DATE

default: tools check test build

tools: $(shell go env GOPATH)/bin/golangci-lint $(shell go env GOPATH)/bin/goreleaser

$(shell go env GOPATH)/bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.23.6
	golangci-lint --version

$(shell go env GOPATH)/bin/goreleaser:
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | bash -s -- -b "${GOPATH}/bin"
	goreleaser --version

test:
	go test -v -cover ./...

clean:
	rm -rf dist/

build: clean
	@echo Version: $(VERSION)
	go build -v -ldflags '-X "main.Version=${VERSION}" -X "main.ShortCommit=${SHA}" -X "main.Date=${DATE}"' .

check:
	golangci-lint run

docker-image:
	docker build -t containous/acme-fixer:$(VERSION) .

release: tools
	goreleaser release

release-snapshot: tools
	goreleaser release --rm-dist --snapshot --skip-publish

doc:
	go run . doc
