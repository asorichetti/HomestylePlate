# Build stage: compile Go Backend
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy Go source code from GoCode folder inside Backend
COPY go.mod go.sum ./
RUN go mod download

COPY . .

#Build meals backend Binary from randomizer.go
RUN go build -o meals-backend randomizer.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/meals-backend .

EXPOSE 8080
CMD ["./meals-backend"]

