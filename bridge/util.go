package main

import (
	helloservice "wit-grpc-test/bridge/bindings/local/wit_grpc_test/hello_service"
	"wit-grpc-test/proto"
)

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
