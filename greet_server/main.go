package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"greet_server/greetpb"

	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Println("Invoked with: ", req)

	firstname := req.GetGreeting().GetFirstName()
	lastname := req.GetGreeting().GetLastName()
	result := "Hello " + firstname + " " + lastname

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("Invoked with: ", req)
	firstname := req.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + firstname + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	request := 0
	result := " "

	for {
		request++
		log.Printf("Request: %v", request)
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Print("Client stop the streaming")
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		firstName := msg.GetGreeting().GetFirstName()
		result += firstName + "! "
	}
}

func (s *server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone is invoked")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("Client stop the streaming")
			return nil
		}
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName
		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error: %v", sendErr)
		}

	}
}

func main() {
	fmt.Println("Hello gRPC")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v", err)
	}
}
