// main application file, see README.md
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/jessevdk/go-flags"

	mapper "github.com/birkirb/loggers-mapper-logrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/birkirb/loggers.v1"

	"SELF/api"
	"SELF/api/pb"
)

// DBConfig holds cli part of pg.Options
type DBConfig struct {
	Addr     string `long:"addr"  default:"localhost:5432" description:"host:port"`
	Driver   string `long:"driver" default:"postgres" description:"DB driver"`
	User     string `long:"user" description:"User name"`
	Password string `long:"password" description:"User password"`
	Database string `long:"name" description:"Database name"`
	Options  string `long:"opts" default:"sslmode=disable" description:"Database connect options"`
}

// Config holds all config vars
type Config struct {
	Addr string `long:"addr" default:"localhost:7070"  description:"Listen address"`
	MetricAddr string `long:"metric_addr"  default:"localhost:8080" description:"prometheus service host:port"`
	MetricURL  string `long:"metric_url"  default:"/metrics" description:"prometheus service URL"`
	LogLevel   string `long:"log_level" default:"debug" description:"Log level"`

	API api.Config `group:"API Options" namespace:"api"`
	DB  DBConfig   `group:"DB Options" namespace:"db"`
}

var (
	// ErrGotHelp returned after showing requested help
	ErrGotHelp = errors.New("help printed")
	// ErrBadArgs returned after showing command args error message
	ErrBadArgs = errors.New("option error printed")
)

// Actual version value will be set at build time
var version = "0.0-dev"

// main может быть вынесена в отдельный фал и исключена из расчета code coverage
func main() {
	run(os.Exit)
}

// код основной функции с поддержкой тестов
func run(exitFunc func(code int)) {
	log.Printf("gRPC sample %s. gRPC sample service", version)
	var err error
	var cfg *Config
	defer func() { shutdown(exitFunc, err) }()
	cfg, err = setupConfig()
	if err != nil {
		return
	}
	l := setupLog(cfg)
	serve(cfg, l)
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

// setupConfig loads flags from args (if given) or command flags and ENV otherwise
func setupConfig(args ...string) (*Config, error) {
	cfg := &Config{}
	p := flags.NewParser(cfg, flags.Default)
	var err error
	if len(args) == 0 {
		_, err = p.Parse()
	} else {
		_, err = p.ParseArgs(args)
	}
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return nil, ErrGotHelp
		}
		return nil, ErrBadArgs
	}
	return cfg, nil
}

// setupLog creates logger
func setupLog(cfg *Config) loggers.Contextual {
	l := logrus.New()
	ll, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		l.Fatal(err)
	}
	if ll == logrus.DebugLevel {
		l.SetReportCaller(true)
	}
	l.SetLevel(ll)
	return &mapper.Logger{Logger: l} // Same as mapper.NewLogger(l) but without info log message
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
	db.LogMode(cfg.LogLevel == "debug")
	db.SetLogger(gorm.Logger{LogWriter: log})
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
