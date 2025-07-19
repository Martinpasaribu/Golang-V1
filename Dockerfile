# Gunakan image dasar Go
FROM golang:1.21-alpine

# Buat direktori kerja
WORKDIR /app

# Copy semua file ke dalam container
COPY .env .env

# Build binary
RUN go build -o main .

# Render akan mengatur PORT lewat env
ENV PORT=8080

# Expose port (optional, Render abaikan ini)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]
