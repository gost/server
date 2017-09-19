FROM golang:1.9.0-alpine3.6 AS gost-build
WORKDIR /go/src/github.com/gost/server
ADD . .
RUN apk add --update --no-cache git \
    && go get -u github.com/golang/dep/cmd/dep \
    && dep ensure \
    && go build -o /gostserver/gost github.com/gost/server \
    && cp config.yaml /gostserver/config.yaml


# final stage
FROM alpine
WORKDIR /app
COPY --from=gost-build /gostserver /app/
ENTRYPOINT ./gost
EXPOSE 8080
