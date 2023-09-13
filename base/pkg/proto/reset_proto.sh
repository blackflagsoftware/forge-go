#! /bin/zsh
cd .. && protoc --go_out=./pkg/proto --go-grpc_out=./pkg/proto ./pkg/proto/finance.proto