# TAHAP 1: Builder
# Menggunakan Golang versi Alpine (Sangat Ringan)
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api/main.go

FROM scratch

WORKDIR /

COPY --from=builder /app/api /api

EXPOSE 8080

ENTRYPOINT ["/api"]
