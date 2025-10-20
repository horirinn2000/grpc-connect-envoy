module github.com/horirinn2000/grpc-connect-envoy/server

go 1.25.3

require github.com/horirinn2000/grpc-connect-envoy/proto v0.0.0

require (
	connectrpc.com/connect v1.19.1 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

replace github.com/horirinn2000/grpc-connect-envoy/proto => ../proto
