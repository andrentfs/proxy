FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o proxy cmd/proxy/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o relay cmd/relay/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/proxy .
COPY --from=builder /app/relay .

EXPOSE 1080 8080

CMD ["./proxy"]
