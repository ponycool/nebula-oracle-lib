.PHONY: test
test:
	go test -v ./test/...

.PHONY: deploy
deploy:
	cp .env.dev .env