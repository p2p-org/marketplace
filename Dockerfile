FROM golang:1.12-alpine3.10

ENV APP=marketplace
RUN apk update
RUN apk upgrade
RUN apk add --no-cache bash ca-certificates git libc-dev make build-base

ENV PATH /go/bin:$PATH
ENV GOPATH /go
ENV MARKETPLACEPATH /go/src/github.com/dgamingfoundation/marketplace
RUN mkdir -p $MARKETPLACEPATH

COPY . $MARKETPLACEPATH

WORKDIR $MARKETPLACEPATH

ENV GO111MODULE=on

RUN $MARKETPLACEPATH/run.sh --no_run

EXPOSE 26657

ENTRYPOINT mpd start