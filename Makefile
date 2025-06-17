APP_NAME=merch
PKG=./...
COVERFILE=coverage.out

.PHONY: test test-unit test-integration run cover

## Запуск всех тестов
test:
	go test -v -race -cover -coverprofile=$(COVERFILE) $(PKG)

## Только юнит-тесты
test-unit:
	go test -v -race -cover -coverprofile=$(COVERFILE) -tags='!e2e' $(PKG)

## Только e2e тесты
test-integration:
	go test -v -race -cover -coverprofile=$(COVERFILE) -tags=e2e $(PKG)
run:
	go run ./cmd/server
cover:
	go tool cover -func=$(COVERFILE)
## Запуск приложения и seed скрипта
up:
	docker compose up --build -d
	docker-compose up migrate
	docker compose run --rm seed
down:
	docker compose down -v