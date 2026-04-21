FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

FROM debian:11.7-slim

WORKDIR /bin

COPY --from=builder /app/api .
COPY --from=builder /app/.env .

EXPOSE 8080

ENTRYPOINT [ "./api"]