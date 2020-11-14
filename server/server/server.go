package server

import (
	"context"
	"fmt"
	"greet_server/greetpb"
	"io"
	"log"
	"math"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server struct
type Server struct {
	Network string
	Address string
}

// NewServerInstance function
func NewServerInstance() *Server {
	return &Server{
		Network: "tcp",
		Address: "localhost:50051",
	}
}

// Greet function
// example of unary
func (*Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	log.Println("Invoked with: ", req)

	// This is just
	// to test client deadlines
	// time.Sleep(25 * time.Second)

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

// GreetManyTimes is
// example of server streaming
func (*Server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
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

// LongGreet function is
// example of client streaming
func (*Server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
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

// GreetEveryone function is
// example of Bi Directional Streaming
// Client and Server Streaming
func (*Server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
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

// SquareRoot function
// this example of error handling
func (*Server) SquareRoot(ctx context.Context, req *greetpb.SquareRootRequest) (*greetpb.SquareRootResponse, error) {
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
