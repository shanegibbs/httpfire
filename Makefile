SHELL := /bin/bash

-include Makefile.overrides.mk

KUBECTL=kubectl
IMAGE=httpfire

agent:
	go run ./cmd/agent/main.go

build:
	docker build -t $(IMAGE) .

run:
	docker run --rm $(IMAGE)

apply:
	$(KUBECTL) apply -f resources

delete:
	$(KUBECTL) delete -f resources
