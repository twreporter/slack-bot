# all: bin/bot
# test: lint test
build: build

PLATFORM=local

.PHONY: help
help:	## Show this help
	@grep -E '^[a-z.A-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# .PHONY: bin/bot
# bin/bot:	## Build the bin file for the bot with desired platform
# 	@docker build . --target bin \
# 	--output bin/ \
# 	--platform ${PLATFORM}

# .PHONY: test
# unit-test:	## Run the unit test
# 	@docker build . --target unit-test

# .PHONY: lint
# lint:	## Run the lint
# 	@docker build . --target lint

.PHONY: build
build:	## Build a Docker image for containerized Slack bot
	@docker build -t slack-bot \
	--build-arg SLACK_AUTH_TOKEN="" \
	--build-arg SLACK_APP_TOKEN="" \
	--build-arg SLACK_CHANNEL_ID="" \
	--build-arg APP_ENV="dev" \
	.
