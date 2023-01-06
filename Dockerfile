FROM golang:1.16

RUN apt update

ENV GO111MODULE on
WORKDIR /go/src/work

# install go tools（自動補完等に必要なツールをコンテナにインストール）
RUN go get github.com/uudashr/gopkgs/v2/cmd/gopkgs
RUN go get github.com/ramya-rao-a/go-outline
RUN go get github.com/nsf/gocode
RUN go get github.com/acroca/go-symbols
RUN go get github.com/fatih/gomodifytags
RUN go get github.com/josharian/impl
RUN go get github.com/haya14busa/goplay/cmd/goplay
RUN go get github.com/go-delve/delve/cmd/dlv
RUN go get golang.org/x/lint/golint
RUN go get golang.org/x/tools/gopls