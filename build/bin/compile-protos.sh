#!/bin/sh

proto_imports=".:${GOPATH}/src/github.com/gogo/protobuf/protobuf:${GOPATH}/src/github.com/gogo/protobuf:${GOPATH}/src"

protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,c1-interface.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/nb,plugins=grpc:. api/nb/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,e2-interface.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/sb,plugins=grpc:. api/sb/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,e2sm.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/sb/e2sm,plugins=grpc:. api/sb/e2sm/*.proto
protoc -I=$proto_imports --doc_out=docs/api --doc_opt=markdown,e2ap.md --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/sb/e2ap,plugins=grpc:. api/sb/e2ap/*.proto

protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/store/indications,plugins=grpc:. api/store/indications/*.proto
protoc -I=$proto_imports --gogofaster_out=Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,import_path=ran/store/requests,plugins=grpc:. api/store/requests/*.proto
