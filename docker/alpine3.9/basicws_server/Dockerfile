FROM alpine:3.9

COPY build/alpine3.9/basicws_server /usr/bin
RUN chmod +x /usr/bin/basicws_server

RUN apk add --update --no-cache ca-certificates fuse

# ENTRYPOINT ["/usr/bin/basicws_server"]
