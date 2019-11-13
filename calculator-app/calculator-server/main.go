package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/serhii12/grpc-go/calculator-app/calculatorpb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement CalculatorServiceServer
type server struct {
	pb.UnimplementedCalculatorServiceServer
}

// Sum implements calculator.CalculatorServiceServer
func (s *server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	log.Printf("Received numbers: %v, %v", in.GetFirstNumber(), in.GetSecondNumber())

	sum := in.GetFirstNumber() + in.GetSecondNumber()

	return &pb.SumResponse{
		SumResult: sum,
	}, nil
}

func (s *server) PrimeNumberDecomposition(in *pb.PrimeNumberDecompositionRequest, stream pb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", in)

	n := in.GetNumber()
	divisor := int64(2)
	for n > 1 {
		if n%divisor == 0 {
			stream.Send(&pb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})

			n = n / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increased to %v\n", divisor)
		}
	}
	return nil
}

func (s *server) ComputeAverage(stream pb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with %v\n", stream)

	var sum int32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Finished reading
			average := float64(sum) / float64(count)
			return stream.SendAndClose(&pb.ComputeAverageResponse{
				Average: average,
			})
		}

		if err != nil {
			log.Fatalf("Error reading file stream: %v", err)
		}

		sum += req.GetNumber()
		count++
	}
}

func (s *server) FindMaximum(stream pb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with %v\n", stream)

	var maximum int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Finished reading
			return nil
		}
		if err != nil {
			log.Fatalf("Error reading file stream: %v", err)
			return err
		}

		n := req.GetNumber()

		if n > maximum {
			maximum = n

			if err := stream.Send(&pb.FindMaximumResponse{
				Maximum: maximum,
			}); err != nil {
				log.Fatalf("Failed to send data to client: %v", err)
				return err
			}
		}
	}
}

func main() {
	fmt.Println("Calculator Server")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
