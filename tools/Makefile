SHELL := $(shell which bash)

.PHONY: obs-run ws-run rabbit-run request-test

http-connector:
	cd httpConnector && \
		go run main.go httpClientWrapper.go

ws-connector:
	cd wsConnector && \
		go run main.go

rabbit-publisher:
	cd rabbitPublisher && \
		python3 receive.py

ws-publisher:
	cd wsPublisher && \
		go run main.go


# #########################
# Test Data
# #########################

test_data_file = ./requests.json

# it can be "push", "revert", "finalized"
event_type = "push"

request-test:
	curl -vvv \
		-X POST \
		-H "Content-Type: application/json" \
		-u aaaa:aaa \
		-d @${test_data_file} \
		http://localhost:5000/events/${event_type}
