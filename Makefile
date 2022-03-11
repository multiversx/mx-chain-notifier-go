SHELL := $(shell which bash)

debugger := $(shell which dlv)

cmd_dir = cmd
binary = event-notifier

.DEFAULT_GOAL := help

.PHONY: help build run runb kill debug debug-ath

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


# #####################
# Redis
# #####################

.PHONY: redis-new redis-start redis-stop

redis_name = main-redis
redis_port = 6379

redis-new:
	docker run \
		--name ${redis_name} \
		-d \
		-p ${redis_port}:6379 \
		redis:latest

redis-start:
	docker start ${redis_name}

redis-stop:
	docker stop ${redis_name}

redis-rm: redis-stop
	docker rm ${redis_name}

# #####################
# RabbitMQ
# #####################

rabbitmq_name = main-rabbit
rabbitmq_port = 5672

rabbitmq-new:
	docker run \
		-d \
		--name ${rabbitmq_name} \
		-p ${rabbitmq_port}:5672 \
		rabbitmq:3

rabbitmq-start:
	docker start ${rabbitmq_name}

rabbitmq-stop:
	docker stop ${rabbitmq_name}

rabbitmq-rm: rabbitmq-stop
	docker rm ${rabbitmq_name}
