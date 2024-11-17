//go:generate wit-bindgen-wrpc go --out-dir bindings --package wit-grpc-test/bridge/bindings wit

package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	slogmulti "github.com/samber/slog-multi"
	"go.opentelemetry.io/contrib/bridges/otelslog"

	"go.opentelemetry.io/otel"
	"go.wasmcloud.dev/provider"
)

var tracer = otel.Tracer("wit-grpc-bridge")

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	b := Bridge{}

	// Init wasmcloud Provider
	p, err := provider.New(
		provider.HealthCheck(b.healthcheck),
		provider.TargetLinkPut(b.createLink),
		provider.TargetLinkDel(b.deleteLink),
	)
	if err != nil {
		return err
	}

	// Forward logs to Otel
	if p.HostData().OtelConfig.EnableObservability || p.HostData().OtelConfig.EnableLogs {
		p.Logger = slog.New(slogmulti.Fanout(
			p.Logger.Handler(),
			otelslog.NewLogger("bridge").Handler(),
		))
	}

	// Init gRPC Server
	b.provider = p
	err = b.Init()
	if err != nil {
		return err
	}

	// Setup two channels to await RPC and control interface operations
	providerCh := make(chan error, 1)
	signalCh := make(chan os.Signal, 1)

	// Handle control interface operations
	go func() {
		err := p.Start()
		providerCh <- err
	}()

	go func() {
		err := b.Serve()
		providerCh <- err
	}()

	// Shutdown on SIGINT
	signal.Notify(signalCh, syscall.SIGINT)

	// Run provider until either a shutdown is requested or a SIGINT is received
	select {
	case err = <-providerCh:
		b.Stop()
		return err
	case <-signalCh:
		b.Stop()
		p.Shutdown()
	}

	return nil
}
