FROM golang:1.14-alpine

RUN apk add --update --no-cache bash ca-certificates git libc-dev make build-base
#ENV PATH /go/bin:$PATH
#ENV GOPATH /go
ENV MARKETPLACEPATH /go/src/github.com/corestario/marketplace
#RUN mkdir -p $MARKETPLACEPATH
WORKDIR $MARKETPLACEPATH

#COPY ./go.mod $MARKETPLACEPATH
#COPY ./modules $MARKETPLACEPATH
COPY . .

ENV GO111MODULE=on

RUN make install

#COPY . $MARKETPLACEPATH

EXPOSE 26657
EXPOSE 1317

ENTRYPOINT $MARKETPLACEPATH/run.sh