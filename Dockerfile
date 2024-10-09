# Stage 1: Build the Go application
FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o App -ldflags "-s -w -extldflags '-static' -X 'github.com/GoldenSheep402/Hermes/conf.SysVersion=$(git show -s --format=%h)'" main.go

# Stage 2: Build the frontend using Node.js and pnpm
FROM node:latest AS frontend-builder

WORKDIR /frontend

COPY ./frontend .

RUN npm install -g pnpm \
    && pnpm install \
    && pnpm run build

FROM alpine:latest

WORKDIR /Serve

RUN apk update && apk add tzdata caddy supervisor \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone


#COPY config.yaml /Serve/config.yaml
#COPY Caddyfile /etc/caddy/Caddyfile
#
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

COPY --from=builder /build/App /Serve/App
COPY --from=frontend-builder /frontend/dist /Serve/frontend

EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]

