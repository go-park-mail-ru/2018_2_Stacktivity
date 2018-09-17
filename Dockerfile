FROM golang:1.9-alpine
WORKDIR /go/src

COPY . .
WORKDIR /go/src/cmd/server-public-api

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o main

CMD ["/go/src/cmd/server-public-api/main"]
EXPOSE 3000