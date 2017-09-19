FROM golang:1.9.0-alpine3.6
WORKDIR /go/src/github.com/gost/server
ADD . .
RUN apk add --update --no-cache git \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && go build -o /gostserver/gost github.com/gost/server \
    && cp config.yaml /gostserver/config.yaml
ENTRYPOINT ["/gostserver/gost"]
EXPOSE 8080
