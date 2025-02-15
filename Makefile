-include .env

export GO111MODULE=on
export GOPROXY=direct

.PHONY: deps
deps:
	@echo 'install dependencies'
	go mod tidy -v


.PHONY: test
test:
	go test -v -short ./...

.PHONY: up
up:
	docker-compose up -d --build

.PHONY: down
down:
	docker-compose down
