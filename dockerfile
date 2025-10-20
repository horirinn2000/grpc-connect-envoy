FROM golang

ENV CGO_ENABLED=0

# formatter
RUN apt update
RUN apt-get install protobuf-compiler -y
RUN apt-get install clang-format -y

# connect
RUN go install github.com/bufbuild/buf/cmd/buf@latest
RUN go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

WORKDIR /app
