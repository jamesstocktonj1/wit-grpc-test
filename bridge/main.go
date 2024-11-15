//go:generate wit-bindgen-wrpc go --out-dir bindings --package wit-grpc-test/bridge/bindings wit

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.wasmcloud.dev/provider"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// Init wasmcloud Provider
	p, err := provider.New()
	if err != nil {
		return err
	}

	// Init gRPC Server
	b := Bridge{
		provider: p,
	}
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
		// stopFunc()
		return err
	case <-signalCh:
		p.Shutdown()
		// stopFunc()
	}

	return nil
}
