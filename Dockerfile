FROM golang:latest
RUN apt-get update
RUN mkdir -p /go/src/github.com/gost/server/
ADD . /go/src/github.com/gost/server
RUN go get ./...
RUN go build -o /go/bin/gost/gost github.com/gost/server
RUN cp /go/src/github.com/gost/server/config.yaml /go/bin/gost/config.yaml
WORKDIR /go/bin/gost
ENTRYPOINT ["/go/bin/gost/gost"]
EXPOSE 8080
