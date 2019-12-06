PLATFORMS := linux/amd64 linux/arm linux/arm64 windows/amd64 darwin/amd64

install:
	# More info: https://github.com/grpc-ecosystem/grpc-gateway
	go get -u \
		github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
		github.com/golang/protobuf/protoc-gen-go

generate:
	# Generate Go and gRPC-Gateway output.
	#
	# --go_out generates Go protobuf output with gRPC plugin enabled.
	# --grpc-gateway_out generates gRPC-Gateway output.
	# proto/commands.proto is the location of the protofile in use.
	#
	# More info: https://github.com/grpc-ecosystem/grpc-gateway
	protoc \
		-I proto \
		-I ${GOPATH}\src\github.com\grpc-ecosystem\grpc-gateway\third_party\googleapis \
		--go_out=plugins=grpc,paths=source_relative:./proto \
		--grpc-gateway_out=./proto \
		proto/commands.proto

setup: install generate

build: ./build.sh