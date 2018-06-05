# build
FROM golang AS build-env
ADD . /go/src/storj.io/storj
RUN cd /go/src/storj.io/storj/cmd/overlay && go build -o overlay


# final stage
FROM golang
WORKDIR /app
COPY --from=build-env /go/src/storj.io/storj/cmd/overlay/overlay /app/

EXPOSE 8080
ENTRYPOINT ./overlay -tlsHosts=${TLS_HOSTS} -tlsCertPath=${TLS_CERT_PATH} -tlsKeyPath=${TLS_KEY_PATH} -redisAddress=${REDIS_ADDRESS} -redisPassword=${REDIS_PASSWORD} -db=${REDIS_DB} -srvPort=${OVERLAY_PORT} -httpPort=${HTTP_PORT}