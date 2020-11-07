generate:
	protoc greetpb/greet.proto --go_out=plugins=grpc:server
	protoc greetpb/greet.proto --go_out=plugins=grpc:client
certificates:
	sudo chmod +x ./generateSSL.sh
	./generateSSL.sh
