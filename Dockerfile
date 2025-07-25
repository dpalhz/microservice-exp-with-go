# Tahap 1: Build
FROM golang:1.24-alpine AS builder

# Argumen untuk menentukan layanan mana yang akan di-build
ARG SERVICE_NAME

WORKDIR /app

# Salin file mod dan sum untuk men-cache dependensi
COPY go.mod go.sum./
RUN go mod download

# Salin semua kode sumber
COPY .. /app

# Build biner untuk layanan yang ditentukan
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main./cmd/${SERVICE_NAME}

# Tahap 2: Final Image
FROM alpine:latest

WORKDIR /app

# Salin biner yang telah dikompilasi dari tahap builder
COPY --from=builder /app/main /app/main

# Salin file konfigurasi yang relevan
COPY ./configs /app/configs

# Perintah untuk menjalankan layanan
CMD ["./main"]