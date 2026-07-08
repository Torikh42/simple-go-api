# TAHAP 1: Builder
# Menggunakan Golang versi Alpine (Sangat Ringan)
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency terlebih dahulu untuk Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh file proyek
COPY . .

# Build Binary secara Statis (Tanpa butuh library eksternal C/C++)
# Flag -ldflags="-w -s" digunakan untuk membuang debug info sehingga ukuran binary jauh lebih kecil
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api ./cmd/api/main.go

# TAHAP 2: Production (Image Final)
# Menggunakan OS 'scratch' (OS 100% kosong tanpa apa-apa, ukuran 0 MB!)
FROM scratch

WORKDIR /

# Salin binary file hasil build dari tahap 1 ke tahap 2
COPY --from=builder /app/api /api

# EXPOSE port (Hanya sebagai dokumentasi untuk Docker)
EXPOSE 8080

# Jalankan API
ENTRYPOINT ["/api"]
