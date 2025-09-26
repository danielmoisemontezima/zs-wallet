# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /wallet-backend ./cmd/api

# Final stage
FROM alpine:latest
WORKDIR /
COPY --from=builder /wallet-backend /wallet-backend
# Add this line to copy the migrations directory into the final image
COPY migrations ./migrations

EXPOSE 8080
CMD ["/wallet-backend"]