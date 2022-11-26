FROM golang:1.17.13-alpine3.16

WORKDIR /data/gin_blog/
COPY . /data/gin_blog

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"

RUN go mod download && \
    go build -o main

EXPOSE 8000

ENTRYPOINT ["./main", "-h", "0.0.0.0"]