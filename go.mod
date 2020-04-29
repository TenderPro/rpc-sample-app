module github.com/TenderPro/rpc-sample-app

go 1.14

replace (
	SELF => ./
	github.com/TenderPro/rpckit => ./../rpckit
)

require (
	SELF v0.0.0-00010101000000-000000000000
	github.com/TenderPro/rpckit v0.0.0-00010101000000-000000000000
	github.com/UNO-SOFT/grpcer v0.4.5
	github.com/UNO-SOFT/soap-proxy v0.8.3 // indirect
	github.com/birkirb/loggers-mapper-logrus v0.0.0-20180326232643-461f2d8e6f72
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/envoyproxy/protoc-gen-validate v0.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/gddo v0.0.0-20200324184333-3c2cc9a6329d // indirect
	github.com/golang/protobuf v1.4.0
	github.com/golangci/golangci-lint v1.25.0 // by hand
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.4
	github.com/jackc/pgx/v4 v4.6.0
	github.com/jessevdk/go-flags v1.4.0
	github.com/jinzhu/gorm v1.9.10
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/nats-io/gnatsd v1.4.1
	github.com/nats-io/go-nats v1.7.2
	github.com/nats-io/nats-server/v2 v2.1.6 // indirect
	github.com/nats-io/nats.go v1.9.2
	github.com/nats-rpc/nrpc v0.0.0-20190920042445-3ae2c6c6aad4
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pgmig/pgmig v0.35.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/rakyll/statik v0.1.7
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc
	github.com/uber/jaeger-client-go v2.23.0+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215
	google.golang.org/grpc v1.29.1
	gopkg.in/birkirb/loggers.v1 v1.1.0
)
