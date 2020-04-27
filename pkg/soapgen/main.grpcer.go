// Generated with protoc-gen-grpcer
//	from "main.proto"
//	at   2020-04-27T08:41:44Z
//
// DO NOT EDIT!

package soapgen

import (
	"io"
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	errors "golang.org/x/xerrors"
	grpcer "github.com/UNO-SOFT/grpcer"

	pb "SELF/pkg/pb"
	"github.com/TenderPro/rpckit/app/ticker"
	
)



type client struct {
	pb.PingServiceClient
	m map[string]inputAndCall
}

func (c client) List() []string {
	names := make([]string, 0, len(c.m))
	for k := range c.m {
		names = append(names, k)
	}
	return names
}

func (c client) Input(name string) interface{} {
	iac := c.m[name]
	if iac.Input == nil {
		return nil
	}
	return iac.Input()
}

func (c client) Call(name string, ctx context.Context, in interface{}, opts ...grpc.CallOption) (grpcer.Receiver, error) {
	iac := c.m[name]
	if iac.Call == nil {
		return nil, errors.Errorf("name %q not found", name)
	}
	return iac.Call(ctx, in, opts...)
}
func NewClient(cc *grpc.ClientConn) grpcer.Client {
	c := pb.NewPingServiceClient(cc)
	return client{
		PingServiceClient: c,
		m: map[string]inputAndCall{
		"Ping": inputAndCall{
			Input: func() interface{} { return new(pb.PingRequest) },
			Call: func(ctx context.Context, in interface{}, opts ...grpc.CallOption) (grpcer.Receiver, error) {
				input := in.(*pb.PingRequest)
				res, err := c.Ping(ctx, input, opts...)
				if err != nil {
					return &onceRecv{Out:res}, err
				}
				return &onceRecv{Out:res}, err
				
			},
		},
		"PingEmpty": inputAndCall{
			Input: func() interface{} { return new(pb.Empty) },
			Call: func(ctx context.Context, in interface{}, opts ...grpc.CallOption) (grpcer.Receiver, error) {
				input := in.(*pb.Empty)
				res, err := c.PingEmpty(ctx, input, opts...)
				if err != nil {
					return &onceRecv{Out:res}, err
				}
				return &onceRecv{Out:res}, err
				
			},
		},
		"PingError": inputAndCall{
			Input: func() interface{} { return new(pb.PingRequest) },
			Call: func(ctx context.Context, in interface{}, opts ...grpc.CallOption) (grpcer.Receiver, error) {
				input := in.(*pb.PingRequest)
				res, err := c.PingError(ctx, input, opts...)
				if err != nil {
					return &onceRecv{Out:res}, err
				}
				return &onceRecv{Out:res}, err
				
			},
		},
		"PingList": inputAndCall{
			Input: func() interface{} { return new(pb.PingRequest) },
			Call: func(ctx context.Context, in interface{}, opts ...grpc.CallOption) (grpcer.Receiver, error) {
				input := in.(*pb.PingRequest)
				res, err := c.PingList(ctx, input, opts...)
				if err != nil {
					return &onceRecv{Out:res}, err
				}
				return multiRecv(func() (interface{}, error) { return res.Recv() }), nil
				
			},
		},
		"TimeService": inputAndCall{
			Input: func() interface{} { return new(ticker.TimeRequest) },
			Call: func(ctx context.Context, in interface{}, opts ...grpc.CallOption) (grpcer.Receiver, error) {
				input := in.(*ticker.TimeRequest)
				res, err := c.TimeService(ctx, input, opts...)
				if err != nil {
					return &onceRecv{Out:res}, err
				}
				return multiRecv(func() (interface{}, error) { return res.Recv() }), nil
				
			},
		},
		
		},
	}
}

type inputAndCall struct {
	Input func() interface{}
	Call func(ctx context.Context, in interface{}, opts ...grpc.CallOption) (grpcer.Receiver, error)
}

type onceRecv struct {
	Out interface{}
	done bool
}
func (o *onceRecv) Recv() (interface{}, error) {
	if o.done {
		return nil, io.EOF
	}
	out := o.Out
	o.done, o.Out = true, nil
	return out, nil
}

type multiRecv func() (interface{}, error)
func (m multiRecv) Recv() (interface{}, error) {
	return m()
}

var _ = multiRecv(nil) // against "unused"

