
FROM golang:1.12.6-alpine3.9 as builder

WORKDIR /opt/companyserv
RUN apk --update add curl git

# Cached layer
COPY ./go.mod ./go.sum ./
RUN go mod download

# Sources dependent layer
COPY ./ ./
RUN go build .
# Для версии из git:
# RUN go build -ldflags "-X main.version=`git describe --tags`" .

FROM alpine:3.9

ENV DOCKERFILE_VERSION  190730

WORKDIR /opt/companyserv

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/companyserv/ grpcsample/usr/bin/companyserv

ENTRYPOINT ["/usr/bin/companyserv"]
