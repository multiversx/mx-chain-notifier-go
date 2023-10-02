SHELL := $(shell which bash)

debugger := $(shell which dlv)

.DEFAULT_GOAL := help

# #########################
# Base commands
# #########################

test:
	@echo "  >  Running unit tests"
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./...



# #########################
# Manage Notifier locally
# #########################

.PHONY: help build run runb kill debug debug-ath

cmd_dir = cmd/notifier
binary = event-notifier

help:
	@echo -e ""
	@echo -e "Make commands:"
	@grep -E '^[a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":"}; {printf "\t\033[36m%-30s\033[0m\n", $$1}'
	@echo -e ""

build:
	cd ${cmd_dir} && \
		go build -v -ldflags="-X main.appVersion=$(git describe --tags --long --dirty)" -o ${binary}

publisher_type="rabbitmq"
run: build
	cd ${cmd_dir} && \
		./${binary} --publisher-type=${publisher_type} --log-level="*:DEBUG"

runb: build
	cd ${cmd_dir} && \
		(./${binary} --publisher-type=${publisher_type} & echo $$! > ./${binary}.pid)

kill:
	kill $(shell cat ${cmd_dir}/${binary}.pid)

debug: build
	cd ${cmd_dir} && \
		${debugger} exec ./${binary} -- --publisher-type=${publisher_type}

debug-ath:
	${debugger} attach $$(cat ${cmd_dir}/${binary}.pid)

# Run local instance with Docker
image = "notifier"
image_tag = "latest"
container_name = notifier

dockerfile = Dockerfile

docker-build:
	docker build \
		-t ${image}:${image_tag} \
		-f ${dockerfile} \
		.

docker-new: docker-build
	docker run  \
		--detach \
		--network "host" \
		--name ${container_name} \
		${image}:${image_tag} \
		--publisher-type ${publisher_type} --log-level="*:DEBUG"

docker-start:
	docker start ${container_name}

docker-stop:
	docker stop ${container_name}

docker-logs:
	docker logs -f ${container_name}

docker-rm: docker-stop
	docker rm ${container_name}


# #########################
# System testing
# #########################

.PHONY: compose-new compose-start compose-stop

# Notifier with Redis sentinel and RabbitMQ
compose-new:
	docker-compose up -d

compose-start:
	docker-compose start

compose-stop:
	docker-compose stop

compose-rm:
	docker-compose down


# #########################
# Test Data
# #########################

test_data_file = ./examples/requests.json

# it can be "push", "revert", "finalized"
event_type = "push"

request-test:
	curl \
		-X POST \
		-H "Content-Type: application/json" \
		-u aaaa:aaa \
		-d @${test_data_file} \
		http://localhost:5000/events/${event_type}
