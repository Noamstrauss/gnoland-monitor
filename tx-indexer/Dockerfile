FROM golang:1.23-alpine

RUN apk add --no-cache git make build-base curl

WORKDIR /app
RUN git clone https://github.com/gnolang/tx-indexer.git . && \
    make build

RUN rm -rf ${DB_PATH:-"/data/indexer-db"}

VOLUME ["/data"]

CMD ./build/tx-indexer start \
    --remote "${REMOTE_NODE:-gnoland:26657}" \
    --listen-address "${LISTEN_ADDRESS:-tx-indexer:8546}" \
    --db-path "${DB_PATH:-/data/indexer-db}"
