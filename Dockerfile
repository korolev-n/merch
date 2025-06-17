# Stage 1 — builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка основного приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o merch ./cmd/server/main.go

# Сборка seed-скрипта
RUN CGO_ENABLED=0 GOOS=linux go build -o seed ./cmd/seed/main.go

# Stage 2 — runtime
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/merch .
COPY --from=builder /app/seed .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./merch"]
