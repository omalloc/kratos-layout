FROM golang:1.21 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

FROM alpine:3.14

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
     && apk add --no-cache tzdata curl

COPY --from=builder /src/bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

ENV TZ=Asia/Shanghai

HEALTHCHECK --interval=5s --timeout=5s --start-period=3s --retries=3 \
    CMD curl -sS 'http://127.0.0.1:8000/health' || exit 1

CMD ["./server", "--conf", "/data/conf"]
