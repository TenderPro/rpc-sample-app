
FROM golang:1.14.0-alpine3.11

# Speed up build if proxy given
ARG GOPROXY
RUN echo $GOPROXY

WORKDIR /opt/app
RUN apk --update add curl git

# Cached layer
COPY ./go.mod ./go.sum ./
RUN go mod download

# Sources dependend layer
COPY ./ ./
RUN GOOS=linux go build -ldflags "-X main.version=`git describe --tags --always`" -a -o app .

FROM alpine:3.11.2

ENV DOCKERFILE_VERSION  200423

WORKDIR /opt/app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/app/app /usr/bin/app

ENTRYPOINT ["/usr/bin/app"]
