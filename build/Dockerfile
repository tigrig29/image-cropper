FROM golang:1.19.4-alpine

RUN apk update && apk add git

WORKDIR /go/src

ADD . /go/src

RUN go install golang.org/x/tools/cmd/goimports@latest