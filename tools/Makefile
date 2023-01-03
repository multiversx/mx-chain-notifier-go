SHELL := $(shell which bash)

.PHONY: ws-run rabbit-run request-test

obs-run:
	cd obs && \
		go run main.go

ws-run:
	cd wsClient && \
		go run main.go

rabbit-run:
	cd rabbitClient && \
		python3 receive.py


# #########################
# Test Data
# #########################

test_data_file = ./requests2.json

# it can be "push", "revert", "finalized"
event_type = "push"

request-test:
	curl -vvv \
		-X POST \
		-H "Content-Type: application/json" \
		-u aaaa:aaa \
		-d @${test_data_file} \
		http://localhost:5000/events/${event_type}
