package app

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/nats-io/nats.go"
	"github.com/nats-rpc/nrpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tmc/grpc-websocket-proxy/wsproxy"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	go_grpc "google.golang.org/grpc"

	// Register gzip system wide
	_ "google.golang.org/grpc/encoding/gzip"

	"github.com/TenderPro/rpckit/app/ticker"
	"github.com/TenderPro/rpckit/config"
	"github.com/TenderPro/rpckit/debug"
	"github.com/TenderPro/rpckit/grpc"
	"github.com/TenderPro/rpckit/soap"
	"github.com/TenderPro/rpckit/static"

	"SELF/pkg/nrpcgen"
	"SELF/pkg/pb"
	"SELF/pkg/service"
	"SELF/pkg/soapgen"
	"SELF/pkg/staticgen"
)

// Config holds Application config data
type Config struct {
	config.Config

	TickerSubject string `long:"ticker_subject" default:"stream.time" description:"Ticker channel subject"`

	Service service.Config `group:"Service Options" namespace:"srv"`
	GRPC    grpc.Config    `group:"gRPC Options" namespace:"grpc"`
	Static  static.Config  `group:"Static files Options" namespace:"html"`
	SOAP    soap.Config    `group:"SOAP Options" namespace:"soap"`
	Trace   debug.Config   `group:"Trace Options" namespace:"trace"`
}

// Application holds Application name for logging
const Application = "RPC Sample app"

// Run func called from main()
func Run(version string, exitFunc func(code int)) {
	var err error
	defer func() { config.Close(exitFunc, err) }()

	var cfg Config
	if err = config.New(&cfg); err != nil {
		return
	}

	log := debug.NewLogger(cfg.Debug)
	defer log.Sync()

	log.Info(Application, zap.String("version", version))
	log.Debug("Config", zap.Reflect("cfg", cfg))

	tracer, closer, er := debug.New(cfg.Trace, log)
	if er != nil {
		err = er
		return
	}
	defer closer.Close()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var nc *nats.Conn
	for {
		nc, err = nats.Connect(cfg.MQ, nats.Timeout(5*time.Second))
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	var sub *nats.Subscription

	group, gctx := errgroup.WithContext(ctx)

	if cfg.UseServer() {
		log.Info("Start ticker")
		group.Go(func() error {
			return ticker.Run(gctx, log, nc, cfg.TickerSubject)
		})
		if cfg.UseNRPC() {
			log.Info("Start NATS handler")
			pool := nrpc.NewWorkerPool(context.Background(), 200, 5, 4*time.Second)
			h := nrpcgen.NewPingServiceConcurrentHandler(pool, nc, service.New(cfg.Service, log, nc, cfg.TickerSubject, cfg.Trace.Name))

			log.Info("NATS subscriber", zap.String("subject", h.Subject()))
			sub, err = nc.Subscribe(h.Subject(), h.Handler)
			if err != nil {
				return
			}
		}
	}

	if sub != nil {
		log.Info("NATS subscriber is Ready")
		defer sub.Unsubscribe()
	}

	if cfg.UseProxy() {
		grpcService := grpc.New(cfg.GRPC, log, tracer)
		if cfg.UseNRPC() {
			log.Info("Start gRPC NATS client")
			client := nrpcgen.NewPingServiceClient(nc, "insta")
			pb.RegisterPingServiceServer(grpcService.Server, client)
		} else {
			log.Info("Start gRPC Native service")
			pb.RegisterPingServiceServer(grpcService.Server, service.New(cfg.Service, log, nc, cfg.TickerSubject, cfg.Trace.Name))
		}
		group.Go(func() error {
			<-gctx.Done()
			log.Info("Stop gRPC service")
			grpcService.Shutdown()
			return gctx.Err()
		})
		group.Go(func() error {
			return grpcService.ListenAndServe()
		})

		// JSON Gateway

		// Register gRPC server endpoint
		// Note: Make sure the gRPC server is running properly and accessible
		gwm := runtime.NewServeMux()
		opts1 := []go_grpc.DialOption{go_grpc.WithInsecure()}
		err = pb.RegisterPingServiceHandlerFromEndpoint(ctx, gwm, cfg.GRPC.Bind, opts1)
		if err != nil {
			return
		}

		// HTTP
		mux := http.NewServeMux()

		// Static
		staticService := static.New(cfg.Static, log,
			staticgen.AssetNames,
			staticgen.Asset,
			staticgen.AssetInfo,
		)
		staticService.SetupRouter(mux)

		// API handler
		mux.Handle("/v1/", static.OptionsProxy(wsproxy.WebsocketProxy(gwm)))
		// Prom5s handler
		mux.Handle("/metrics", promhttp.Handler())

		// SOAP
		soapService, err := soap.New(cfg.SOAP, log, cfg.GRPC.Bind, cfg.OutsideHost)
		if err != nil {
			return
		}
		soapService.SetupRouter(mux, soapgen.NewClient(soapService.Server), soapgen.WSDLgzb64)

		httpServer := http.Server{
			Addr:    cfg.BindHTTP,
			Handler: mux, //http.HandlerFunc(handler),
		}
		// Start HTTP server (and proxy calls to gRPC server endpoint)
		log.Info("Start HTTP service")
		group.Go(func() error {
			<-gctx.Done()
			log.Info("Stop HTTP service")
			if er := httpServer.Shutdown(ctx); er != nil {
				log.Error("Could not gracefully shutdown the server", zap.Error(er))
			}
			return gctx.Err()
		})
		group.Go(func() error {
			if er := httpServer.ListenAndServe(); er != http.ErrServerClosed {
				return er
			}
			return nil
		})
	}
	log.Info("Server started")
	defer log.Info("Server exited")
	group.Go(func() error {
		er := config.WaitSignal(gctx)
		log.Info("Got Signal", zap.Error(er))
		return er
	})
	err = group.Wait()
	return
}
