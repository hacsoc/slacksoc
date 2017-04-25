FROM golang:alpine

RUN apk add --update git && rm -rf /var/cache/apk/*

COPY . $GOPATH/src/github.com/hacsoc/slacksoc

RUN set -x \
	&& go get -v github.com/brenns10/slacksoc/slacksoc

ADD startup.sh /usr/bin/
ADD slacksoc.yaml /etc/

CMD ["startup.sh"]
