FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/api -o pack-calculator .

# Create final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the binary and static folder from builder stage
COPY --from=builder /app/pack-calculator .
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["./pack-calculator"]
