#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src"
grpc_gateway_out=logtostderr=true,grpc_api_configuration=api/nb/a1/a1_config.yaml:.


## c1 related protobuf files
protoc -I=$proto_imports  --doc_out=docs/api --doc_opt=markdown,c1-interface.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/nb,plugins=grpc:. api/nb/*.proto

## e2 related protobuf files
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,e2-interface.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/sb,plugins=grpc:. api/sb/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,e2sm.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/sb/e2sm,plugins=grpc:. api/sb/e2sm/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,e2ap.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/sb/e2ap,plugins=grpc:. api/sb/e2ap/*.proto

## a1 related protobuf files
protoc -I=$proto_imports --grpc-gateway_out=$grpc_gateway_out  --doc_out=docs/api --doc_opt=markdown,a1-interface.md     --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/nb/a1,plugins=grpc:. api/nb/a1/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,a1-interface-types.md  --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/nb/a1/types,plugins=grpc:. api/nb/a1/types/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,a1-interface-tsp.md  --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/nb/a1/tsp,plugins=grpc:. api/nb/a1/a1-p/tsp/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,a1-interface-qos.md  --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/nb/a1/qos,plugins=grpc:. api/nb/a1/a1-p/qos/*.proto
