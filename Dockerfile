# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .

# Build the binary (use /app/server to avoid conflict with chatvault/ helm chart dir)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# Final stage
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/chatvault

EXPOSE 8080

CMD ["/app/chatvault"]
