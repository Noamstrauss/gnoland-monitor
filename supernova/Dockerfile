FROM golang:1.23-alpine

RUN apk add --no-cache git make build-base curl

WORKDIR /app
RUN git clone https://github.com/gnolang/supernova.git . && \
    make build

CMD ./build/supernova \
    -sub-accounts ${SUB_ACCOUNTS:-5} \
    -transactions ${TRANSACTIONS:-100} \
    -url ${GNO_URL:-"http://gnoland:26657"} \
    -mnemonic "${MNEMONIC:-source bonus chronic canvas draft south burst lottery vacant surface solve popular case indicate oppose farm nothing bullet exhibit title speed wink action roast}" \
    -output ${OUTPUT_FILE:-"result.json"} \
    -mode PACKAGE_DEPLOYMENT
