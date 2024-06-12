#!/usr/bin/env make

.PHONY: dc up down start stop restart status

dc:
	@docker compose ${ARGS}

up:
	@docker compose up -d --remove-orphans

down:
	@docker compose down

start:
	@docker compose start

stop:
	@docker compose stop

restart: stop start

status:
	@docker compose ps
