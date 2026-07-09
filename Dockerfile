# TAHAP 1: Builder
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api/main.go

# TAHAP 2: Runner (alpine bukan scratch agar health check Railway berjalan)
FROM alpine:latest

# Instal CA certificates agar bisa konek ke HTTPS (Neon butuh ini)
RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/api /app/api

EXPOSE 8080

ENTRYPOINT ["/app/api"]
