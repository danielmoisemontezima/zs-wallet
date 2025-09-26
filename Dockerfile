# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet-backend ./cmd/api

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /wallet-backend .

# Copy config files if needed, e.g., .env
# COPY .env .

EXPOSE 8080

CMD ["./wallet-backend"]
