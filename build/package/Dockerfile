# ===========================
# Stage 1: Build binary
# ===========================
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build binary ke /app/bin/dashboard
RUN go build -o bin/dashboard ./cmd/dashboard

# ===========================
# Stage 2: Runtime (lebih kecil)
# ===========================
FROM alpine:3.19

WORKDIR /app

# Tambahkan user non-root demi keamanan
RUN adduser -D appuser

# Copy binary dan folder yang dibutuhkan
COPY --from=builder /app/bin/dashboard .
COPY web ./web
COPY assets ./assets

USER appuser
EXPOSE 3000

CMD ["./dashboard"]
