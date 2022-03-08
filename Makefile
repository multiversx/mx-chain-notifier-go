SHELL := $(shell which bash)

debugger := $(shell which dlv)

cmd_dir = cmd
binary = event-notifier

.DEFAULT_GOAL := help

.PHONY: help build run runb kill debug debug-ath

help:
	@echo -e ""
	@echo -e "Make commands:"
	@grep -E '^[a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":"}; {printf "\t\033[36m%-30s\033[0m\n", $$1}'
	@echo -e ""

build:
	cd ${cmd_dir} && \
		go build -o ${binary}

api_type="notifier"
run: build
	cd ${cmd_dir} && \
		./${binary} --api-type=${api_type}

runb: build
	cd ${cmd_dir} && \
		(./${binary} --api-type=${api_type} & echo $$! > ./${binary}.pid)

kill:
	kill $(shell cat ${cmd_dir}/${binary}.pid)

debug: build
	cd ${cmd_dir} && \
		${debugger} exec ./${binary} -- --api-type=${api_type}

debug-ath:
	${debugger} attach $$(cat ${cmd_dir}/${binary}.pid)
