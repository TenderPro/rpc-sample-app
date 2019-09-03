package server

import (
	"errors"

	"github.com/jessevdk/go-flags"

	mapper "github.com/birkirb/loggers-mapper-logrus"
	"github.com/sirupsen/logrus"
	"gopkg.in/birkirb/loggers.v1"

	"SELF/api"
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
	Addr       string `long:"addr" default:"localhost:7070"  description:"Listen address"`
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
func setupLog(cfg *Config) (loggers.Contextual, error) {
	l := logrus.New()
	ll, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	if ll == logrus.DebugLevel {
		l.SetReportCaller(true)
	}
	l.SetLevel(ll)
	return &mapper.Logger{Logger: l}, nil // Same as mapper.NewLogger(l) but without info log message
}
