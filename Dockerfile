# --- Build stage ---
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/bin/payment-splitter ./cmd/api

# --- Runtime stage ---
FROM alpine:3.19

RUN apk --no-cache add ca-certificates
RUN adduser -D -g '' appuser

COPY --from=builder /app/bin/payment-splitter /usr/local/bin/payment-splitter

USER appuser
EXPOSE 8080

ENTRYPOINT ["payment-splitter"]
