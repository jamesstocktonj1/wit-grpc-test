grpcurl -plaintext -proto proto/hello.proto \
    -d '{ "message": "World" }' \
    localhost:8080 witgrpctest.HelloService/Hello