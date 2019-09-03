
FROM golang:1.12.6-alpine3.9 as builder

WORKDIR /opt/app
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

WORKDIR /opt/app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/app/grpcsample /usr/bin/grpcsample

ENTRYPOINT ["/usr/bin/grpcsample"]
