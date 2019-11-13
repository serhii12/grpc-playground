package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/serhii12/grpc-go/greet-app/greetpb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func doUnaryAPI(cxt context.Context, c pb.GreetServiceClient) {
	fmt.Println("Testing grpc Unary")

	req := &pb.GreetRequest{
		Greeting: &pb.Greeting{
			FirstName: "Alex",
			LastName:  "Brown",
		},
	}

	resp, err := c.Greet(cxt, req)
	if err != nil {
		log.Fatalf("could not sum: %v", err)
	}

	log.Printf("Response from Greet: %v", resp.Result)
}

func doServerStream(cxt context.Context, c pb.GreetServiceClient) {
	fmt.Println("Testing grpc server stream")

	req := &pb.GreetManyTimesRequest{
		Greeting: &pb.Greeting{
			FirstName: "Alex",
			LastName:  "Brown",
		},
	}

	resp, err := c.GreetManyTimes(cxt, req)
	if err != nil {
		log.Fatalf("could not sum: %v", err)
	}

	for {
		msg, err := resp.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("error reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes client: %v", msg.GetResult())
	}

}

func doClientStream(cxt context.Context, c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*pb.LongGreetRequest{
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Stephane",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "John",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Lucy",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Mark",
			},
		},
		&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: "Piper",
			},
		},
	}

	stream, err := c.LongGreet(cxt)
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiStream(cxt context.Context, c pb.GreetServiceClient) {
	fmt.Println("Starting to do a BI Streaming RPC...")

	stream, err := c.GreetEveryone(cxt)
	if err != nil {
		log.Fatalf("error while calling GreetEveryone: %v", err)
	}

	requests := []*pb.GreetEveryoneRequest{
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Stephane",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "John",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Lucy",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Mark",
			},
		},
		&pb.GreetEveryoneRequest{
			Greeting: &pb.Greeting{
				FirstName: "Piper",
			},
		},
	}

	waitc := make(chan struct{})
	// we iterate over our slice and send each message individually
	go func(requests []*pb.GreetEveryoneRequest) {
		for _, req := range requests {
			fmt.Printf("Sending req: %v\n", req)

			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}(requests)

	// we receive a bunch of messages from the client (go routine)
	go func(stream pb.GreetService_GreetEveryoneClient) {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}(stream)

	// block until everything is done
	<-waitc
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	ctx := context.Background()
	// doUnaryAPI(ctx, c)

	// doServerStream(ctx, c)

	// doClientStream(ctx, c)

	doBiStream(ctx, c)
}
