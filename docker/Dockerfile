FROM golang:latest
RUN apt-get update
RUN apt-get install -y mosquitto
RUN mkdir -p /go/src/github.com/geodan/gost/
ADD . /go/src/github.com/geodan/gost/
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/yaml.v2
RUN go get github.com/lib/pq
RUN go get github.com/eclipse/paho.mqtt.golang
RUN go build github.com/geodan/gost
ENTRYPOINT /go/bin/gost
EXPOSE 8080
