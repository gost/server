FROM alpine:3.4
MAINTAINER Geodan <info@geodan.nl>

ADD . /go/src/github.com/geodan/gost/src/
ENV GOPATH=/go/
ENV CGO_ENABLED=0

    # Make directories
RUN apk add --no-cache --virtual .build-deps bash git openssh go && \
    mkdir -p /go/src/github.com/geodan/gost/ /go/bin/gost/client/ && \

    # Resolve dependencies
    go get github.com/gorilla/mux && \
    go get gopkg.in/yaml.v2 && \
    go get github.com/lib/pq && \
    go get github.com/eclipse/paho.mqtt.golang && \

    # Build application
    go build -a -ldflags '-extldflags "-static"' -o /go/bin/gost/gost github.com/geodan/gost/src && \
    cp -avr /go/src/github.com/geodan/gost/src/client /go/bin/gost && \
    cp /go/src/github.com/geodan/gost/src/config.yaml /go/bin/gost/config.yaml && \

    # Cleanup
    apk del .build-deps && \
    rm -rf /usr/local/go

WORKDIR /go/bin/gost
EXPOSE 8080
ENTRYPOINT /go/bin/gost/gost
