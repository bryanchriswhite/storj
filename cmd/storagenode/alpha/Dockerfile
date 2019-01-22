# build
FROM golang:1.11 as build-env

ARG command="run"
ARG go_os="linux"
ARG go_arch=""

ENV GOOS=$go_os
ENV GOARCH=$go_arch

ENV GOARM=6
ENV CGO_ENABLED=1
ENV CC=arm-linux-gnueabihf-gcc

ADD . /go/src/storj.io/storj
WORKDIR /go/src/storj.io/storj/cmd/storagenode

# dependencies + binary
RUN apt update && apt install -y git gcc gcc-arm-linux-gnueabihf musl-dev && unset GOPATH && go mod vendor && go build -a -installsuffix cgo -o storagenode . && mkdir config && mkdir identity

# final stage
FROM alpine

EXPOSE 28967

ENV STORJ_KADEMLIA_BOOTSTRAP_ADDR="bootstrap.storj.io:8888"
ENV STORJ_METRICS_APP_SUFFIX="-alpha"
ENV STORJ_METRICS_INTERVAL="30m"
ENV STORJ_SERVER_USE_PEER_CA_WHITELIST="true"

WORKDIR /app
COPY --from=build-env /go/src/storj.io/storj/cmd/storagenode/storagenode /app/
COPY --from=build-env /go/src/storj.io/storj/cmd/storagenode/config /app/
COPY --from=build-env /go/src/storj.io/storj/cmd/storagenode/identity /app/

ENTRYPOINT ./storagenode $command --config-dir="config" --identity-dir="identity"