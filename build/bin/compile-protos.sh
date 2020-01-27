#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src"

# FIXME: Temporarily pull the proto files from xran repo
rm -f api/nb/*.proto api/sb/*.proto
wget https://raw.githubusercontent.com/OpenNetworkingFoundation/xranc/master/protos/c1-interface.proto -P api/nb/
wget https://raw.githubusercontent.com/OpenNetworkingFoundation/xranc/master/protos/e2-interface.proto -P api/sb/

protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/nb,plugins=grpc:. api/nb/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/sb,plugins=grpc:. api/sb/*.proto
