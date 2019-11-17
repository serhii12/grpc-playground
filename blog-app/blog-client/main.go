package main

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/serhii12/grpc-go/blog-app/blogpb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	fmt.Println("Blog Client")

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	c := pb.NewBlogServiceClient(conn)

	// --- Create Blog START ---
	fmt.Println("Creating Blog")

	blog := &pb.Blog{
		AuthorId: "1",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	resp, err := c.CreateBlog(context.Background(), &pb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Could create a blog: %v\n", err)
	}

	fmt.Printf("Blog has been created %v\n", resp)
	// --- Create Blog FINISHED ---

	// --- Read Blog START ---
	fmt.Println("Reading Blog")

	readResp, err := c.ReadBlog(context.Background(), &pb.ReadBlogRequest{BlogId: resp.GetBlog().GetId()})
	if err != nil {
		log.Fatalf("Could read a blog: %v\n", err)
	}

	fmt.Printf("Blog has been read %v\n", readResp)
	// --- Read Blog FINISHED ---

	// --- Update Blog START ---
	fmt.Println("Updating Blog")

	newBlog := &pb.Blog{
		Id:       resp.GetBlog().GetId(),
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}

	updateRes, err := c.UpdateBlog(context.Background(), &pb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("Error happened while updating: %v \n", err)
	}

	fmt.Printf("Blog was updated: %v\n", updateRes)
	// --- Update Blog FINISHED ---

	// --- Delete Blog START ---
	deleteRes, err := c.DeleteBlog(context.Background(), &pb.DeleteBlogRequest{BlogId: resp.GetBlog().GetId()})
	if err != nil {
		fmt.Printf("Error happened while deleting: %v \n", err)
	}

	fmt.Printf("Blog was deleted: %v \n", deleteRes)
	// --- Delete Blog FINISHED ---

	// --- List Blog START ---
	stream, err := c.ListBlog(context.Background(), &pb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}

		fmt.Println(res.GetBlog())
	}
	// --- List Blog FINISHED ---
}
