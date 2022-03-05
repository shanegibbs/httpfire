SHELL := /bin/bash

-include Makefile.overrides.mk

KUBECTL=kubectl
CURL=curl

IMAGE=httpfire
LISTEN_ADDR=127.0.0.1:8080

test:
	go test ./...

local:
	go run ./cmd/local/main.go

single:
	go run ./cmd/single/main.go

agent:
	go run ./cmd/agent/main.go

director:
	HTTPFIRE_CONFIG=examples/director/config.yaml go run ./cmd/director/main.go

docker:
	docker build -t $(IMAGE) .

up:
	docker-compose --project-name httpfire up --force-recreate --build

down:
	docker-compose --project-name down

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
	$(CURL) --data @examples/plans/default.json http://$(LISTEN_ADDR)/start

clean: down

ignore-overrides-file:
	git update-index --assume-unchanged Makefile.overrides.mk

nginx:
	docker run --rm --name httpfire-nginx -p 9090:80 nginx

http-echo:
	docker run --rm --name httpfire-echo -p 9090:80 mendhak/http-https-echo
