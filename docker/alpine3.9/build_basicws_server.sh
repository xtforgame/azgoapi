
# from https://github.com/restic/restic/tree/master/docker
#!/bin/sh

set -e

echo "Build binary using golang docker image"
docker run --rm -ti \
  -v $(pwd):/go/src/github.com/xtforgame/gbol \
  -w /go/src/github.com/xtforgame/gbol \
  -e CGO_ENABLED=1 \
  -e GOOS=linux \
  -e GO111MODULE=on \
  golang:1.12-alpine3.9 go build -mod=vendor -o ./build/alpine3.9/basicws_server main/basicws_server.go

echo "Build docker image xtforgame/basicws_server:latest"
docker build --rm -t xtforgame/basicws_server:latest -f docker/alpine3.9/basicws_server/Dockerfile .

# docker run --rm -ti \
#   -p 8080:8080 \
#   -v $(pwd)/tmp:/usr/gbol \
#   -w /usr/gbol \
#   xtforgame/basicws_server:latest gbol ./forweb ./pgbackrest-backup ./output
