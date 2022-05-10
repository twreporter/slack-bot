all: bin/bot
test: lint test

PLATFORM=local

.PHONY: bin/bot
bin/bot:
	@docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM}

.PHONY: test
unit-test:
	@docker build . --target unit-test

.PHONY: lint
lint:
	@docker build . --target lint