FROM golang:alpine

RUN apk add --update git && rm -rf /var/cache/apk/*

COPY . $GOPATH/src/github.com/hacsoc/slacksoc

RUN set -x \
	&& cd $GOPATH/src/github.com/hacsoc/slacksoc \
	&& go get -d -v github.com/hacsoc/slacksoc \
	&& go build -o /usr/bin/slacksoc . \
	&& rm -rf $GOPATH

ADD startup.sh /usr/bin/

CMD ["startup.sh"]
