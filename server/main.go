package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"greet_server/greetpb"
)

type server struct{}

func main() {
	fmt.Println("Hello from gRPC's Server")

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

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Println("Invoked with: ", req)

	// This is just
	// to test client deadlines
	time.Sleep(25 * time.Second)

	// This if-statement
	// handler whether the client
	// has cancel the request
	if ctx.Err() == context.Canceled {
		fmt.Println("Client has canceled the request")
		return nil, status.Errorf(codes.Canceled, "The Client has canceled the request")
	}

	firstname := req.GetGreeting().GetFirstName()
	lastname := req.GetGreeting().GetLastName()
	result := "Hello " + firstname + " " + lastname

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Println("Invoked with: ", req)

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

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log.Println("Long Greet Invoked")

	var request int
	var result string

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

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	log.Println("GreetEveryone is invoked")

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

// Unary rpc
// but with error handling
func (*server) SquareRoot(ctx context.Context, req *greetpb.SquareRootRequest) (*greetpb.SquareRootResponse, error) {
	log.Println("SquareRoot invoked")

	inputNumber := req.GetNumber()

	if inputNumber < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", inputNumber),
		)
	}

	return &greetpb.SquareRootResponse{
		Root: math.Sqrt(float64(inputNumber)),
	}, nil
}
