**Цель проекта:** улучшение навыков разработки REST API (CRUD+JWT). 

**Merch** - это магазин, где сотрудники могут приобретать товары за монеты (coin). Новому пользователю выделяется 1000 монет, которые можно использовать для покупки товаров или перевода другим пользователям. Также можно посмотреть информацию о купленных товарах, текущем балансе и истории переводов.

## Структура проекта

В данном проекте реализована слоистая архитектура (Layered Architecture).

```
cmd/
├── server/       # Точка входа основного приложения
└── seed/         # Скрипт для заполнения БД

internal/
├── config/       # Конфигурация приложения
├── db/           # Работа с БД (миграции, сиды)
├── domain/       # Бизнес-сущности
├── logger/       # Логирование
├── repository/   # Работа с данными
│   ├── mocks/    # Моки для тестирования
├── server/       # Инициализация сервера
├── service/      # Бизнес-логика
└── transport/    # HTTP-обработчики
    └── http/
        ├── middleware/
        ├── request/
        └── response/
```

## API

Для просмотра документации API, перейдите по [ссылке на Swagger UI](https://editor.swagger.io/?url=https://raw.githubusercontent.com/avito-tech/tech-internship/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025/schema.json) 

## Технологии

- **Язык**: Go 1.24
- **Фреймворк**: Gin
- **База данных**: PostgreSQL
- **Аутентификация**: JWT
- **Логирование**: slog
- **Контейнеризация**: Docker
- **Нагрузочное тестирование**: k6

## Конфигурация

Для локального запуска приложения необходимо определить следующие переменные окружения в файле `.env`:
- **MERCH_DB_DSN**: DSN (Data Source Name) для подключения к базе данных PostgreSQL. Пример значения:

```.env
MERCH_DB_DSN=postgres://user:password@localhost:5432/mydatabase
```

- **JWT_SECRET**: Секретный ключ для подписи JWT (JSON Web Tokens). Пример значения:

```.env
JWT_SECRET=mysecretkey
```

## Миграции базы данных

Убедитесь, что переменная `MERCH_DB_DSN` определена в вашем файле `.env` перед выполнением миграции.
Для применения миграций к базе данных выполните следующую команду:

```bash
make migrate
```

## Инструкция по запуску

- Go (версия 1.24 или выше)
- Docker и Docker Compose

### Команды Makefile

#### Запуск Docker контейнеров

- Для запуска контейнеров:

```bash
make up
```

- Для остановки контейнеров:

```bash
make down
```

- Для остановки контейнеров и удаления томов:

```bash
make clean
```

#### Запуск приложения локально

```bash
make run
```

#### Запуск тестов

- Для запуска всех тестов:

```bash
make test
```

- Для запуска только юнит-тестов

```bash
make test-unit
```

- Для запуска e2e тестов:

```bash
make test-e2e
```

## Примеры запросов

### Регистрация

```bash
curl --location 'http://localhost:8080/api/auth' \
--header 'Content-Type: application/json' \
--data '{"username": "user_1", "password": "Password123"}'
```

### Перевод монет

```bash
curl --location 'localhost:8080/api/sendCoin' \
--header 'Authorization: Bearer ' \
--header 'Content-Type: application/json' \
--data '{
  "toUser": "user_1",
  "amount": 100
}'
```

### Покупка товара

```bash
curl --location 'localhost:8080/api/buy/t-shirt' \
--header 'Authorization: Bearer '
```

### Информация о пользователе

```bash
curl --location 'http://localhost:8080/api/info' \
--header 'Authorization: Bearer '
```
