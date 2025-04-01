# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main .

FROM alpine:3.18
WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /app/main
COPY --from=builder /app/.env.example /app/.env

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD wget --spider http://localhost:8080/health || exit 1

EXPOSE 8080
ENTRYPOINT ["/app/main"]