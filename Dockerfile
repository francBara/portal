FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o portal .

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/portal .

EXPOSE 8080

CMD ["./portal", "patcher"]
