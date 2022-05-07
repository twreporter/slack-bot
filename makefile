all: bin/bot

PLATFORM=local

.PHONY: bin/bot
bin/bot:
	@docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM}