FROM golang:1.5.3

COPY . $GOPATH/src
WORKDIR $GOPATH/src/slacksoc

RUN go build .
CMD ["./slacksoc"]
