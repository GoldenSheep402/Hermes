FROM golang:latest AS builder

WORKDIR /build

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o App -ldflags "-s -w -extldflags '-static' -X 'github.com/GoldenSheep402/Hermes/conf.SysVersion=$(git show -s --format=%h)'" main.go

# Build the frontend using Node.js and pnpm
FROM node:latest AS frontend-builder

WORKDIR /frontend

COPY ./frontend .

RUN npm install -g pnpm \
    && pnpm install \
    && pnpm run build

# Caddy server stage
FROM caddy:latest AS frontend-server

WORKDIR /srv

COPY --from=frontend-builder /frontend/dist /srv

COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80

# Final stage with Alpine
FROM alpine:latest

WORKDIR /Serve
RUN mkdir -p /Serve/config

# Copy the Go app from the builder stage
COPY --from=builder /build/App /Serve/App

# Copy the frontend files from the frontend-server stage
COPY --from=frontend-server /srv /Serve/frontend

# Copy the Caddy binary from the frontend-server stage
COPY --from=frontend-server /usr/bin/caddy /usr/bin/caddy

# Install timezone data and set time zone
RUN apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# Expose port 80 for the Caddy server
EXPOSE 80

# Final command to start the Go application and the Caddy server
CMD sh -c "caddy run --config /etc/caddy/Caddyfile --adapter caddyfile & /Serve/App server -c /Serve/config/config.yaml"
