package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jinzhu/gorm"

	mapper "github.com/birkirb/loggers-mapper-logrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/birkirb/loggers.v1"

	"SELF/api"
	"SELF/api/pb"
)

// код основной функции с поддержкой тестов
func Run(version string, exitFunc func(code int)) {
	var err error
	var cfg *Config
	defer func() { shutdown(exitFunc, err) }()
	cfg, err = setupConfig()
	if err != nil {
		return
	}
	log.Printf("gRPC sample %s. gRPC sample service", version)
	l, err := setupLog(cfg)
	if err != nil {
		log.Println(err.Error())
	} else {
		serve(cfg, l)
	}
}

// exit after deferred cleanups have run
func shutdown(exitFunc func(code int), e error) {
	if e != nil {
		var code int
		switch e {
		case ErrGotHelp:
			code = 3
		case ErrBadArgs:
			code = 2
		default:
			code = 1
			log.Printf("Run error: %s", e.Error())
		}
		exitFunc(code)
	}
}

// serve creates and starts service
func serve(cfg *Config, log loggers.Contextual) {
	url := fmt.Sprintf("%s://%s:%s@%s/%s?%s",
		cfg.DB.Driver,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Addr,
		cfg.DB.Database,
		cfg.DB.Options,
	)
	log.Debugf("Connect: %s", url)
	db, err := gorm.Open(cfg.DB.Driver, url)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer db.Close()
	if cfg.LogLevel == "debug" {
		db = db.Debug()
	}
	db.SetLogger(log)
	// create a listener on TCP port
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc_recovery.Option{
		// go-grpc-middleware@v1.0.0 does not have WithRecoveryHandlerContext
		grpc_recovery.WithRecoveryHandlerContext(func(ctx context.Context, p interface{}) error {
			err := fmt.Errorf("PANIC: %s", p)
			return err
		}),
	}
	logrusEntry := logrus.NewEntry(log.(*mapper.Logger).Logger)
	// create a server instance
	s := api.NewServer(&cfg.API, log, db)
	// create a gRPC server object
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_prometheus.StreamServerInterceptor,
			grpc_logrus.StreamServerInterceptor(logrusEntry),
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)
	// attach the service to the server
	pb.RegisterGreeterServer(grpcServer, s)

	grpc_prometheus.Register(grpcServer)
	http.Handle(cfg.MetricURL, promhttp.Handler())
	go func() {
		err = http.ListenAndServe(cfg.MetricAddr, nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	}()
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
