current_dir = $(shell pwd)
DOCKER_COMPOSE ?= $(shell which docker-compose)

all: build_ui run_docker

build_ui:
	cd $(current_dir)/webapp && yarn build

run_docker:
	$(DOCKER_COMPOSE) build
