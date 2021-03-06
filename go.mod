module github.com/TenderPro/rpc-sample-app

go 1.14

replace SELF => ./

replace github.com/TenderPro/rpckit => ../rpckit

require (
	SELF v0.0.0-00010101000000-000000000000
	github.com/TenderPro/rpckit v0.1.1
	github.com/UNO-SOFT/grpcer v0.4.5
	github.com/apisite/apitpl v0.3.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.0
	github.com/gorilla/mux v1.7.4
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.14.4
	github.com/kevinburke/go-bindata v3.21.0+incompatible
	github.com/mwitkow/go-proto-validators v0.3.0
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-rpc/nrpc v0.0.0-20190920042445-3ae2c6c6aad4
	github.com/opentracing/opentracing-go v1.1.0
	github.com/phogolabs/parcello v0.8.1
	github.com/prometheus/client_golang v1.5.1
	github.com/stretchr/testify v1.5.1
	github.com/tiaguinho/gosoap v1.2.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc
	github.com/uber/jaeger-client-go v2.23.0+incompatible
	go.uber.org/zap v1.15.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543
	google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215
	google.golang.org/grpc v1.29.1
)
