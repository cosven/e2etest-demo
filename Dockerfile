FROM hub.pingcap.net/mirrors/golang:1.16

WORKDIR /apps/e2etest

ADD ./bin/ /apps/e2etest/
