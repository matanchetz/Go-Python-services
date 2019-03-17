FROM golang:1.12.1-alpine3.9

WORKDIR /go/src/app/workers/

RUN apk add git && go get github.com/olivere/elastic && rm -rf /var/cache/apk/*

CMD ["go", "run", "server.go"]

