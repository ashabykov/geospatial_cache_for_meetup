-include .env

export GO111MODULE=on
export GOPROXY=direct

.PHONY: deps
deps:
	go mod tidy -v

.PHONY: test
test:
	go test -v -short ./...

.PHONY: colima-start
colima-start:
	colima start --vm-type vz --mount-type virtiofs --cpu 4 --memory 8

.PHONY: colima-stop
colima-stop:
	colima stop

.PHONY: colima-delete
colima-delete:
	colima delete

.PHONY: up
up:
	docker-compose -f ./docker-compose.yaml up -d --build

.PHONY: down
down:
	docker-compose -f ./docker-compose.yaml down

.PHONY: bench
bench:
	go test -run='^$' -bench=cmd/fan-out-write/BenchmarkGetAll -benchtime=10x -count=6 -timeout 60m > debug.profile
