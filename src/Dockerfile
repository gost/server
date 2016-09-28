FROM golang:latest
RUN apt-get update
RUN apt-get install -y mosquitto
RUN mkdir -p /go/src/github.com/geodan/gost/
ADD . /go/src/github.com/geodan/gost/src
RUN go get github.com/gorilla/mux
RUN go get gopkg.in/yaml.v2
RUN go get github.com/lib/pq
RUN go get github.com/eclipse/paho.mqtt.golang
ENV gost_db_host 192.168.40.10
RUN go build -o /go/bin/gost/gost github.com/geodan/gost/src
RUN mkdir -p /go/bin/gost/client/
RUN cp -avr /go/src/github.com/geodan/gost/src/client /go/bin/gost
RUN cp /go/src/github.com/geodan/gost/src/config.yaml /go/bin/gost/config.yaml
WORKDIR /go/bin/gost
ENTRYPOINT ["/go/bin/gost/gost"]
EXPOSE 8080
