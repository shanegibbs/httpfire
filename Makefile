SHELL := /bin/bash

-include Makefile.overrides.mk

KUBECTL=kubectl
CURL=curl

IMAGE=httpfire
LISTEN_ADDR=127.0.0.1:8080

test:
	go test ./...

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

get-config: curl-get-config
curl-get-%:
	$(CURL) http://$(LISTEN_ADDR)/$* |jq .

post-stop: curl-post-stop
post-restart: curl-post-restart
curl-post-%:
	$(CURL) -XPOST http://$(LISTEN_ADDR)/$*

post-start:
	$(CURL) --data @examples/default.json http://$(LISTEN_ADDR)/start

ignore-overrides-file:
	git update-index --assume-unchanged Makefile.overrides.mk
