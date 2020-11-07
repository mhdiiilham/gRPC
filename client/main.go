package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"greet_client/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Greet Client")

	certFile := "../ssl/server.crt"
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "localhost")
	if sslErr != nil {
		log.Fatalf("SSL Error. Error: %v", sslErr)
	}
	opts := grpc.WithTransportCredentials(creds)
	conn, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Cannot connect. Error: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	doUnary(c, 30*time.Second) // This should be success
	// doUnary(c, 15*time.Second) // This should be fail
	// doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)
	// doSquareRoot(c)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming")

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilham",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nando",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Faway",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ester C",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hanafi",
			},
		},
	}

	// Create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	waitc := make(chan struct{})
	// Send bunch of messages to the client
	go func() {
		for _, req := range requests {
			fmt.Println("Sending ", req.GetGreeting().GetFirstName())
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive bunch of messages to the client
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				log.Print("Server stop the streaming")
				break
			}
			if err != nil {
				log.Fatalf("Error: %v", err)
				break
			}
			fmt.Printf("Receiving from server: %v\n", res.GetResult())
		}
		close(waitc)

	}()
	<-waitc
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilham",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Nando",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Faway",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ester C",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Hanafi",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req.GetGreeting().GetFirstName())
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Response from server: %v", res.GetResult())

}

func doUnary(c greetpb.GreetServiceClient, timeout time.Duration) {

	greetRequest := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Muhammad",
			LastName:  "Ilham",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.Greet(ctx, greetRequest)

	if err != nil {

		statusErr, ok := status.FromError(err)

		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit!")
			} else {
				fmt.Println("Unexpected error: ", err)
			}
		} else {
			fmt.Printf("Error: %s", err)
		}
		return
	}

	fmt.Println(res.GetResult())
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Muhammad",
			LastName:  "Ilham",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error why streaming: %v", err)
		}
		log.Printf("Response froom GreetManyTimes: %v", msg.GetResult())
	}
}

func doSquareRoot(c greetpb.GreetServiceClient) {
	request := &greetpb.SquareRootRequest{
		Number: int32(-4),
	}
	res, err := c.SquareRoot(context.Background(), request)

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
		} else {
			log.Fatalf("Internal Server Error: %v", err)
		}
		return
	}

	fmt.Println("Result: ", res.GetRoot())

}
