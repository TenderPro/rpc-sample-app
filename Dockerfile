
FROM golang:1.12.6-alpine3.9 as builder

WORKDIR /opt/companyserv
RUN apk --update add curl git

# Cached layer
COPY ./go.mod ./go.sum ./
RUN go mod download

# Sources dependent layer
COPY ./ ./
RUN go build -o  grpcsample./cmd/server
# Для версии из git:
# RUN go build -o  grpcsample-ldflags "-X main.version=`git describe --tags`" ./cmd/server

FROM alpine:3.9

MAINTAINER Aleksei Kovrizhkin <lekovr@gmail.com>

ENV DOCKERFILE_VERSION  190730

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /opt/companyserv

COPY --from=builder /opt/companyserv/ grpcsample/usr/bin/companyserv

ENTRYPOINT ["/usr/bin/companyserv"]
