package main

import (
	"context"
	"net"
	helloservice "wit-grpc-test/bridge/bindings/local/wit_grpc_test/hello_service"
	"wit-grpc-test/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.wasmcloud.dev/provider"
	"google.golang.org/grpc"
)

type Bridge struct {
	proto.UnimplementedHelloServiceServer
	provider *provider.WasmcloudProvider
	listen   net.Listener
	server   *grpc.Server
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

func (b *Bridge) Hello(_ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	ctx, span := tracer.Start(_ctx, "Hello")
	defer span.End()

	b.provider.Logger.Info("incoming message", "request", req)
	resp, err := helloservice.Hello(
		injectTraceHeader(ctx),
		b.provider.OutgoingRpcClient("wit_grpc_test-hello"),
		mapRequest(req),
	)

	b.provider.Logger.Info("outgoing message", "response", *resp)
	return mapResponse(resp), err
}

func mapRequest(req *proto.HelloRequest) *helloservice.HelloRequest {
	return &helloservice.HelloRequest{
		Message: req.GetMessage(),
	}
}

func mapResponse(resp *helloservice.HelloResponse) *proto.HelloResponse {
	return &proto.HelloResponse{
		Message: resp.Message,
	}
}
