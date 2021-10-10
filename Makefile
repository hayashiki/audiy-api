GCP_PROJECT = $(shell gcloud config get-value project)
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
SERVICE := datastore-emulator
BUILD_TAG=v0.0.1#$(shell git describe --tags --abbrev=0 HEAD)
BUILD_HASH=$(shell git rev-parse --short HEAD)
BUILD_BRANCH=$(shell git symbolic-ref HEAD |cut -d / -f 3)
BUILD_VERSION=${BUILD_TAG}-${BUILD_HASH}
BUILD_TIME=$(shell date +%F-%H:%m:%S)
PACKAGE=github.com/hayashiki/audiy-api

dev:
	docker compose up

deploy:
	gcloud builds submit --config=cloudbuild.yaml .

deploy-importer:
	gcloud builds submit --config=cloudbuild.importer.yaml .

local-build:
	DOCKER_BUILDKIT=1 docker build -t $(SERVICE) . -f cmd/ffprobe/Dockerfile

local-run:
	docker run -it --rm $(SERVICE)

datastore-build:
	DOCKER_BUILDKIT=1 docker build -t $(SERVICE) . -f deployments/docker/datastore/Dockerfile

build-api:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -ldflags="-w -s -X ${PACKAGE}/etc/version.Version=${BUILD_VERSION} -X ${PACKAGE}/etc/version.BuildTime=$(BUILD_TIME)" -v -o bin/gqlserver cmd/gqlserver/main.go
