# FROM golang:1.12-alpine3.9
FROM alpine:3.9

COPY build/alpine3.9/azgoapi /usr/bin
RUN chmod +x /usr/bin/azgoapi

RUN apk add --update --no-cache ca-certificates

WORKDIR /usr/azgoapi

# ENTRYPOINT ["/usr/bin/azgoapi"]
CMD ["/usr/bin/azgoapi"]
