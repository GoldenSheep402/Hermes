FROM golang:latest AS builder

COPY . /build

WORKDIR /build

RUN set -ex \
    && GO111MODULE=auto CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static' -X 'github.com/juanjiTech/jframe/conf.SysVersion=$(git show -s --format=%h)'" -o App

FROM alpine:latest

WORKDIR /Serve
RUN mkdir "config"

COPY --from=builder /build/App ./App

RUN ls -R

RUN apk update && apk add tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

ENTRYPOINT [ "/Serve/App", "server" ]