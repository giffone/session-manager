DATABASE_URL ?= postgres://user:password@host:port/db-name?sslmode=disable
DOCKER_IMAGE_NAME ?= session_manager_image
REQ_LOG ?= false

.PHONY: migrate build run run_local gen_proto

gen_proto:
	protoc -I=./proto --go_out=. --go-grpc_out=. cadets_time.proto

migrate:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

build:
	docker build -t $(DOCKER_IMAGE_NAME) .

run: migrate build
	docker run --name $(DOCKER_IMAGE_NAME) -d -e DATABASE_URL="$(DATABASE_URL)" -e REQ_LOG="$(REQ_LOG)" -p 9090:8080 -p 9191:8181 $(DOCKER_IMAGE_NAME)

run_local: migrate build
	docker run -d -e DATABASE_URL="$(DATABASE_URL)" -e REQ_LOG="$(REQ_LOG)" --network=host -p 9090:8080 -p 9191:8181 $(DOCKER_IMAGE_NAME)
