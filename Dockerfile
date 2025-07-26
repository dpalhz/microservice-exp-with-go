FROM golang:1.24-alpine AS builder

ARG SERVICE_NAME

WORKDIR /app

COPY go.mod go.sum./
RUN go mod download

COPY .. /app

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main./cmd/${SERVICE_NAME}

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main

COPY ./configs /app/configs

CMD ["./main"]