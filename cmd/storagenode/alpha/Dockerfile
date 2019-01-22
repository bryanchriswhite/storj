FROM alpine

EXPOSE 28967

ARG binary

ENV STORJ_KADEMLIA_BOOTSTRAP_ADDR="bootstrap.storj.io:8888"
ENV STORJ_METRICS_APP_SUFFIX="-alpha"
ENV STORJ_METRICS_INTERVAL="30m"
ENV STORJ_SERVER_USE_PEER_CA_WHITELIST="true"

WORKDIR /app
RUN  mkdir ./{config,identity}
COPY  $binary ./storagenode

ENTRYPOINT ./storagenode $COMMAND --config-dir="config" --identity-dir="identity"