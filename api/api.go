// Package api - реализация gRPC API сервиса pb/api.proto
package api

//go:generate protoc --go_out=plugins=grpc:. pb/api.proto
import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"gopkg.in/birkirb/loggers.v1"

	"SELF/api/pb"
)

// Config содержит настройки, которые могут быть изменены в аргументах сервиса
type Config struct {
}

// Server holds server methods implementation
type Server struct {
	pb.GreeterServer
	db  *gorm.DB
	cfg *Config
	log loggers.Contextual
}

// NewServer создает указатель на структуру Server
func NewServer(cfg *Config, log loggers.Contextual, db *gorm.DB) *Server {
	return &Server{cfg: cfg, log: log, db: db}
}

// SayHello says hello
func (srv *Server) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &pb.HelloReply{
		Message: fmt.Sprintf("Hello, %s!", r.Name),
	}, nil
}
