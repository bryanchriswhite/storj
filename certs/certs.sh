#!/usr/env/bin bash
set -o errexit

git clone https://github.com/storj/storj
cd storj

commands="{./cmd/certificates,storagenode}"

# might be cachable as separate dockerfile command
go get -d -v ${commands}
go install ${commands}
