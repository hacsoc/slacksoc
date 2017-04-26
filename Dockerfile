FROM golang:alpine

RUN apk add --update git && rm -rf /var/cache/apk/*

RUN set -x \
	&& go get -v github.com/hacsoc/slacksoc

ADD startup.sh /usr/bin/
ADD slacksoc.yaml /etc/

CMD ["startup.sh"]
