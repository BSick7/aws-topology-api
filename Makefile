SHELL := /bin/bash

DOCKER_SERVER := docker.mdl.zone
IMAGE_VERSION := $(shell cat VERSION)
IMAGE_BASE := $(DOCKER_SERVER)/deployer/deploy
IMAGE_TAG := $(IMAGE_BASE):$(IMAGE_VERSION)
IMAGE_LATEST := $(IMAGE_BASE):latest

.PHONY: deps build vet test cover

deps:
	go get -u github.com/Masterminds/glide
	go get github.com/jstemmer/go-junit-report
	go get github.com/modocache/gover
	go get github.com/mattn/goveralls
	glide install

build:
	go build -ldflags "-X main.Version=`cat VERSION`"

install:
	go install -ldflags "-X main.Version=`cat VERSION`"

vet:
	glide nv | xargs go vet

test:
	set -o pipefail;glide nv \
		| xargs go test -v \
		| tee /dev/tty \
		| go-junit-report > unit-tests.xml

cover:
	set -e; \
	glide nv \
		| sed 's/\s/\n/g' \
		| sed 's/\/\.\.\.//' \
		| awk '{print "go test -coverprofile="$$1".coverprofile "$$1"...\0"}' \
		| xargs -0 -n1 bash -c; \
	gover;

release:
	go get github.com/mitchellh/gox
	go get github.com/tcnksm/ghr
	gox -os "linux darwin windows" -arch "amd64 386" -ldflags "-X main.Version=`cat VERSION`" -output="dist/deploy_{{.OS}}_{{.Arch}}"
	ghr -t $$GITHUB_TOKEN -u deployer -r deploy --replace `cat VERSION` dist/

docker-build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=`cat VERSION`" -a -installsuffix cgo -o deploy .
	docker build -t $(IMAGE_TAG) -f Dockerfile .
	docker tag $(IMAGE_TAG) $(IMAGE_LATEST)

push:
	@docker login --email="$(DOCKER_EMAIL)" --username="$(DOCKER_USER)" --password="$(DOCKER_PASS)" $(DOCKER_SERVER)
	docker push $(IMAGE_TAG)
	docker push $(IMAGE_LATEST)
