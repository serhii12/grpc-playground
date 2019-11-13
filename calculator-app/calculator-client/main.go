package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/serhii12/grpc-go/calculator-app/calculatorpb"
	pb "github.com/serhii12/grpc-go/calculator-app/calculatorpb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func doUnary(ctx context.Context, c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Sum Unary RPC...")

	req := &pb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 10,
	}

	r, err := c.Sum(ctx, req)
	if err != nil {
		log.Fatalf("could not sum: %v", err)
	}

	log.Printf("Response from Sum: %v", r.GetSumResult())
}

func doServerStreaming(ctx context.Context, c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 12390392840,
	}

	stream, err := c.PrimeNumberDecomposition(ctx, req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}

		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(ctx context.Context, c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage RPC: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, n := range numbers {
		fmt.Printf("Sending number: %v\n", n)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: n,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response: %v", err)
	}

	fmt.Printf("The Average is: %v\n", res.GetAverage())
}

func doBIStreaming(ctx context.Context, c calculatorpb.CalculatorServiceClient) {
	numbers := []int32{3, 5, 9, 54, 23}

	stream, err := c.FindMaximum(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage RPC: %v", err)
	}

	waitc := make(chan struct{})
	// we iterate over our slice and send each message individually
	go func(numbers []int32) {

		for _, n := range numbers {
			fmt.Printf("Sending number: %v\n", n)

			stream.Send(&pb.FindMaximumRequest{
				Number: n,
			})

			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}(numbers)

	// we receive a bunch of messages from the client (go routine)
	go func(stream pb.CalculatorService_FindMaximumClient) {
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
			fmt.Printf("Received: %v\n", res.GetMaximum())
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

	c := pb.NewCalculatorServiceClient(conn)

	// doUnary(context.Background(), c)
	// doServerStreaming(context.Background(), c)
	// doClientStreaming(context.Background(), c)
	doBIStreaming(context.Background(), c)
}
