#!/usr/bin/env make

.PHONY: build build__clear

init:
	cp -f $(ENV_DIST_DIR)/docker.env.dist $(BUILD_DIR)/docker.env

build:
	@docker compose build

build__clear:
	@docker compose build --no-cache
