APP_NAME=merch
PKG=./...
COVERFILE=coverage.out

.PHONY: test test-unit test-integration run cover up down

test:
	go test -v -race -cover -coverprofile=$(COVERFILE) $(PKG)

test-unit:
	go test -v -race -cover -coverprofile=$(COVERFILE) -tags='!e2e' $(PKG)

test-integration:
	go test -v -race -cover -coverprofile=$(COVERFILE) -tags=e2e $(PKG)

run:
	go run ./cmd/server

cover:
	go tool cover -func=$(COVERFILE)

up:
	docker compose up -d --build
down:
	docker compose down
clean:
	docker compose down -v
