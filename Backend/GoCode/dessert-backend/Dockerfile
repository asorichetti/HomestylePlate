FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

#Build desserts backend binary from dessertServer.go
RUN go build -o desserts-backend dessertServer.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/desserts-backend .
EXPOSE 8081
CMD ["./desserts-backend"]