
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate ./cmd/migrate

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/server /usr/local/bin/server
COPY --from=builder /bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/migrations ./migrations

RUN mkdir -p ./uploads/events

EXPOSE 8080

CMD ["server", "DATABASE_URL"]