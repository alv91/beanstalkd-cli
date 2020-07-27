FROM golang:1.13-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
WORKDIR /go/src/app

ADD ./ /go/src/app/
RUN go mod download
RUN go build -o ./beanstalkd-cli .

FROM alpine:latest
ENTRYPOINT ["./beanstalkd-cli"]
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

COPY --from=builder /go/src/app/beanstalkd-cli .
