## Copied from jcmoraisjr/haproxy-ingress

.PHONY: default
default: build

REPO_LOCAL=localhost/kpng-ipvs-backend
#include container.mk

GOOS=linux
GOARCH?=amd64

.PHONY: build
build:
	go mod tidy && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/kpng-backend-ipvs ./cmd/	

image:
	docker build -t $(REPO_LOCAL) .

