FROM golang:1.11.4-alpine3.7 AS builder

RUN apk update && apk add --no-cache git

RUN mkdir -p /go/src/github.com/imulab/drone-webhook-proxy
ADD . /go/src/github.com/imulab/drone-webhook-proxy
WORKDIR /go/src/github.com/imulab/drone-webhook-proxy

RUN go get -u github.com/spf13/cobra/cobra
RUN go get -u github.com/spf13/pflag
RUN go get -u github.com/go-redis/redis
RUN go get -u github.com/sirupsen/logrus

RUN go build -o /usr/local/bin/hook .

FROM alpine:3.7

COPY --from=builder /usr/local/bin/hook /usr/local/bin/hook

ENTRYPOINT ["/usr/local/bin/hook"]