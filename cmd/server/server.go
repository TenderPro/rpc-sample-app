package main

import (
	"errors"
	"net"

	"google.golang.org/grpc"

	"github.com/go-pg/pg"

	"github.com/jessevdk/go-flags"

	mapper "github.com/birkirb/loggers-mapper-logrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/birkirb/loggers.v1"

	api "companyserv/grpccompanyserv"
)

// DBConfig holds cli part of pg.Options
type DBConfig struct {
	Addr     string `long:"addr"  default:"localhost:5432" description:"host:port"`
	User     string `long:"user" description:"User name"`
	Password string `long:"password" description:"User password"`
	Database string `long:"name" description:"Database name"`
	// ApplicationName is the application name. Used in logs on Pg side.
	// Only available from pg-9.0.
	ApplicationName string `long:"app" default:"companyserv" description:"Application name"`
}

// Config holds all config vars
type Config struct {
	Addr string `long:"addr" default:"localhost:7070"  description:"Listen address"`
	// TODO:testing	IsSingle    bool   `long:"single" description:"Run service in single transaction"`
	IsDebugging bool `long:"debug" description:"Print debug logs"`

	API api.Config `group:"API Options" namespace:"api"`
	DB  DBConfig   `group:"DB Options" namespace:"db"`
}

var (
	// ErrGotHelp returned after showing requested help
	ErrGotHelp = errors.New("help printed")
	// ErrBadArgs returned after showing command args error message
	ErrBadArgs = errors.New("option error printed")
)

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
	if cfg.IsDebugging {
		l.SetLevel(logrus.DebugLevel)
		l.SetReportCaller(true)
	}
	return &mapper.Logger{Logger: l} // Same as mapper.NewLogger(l) but without info log message
}

// setupRouter creates gin router
func serve(cfg *Config, log loggers.Contextual) {

	dbcfg := &pg.Options{
		Addr:            cfg.DB.Addr,
		User:            cfg.DB.User,
		Password:        cfg.DB.Password,
		Database:        cfg.DB.Database,
		ApplicationName: cfg.DB.ApplicationName,
	}
	log.Debugf("%+v", dbcfg)
	db := pg.Connect(dbcfg)
	defer db.Close()

	// create a listener on TCP port
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// create a server instance
	s := api.NewCompanyServer(&cfg.API, log, db)
	// create a gRPC server object
	grpcServer := grpc.NewServer()
	// attach the Ping service to the server
	api.RegisterCompanyServiceServer(grpcServer, s)
	// start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
