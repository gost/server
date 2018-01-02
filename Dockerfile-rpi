FROM resin/raspberry-pi-golang AS gost-build
WORKDIR /go/src/github.com/gost/server
ADD . .
RUN go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && go build -o /gostserver/gost github.com/gost/server \
    && cp config.yaml /gostserver/config.yaml


# final stage
FROM hypriot/rpi-alpine
WORKDIR /app
COPY --from=gost-build /gostserver /app/
ENTRYPOINT /app/gost
EXPOSE 8080