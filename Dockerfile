# TAHAP 1: PEMBANGUNAN (Builder)
# Menggunakan mesin Golang untuk meng-compile kode sumber Anda
FROM golang:alpine AS builder

# Membuat direktori kerja di dalam kontainer
WORKDIR /app

# Menyalin seluruh kode Anda ke dalam kontainer
COPY . .

# Mengunduh semua pustaka (Fiber, GORM, dll)
RUN go mod download

# Meng-compile aplikasi menjadi sebuah file bernama 'server'
RUN go build -o server main.go

# TAHAP 2: MENJALANKAN APLIKASI (Runner)
# Menggunakan sistem operasi Alpine yang super ringan (hanya ~5MB)
FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app

# Mengambil file 'server' hasil compile dari Tahap 1
COPY --from=builder /app/server .

# Membuka port 3000 agar bisa diakses
EXPOSE 3000

# Perintah untuk menyalakan mesin
CMD ["./server"]