FROM golang:1.24.5-alpine3.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/bin /app/cmd/server/main.go

FROM alpine:latest AS server

WORKDIR /app

COPY --from=builder /app/bin ./

EXPOSE 8000 9000

ENTRYPOINT [ "./bin" ]