FROM golang:1.9

ADD . /go/src/redisproxy

RUN go get github.com/go-redis/redis
RUN go get github.com/hashicorp/golang-lru

RUN go test -v redisproxy/service
RUN go test -v redisproxy/util
RUN go install -v ./...

ENTRYPOINT ["/go/bin/redisproxy", "-redisIpAndPort=172.17.0.1:6379","-expiry=10","-capacity=100","-port=10000","-concurrentmax=1000","-workers=100"]

EXPOSE 10000