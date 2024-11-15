ghz -c 100 -n 100000 --insecure \
    --proto proto/hello.proto \
    --call witgrpctest.HelloService/Hello \
    -d '{ "message": "World" }' \
    localhost:8080