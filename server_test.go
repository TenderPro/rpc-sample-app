// server_test - внешние тесты сервиса.
// Преварительная версия, которая
// * не очищает БД после внесения изменений
// * работает с БД, в которой выполнен sql/z_addon.sql
//
// Пример вызова:
//   APP_ADDR=localhost:7070 go test -v -count=1 .
package main_test

import (
	"context"
	//	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"google.golang.org/grpc"

	"SELF/api/pb"
)

type ServerSuite struct {
	suite.Suite
	srv pb.GreeterClient
	ctx context.Context
}

var (
	tick = time.Second * time.Duration(1)
)

func TestSuite(t *testing.T) {
	if os.Getenv("APP_ADDR") == "" {
		t.Skip("Skipping testing without app address")
	}

	ctx, cancel := context.WithTimeout(context.Background(), tick)
	defer cancel()
	conn, err := grpc.DialContext(ctx, os.Getenv("APP_ADDR"), grpc.WithInsecure(), grpc.WithBlock())
	assert.Nil(t, err)
	if err == context.DeadlineExceeded {
		// context deadline exceeded
		t.Skip("Service is not available")
	}
	defer conn.Close()
	rand.Seed(time.Now().UTC().UnixNano())
	c := pb.NewGreeterClient(conn)
	myTest := &ServerSuite{
		srv: c,
		ctx: ctx,
	}
	suite.Run(t, myTest)
}

func (ss *ServerSuite) Test01SayHelloBadLen() {
	req := pb.HelloRequest{
		Name: "xx",
	}
	r, err := ss.srv.SayHello(ss.ctx, &req)
	require.Nil(ss.T(), r)
	assert.Equal(ss.T(), "rpc error: code = Unknown desc = invalid HelloRequest.Name: value length must be at least 3 runes", err.Error())
}

func (ss *ServerSuite) Test01SayHelloOK() {
	req := pb.HelloRequest{
		Name: "xxx",
	}
	r, err := ss.srv.SayHello(ss.ctx, &req)
	require.Nil(ss.T(), err)
	assert.Equal(ss.T(), &pb.HelloReply{Message: "Hello, xxx!"}, r)
}
