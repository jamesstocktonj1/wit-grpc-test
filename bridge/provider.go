package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strings"
	helloservice "wit-grpc-test/bridge/bindings/local/wit_grpc_test/hello_service"
	"wit-grpc-test/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.wasmcloud.dev/provider"
	"google.golang.org/grpc"
)

var (
	errComponentAlreadyLinked = errors.New("error: component already linked")
	errComponentNotLinked     = errors.New("error: component not linked")
)

type Bridge struct {
	proto.UnimplementedHelloServiceServer
	provider        *provider.WasmcloudProvider
	listen          net.Listener
	server          *grpc.Server
	linkedComponent string
}

func (b *Bridge) Init() (err error) {
	b.listen, err = net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}

	b.server = grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	proto.RegisterHelloServiceServer(b.server, b)

	return nil
}

func (b *Bridge) Serve() error {
	return b.server.Serve(b.listen)
}

func (b *Bridge) Stop() {
	b.server.GracefulStop()
}

func (b *Bridge) healthcheck() string {
	resp := provider.HealthCheckResponse{
		Healthy: true,
		Message: "healthy",
	}
	msg, err := json.Marshal(&resp)
	if err != nil {
		return "unhealthy"
	}
	return string(msg)
}

func (b *Bridge) createLink(link provider.InterfaceLinkDefinition) error {
	b.provider.Logger.Info("createLink", "link", link)
	if len(b.linkedComponent) > 0 {
		return errComponentAlreadyLinked
	}
	b.linkedComponent = link.SourceID
	return nil
}

func (b *Bridge) deleteLink(link provider.InterfaceLinkDefinition) error {
	b.provider.Logger.Info("deleteLink", "link", link)
	if strings.Compare(b.linkedComponent, link.SourceID) == 0 {
		b.linkedComponent = ""
	}
	return nil
}

func (b *Bridge) Hello(_ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	ctx, span := tracer.Start(_ctx, "Hello")
	defer span.End()

	if len(b.linkedComponent) == 0 {
		span.RecordError(errComponentNotLinked)
		return nil, errComponentNotLinked
	}

	b.provider.Logger.Info("incoming message", "request", req)
	resp, err := helloservice.Hello(
		injectTraceHeader(ctx),
		b.provider.OutgoingRpcClient(b.linkedComponent),
		mapRequest(req),
	)

	b.provider.Logger.Info("outgoing message", "response", *resp)
	return mapResponse(resp), err
}
