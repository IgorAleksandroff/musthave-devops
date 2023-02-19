package proto

//go:generate protoc -I=. -I=${GOPATH}/protoc-gen-validate/validate --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import --validate_out=paths=import,lang=go:./ proto/server.proto
