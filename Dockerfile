FROM golang:latest
RUN mkdir -p /go/src/github.com/gost/server/ && mkdir -p /gostserver/config/
ADD . /go/src/github.com/gost/server
RUN go get ./... && go build -o /gostserver/gost github.com/gost/server && cp /go/src/github.com/gost/server/config.yaml /gostserver/config/config.yaml
ENTRYPOINT ["/gostserver/gost", "-config", "/gostserver/config/config.yaml"]
EXPOSE 8080
