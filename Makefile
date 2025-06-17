.PHONY: test test-unit test-e2e run cover up down clean migrate

migrate:
	migrate -path=./migrations -database=$(MERCH_DB_DSN) up

test:
	go test -v -race -cover -coverprofile=$(COVERFILE) $(PKG)

test-unit:
	go test -v -race -cover -coverprofile=$(COVERFILE) -tags='!e2e' $(PKG)

test-e2e:
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
