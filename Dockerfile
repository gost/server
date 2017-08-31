FROM golang:latest
RUN mkdir -p /go/src/github.com/gost/server/ && mkdir -p /gostserver/config/
ADD . /go/src/github.com/gost/server
WORKDIR /go/src/github.com/gost/server
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure && go build -o /gostserver/gost github.com/gost/server && cp config.yaml /gostserver/config/config.yaml
ENTRYPOINT ["/gostserver/gost", "-config", "/gostserver/config/config.yaml"]
EXPOSE 8080
