FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git build-base

COPY . /storj
WORKDIR /storj

RUN go get -d -v ./cmd/certificates
RUN go build -o /go/bin/certificates ./cmd/certificates

#----------------------
#FROM scratch
#
#COPY --from=builder /go/bin/certificates /certificates
#
#ENTRYPOINT ["/certificates"]
#
# #--> standard_init_linux.go:190: exec user process caused "no such file or directory"