// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// Code based on: https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/testing/pingservice.go

// Package service provides sample application logic
package service

import (
	"context"
	"time"

	"github.com/TenderPro/rpckit/app/ticker"
	"github.com/nats-io/nats.go"
	opentracing "github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"SELF/pkg/pb"
)

const (
	// DefaultResponseValue is the default value used in PingEmpty.
	DefaultResponseValue = "default_response_value"
	// ListResponseCount is the expected number of responses to PingList
	ListResponseCount = 50
)

type Config struct {
}

type Service struct {
	cfg        Config
	subject    string
	spanPrefix string
	ticker     *ticker.Service
	log        *zap.Logger
}

func New(cfg Config, log *zap.Logger, nc *nats.Conn, subject, prefix string) *Service {
	tick := ticker.New(log, nc, subject)
	return &Service{cfg: cfg, log: log, subject: subject, spanPrefix: prefix, ticker: tick}
}

func (srv Service) PingEmpty(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	srv.log.Debug("PingEmpty")
	span, _ := opentracing.StartSpanFromContext(ctx, srv.spanPrefix+"/PingEmpty")
	if span != nil {
		srv.log.Debug("InSpan")
		span.LogKV("event", "service.ping.empty")
		span.Finish()
	}
	return &pb.PingResponse{Value: DefaultResponseValue, Counter: 43}, nil
}

func (srv Service) Ping(ctx context.Context, ping *pb.PingRequest) (*pb.PingResponse, error) {
	srv.log.Debug("Ping")
	span, _ := opentracing.StartSpanFromContext(ctx, srv.spanPrefix+"/Ping")
	if span != nil {
		srv.log.Debug("InSpan")
		span.LogKV("event", "service.ping")
	}
	// Simulate slow call
	time.Sleep(time.Second * 2)
	if span != nil {
		span.Finish()
	}
	return &pb.PingResponse{Value: ping.Value, Counter: 42}, nil
}

func (srv Service) PingError(ctx context.Context, ping *pb.PingRequest) (*pb.Empty, error) {
	srv.log.Debug("PingError")
	code := codes.Code(ping.ErrorCodeReturned)
	span, _ := opentracing.StartSpanFromContext(ctx, srv.spanPrefix+"/PingError")
	if span != nil {
		srv.log.Debug("InSpan")
		span.LogKV("event", "service.ping.error", "code", code)
		span.Finish()
	}
	return nil, status.Error(code, "Userspace error.")
}

func (srv Service) PingList(ping *pb.PingRequest, stream pb.PingService_PingListServer) error {
	srv.log.Debug("PingList")
	if ping.ErrorCodeReturned != 0 {
		return grpc.Errorf(codes.Code(ping.ErrorCodeReturned), "foobar")
	}
	/* TODO: fetch context in STREAM methods
	span, _ := opentracing.StartSpanFromContext(ctx, srv.spanPrefix+"/PingList")
	if span != nil {
		srv.log.Debug("InSpan")
		span.LogKV("event", "service.ping.list")
	}
	*/
	// Send user trailers and headers.
	for i := 0; i < ListResponseCount; i++ {
		//	span.LogKV("loop", i)
		stream.Send(&pb.PingResponse{Value: ping.Value, Counter: int32(i)})
	}
	/*
		if span != nil {
			span.Finish()
		}
	*/
	return nil
}

/*
func (srv Service) PingStream(stream pb.PingService_PingStreamServer) error {
	srv.log.Debug("PingStream")
	count := 0
	for true {
		ping, err := stream.Recv()
		//		span.LogKV("loop", count)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		stream.Send(&pb.PingResponse{Value: ping.Value, Counter: int32(count)})
		count += 1
	}
	return nil
}
*/

// TimeService is a gRPC service for ticker
func (srv Service) TimeService(in *ticker.TimeRequest, stream pb.PingService_TimeServiceServer) error {
	return srv.ticker.TimeService(in, stream)
}
