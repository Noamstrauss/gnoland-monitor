FROM golang:1.23-alpine

RUN apk add --no-cache git make build-base curl

WORKDIR /gno
RUN git clone https://github.com/gnolang/gno.git . && \
    cd gno.land && \
    make build.gnoland

WORKDIR /gno/gno.land
VOLUME ["/gno/home"]

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]