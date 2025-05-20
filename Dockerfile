FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o portal ./cmd

FROM node:20-alpine

WORKDIR /root/

COPY package.json package-lock.json ./

RUN npm install -g vite
RUN npm install

RUN apk add --no-cache git

COPY --from=builder /app/portal .
COPY --from=builder /app/tools tools
COPY --from=builder /app/config.json .
COPY --from=builder /app/users.json .
COPY --from=builder /app/static static

EXPOSE 8080 3000

CMD ["./portal", "patcher"]
