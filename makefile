all: bin/bot
test: unit-test

PLATFORM=local

.PHONY: bin/bot
bin/bot:
	@docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM}

.PHONY: unit-test
unit-test:
	@docker build . --target unit-test