package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	pb "github.com/serhii12/grpc-go/greet-app/greetpb"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement GreetServiceServer
type server struct {
	pb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Println("Calling greet func")

	fN := in.GetGreeting().GetFirstName()

	result := "Hello " + fN

	resp := &pb.GreetResponse{
		Result: result,
	}

	return resp, nil
}

func (s *server) GreetManyTimes(in *pb.GreetManyTimesRequest, stream pb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", in)

	firstName := in.GetGreeting().GetFirstName()

	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)

		resp := &pb.GreetManytimesResponse{
			Result: result,
		}

		stream.Send(resp)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (s *server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with %v\n", stream)

	var result string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// Finished reading
			return stream.SendAndClose(&pb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error reading file stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (s *server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with %v\n", stream)

	var result string
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

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "

		if err := stream.Send(&pb.GreetEveryoneResponse{
			Result: result,
		}); err != nil {
			log.Fatalf("Failed to send data to client: %v", err)
			return err
		}
	}

}

func main() {
	fmt.Println("Greet Server")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
