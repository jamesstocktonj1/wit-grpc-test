//go:generate go run github.com/bytecodealliance/wasm-tools-go/cmd/wit-bindgen-go generate --world hello --out gen ./wit
package main

import (
	"fmt"
	helloservice "wit-grpc-test/hello/gen/local/wit-grpc-test/hello-service"
)

func init() {
	helloservice.Exports.Hello = Hello
}

func Hello(req helloservice.HelloRequest) helloservice.HelloResponse {
	return helloservice.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Message),
	}
}

func main() {}
