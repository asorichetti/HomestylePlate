# Build stage: compile Go Backend
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy Go source code from GoCode folder inside Backend
COPY ./GoCode/ .

# Download dependencies and build binary 
RUN go mod download 
RUN go build -o backend-server .

# Final stage: minimal runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/backend-server .
EXPOSE 8080
CMD ["./backend-server"]