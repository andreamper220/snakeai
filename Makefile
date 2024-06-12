#!/usr/bin/env make

# Optionally include .env (it couldn't exist in case of fresh installation)
-include docker.env
export $(shell sed 's/=.*//' .env)

export COLOR__RED = \n\t\033[01;31m
export COLOR__GREEN = \n\t\033[01;32m
export COLOR__YELLOW = \n\t\033[01;33m
export COLOR__DEFAULT = \033[0m\n

###> IMPORTANT: define the COMPOSE_PROJECT_NAME here so
###> Makefile's targets will be able to perform shell magic.
###> Also, the value here and value in the .env should be the same, that's important
export COMPOSE_PROJECT_NAME=go_snake_ai

BUILD_DIR := .
ENV_DIST_DIR = ./env.d/dotenv.d

-include ./env.d/make.d/operations.mk
-include ./env.d/make.d/build.mk

default: up

# args
FIRST_ARG := $(firstword $(MAKECMDGOALS))
ARGS = $(filter-out $@,$(MAKEOVERRIDES) $(MAKECMDGOALS))
MAKEFILE_PATH := $(abspath $(firstword $(MAKEFILE_LIST)))

%:
	@:
