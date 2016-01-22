FROM golang:1.5.3

COPY . $GOPATH/src/github.com/hacsoc/slacksoc

RUN set -x \
	&& cd $GOPATH/src/github.com/hacsoc/slacksoc \
	&& go get -d -v github.com/hacsoc/slacksoc \
	&& go build -o /usr/bin/slacksoc . \
	&& rm -rf $GOPATH

CMD ["slacksoc"]