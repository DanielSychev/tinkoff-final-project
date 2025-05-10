FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/main ./internal/cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/bin/main .
COPY --from=builder /app/internal/config/.env ./internal/config/.env
COPY --from=builder /app/internal/adapters/adrepo/migrations ./internal/adapters/adrepo/migrations

EXPOSE 8081 1011

CMD ["./main"]