generate:
	protoc greetpb/greet.proto --go_out=plugins=grpc:server
	protoc greetpb/greet.proto --go_out=plugins=grpc:client
