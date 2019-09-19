FROM golang:1.12-alpine3.9

COPY azgoapi /usr/bin
RUN chmod +x /usr/bin/azgoapi

RUN apk add --update --no-cache ca-certificates

WORKDIR /usr/azgoapi

# ENTRYPOINT ["/usr/bin/azgoapi"]
CMD ["/usr/bin/azgoapi"]
