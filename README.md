# grpc-gateway-test

Firestore backend connector with grpc-gateway


### How to generate codes from proto files

```
- protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1/third_party/googleapis -I$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.1/protobuf --go_out=plugins=grpc:. model/*.proto
- protoc -I/usr/local/include -I. -I$GOPATH/src -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.12.1/third_party/googleapis -I$GOPATH/pkg/mod/github.com/gogo/protobuf@v1.3.1/protobuf --grpc-gateway_out=logtostderr=true:. model/*.proto
```
