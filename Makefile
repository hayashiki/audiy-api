GCP_PROJECT := $(shell gcloud config get-value project)
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
SERVICE := datastore-emulator

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
