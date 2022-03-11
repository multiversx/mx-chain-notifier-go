SHELL := $(shell which bash)

debugger := $(shell which dlv)

.DEFAULT_GOAL := help

# #########################
# Manage Notifier locally
# #########################

.PHONY: help build run runb kill debug debug-ath

cmd_dir = cmd
binary = event-notifier

help:
	@echo -e ""
	@echo -e "Make commands:"
	@grep -E '^[a-zA-Z_-]+:.*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":"}; {printf "\t\033[36m%-30s\033[0m\n", $$1}'
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

docker-new:
	docker run  \
		--detach \
		--network "host" \
		--name ${container_name} \
		${image}:${image_tag}

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

.PHONY: light-new light-start light-stop

notifier_name = notifier

compose-build:
	docker-compose build

# Use only notifier
light-new: export API_TYPE = notifier
light-new:
	docker-compose up -d ${notifier_name}

light-start:
	docker-compose start ${notifier_name}

light-stop:
	docker-compose stop ${notifier_name}

# Notifier with Redis sentinel and RabbitMQ
compose-new: export API_TYPE = rabbit-api
compose-new:
	docker-compose up -d

compose-start:
	docker-compose start

compose-stop:
	docker-compose stop

compose-rm:
	docker-compose down
