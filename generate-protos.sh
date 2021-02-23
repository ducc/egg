#!/bin/bash

rm -rf protos/
protoc -I=. --go_out=. --go-grpc_out=. protos.proto
cp -r github.com/ducc/egg/protos protos
rm -rf github.com
