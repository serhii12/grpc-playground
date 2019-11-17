package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/codes"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	pb "github.com/serhii12/grpc-go/blog-app/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

const (
	port = ":50051"
)

// server is used to implement BlogServiceServer
type server struct {
	pb.UnimplementedBlogServiceServer
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Content  string             `bson:"content"`
	Title    string             `bson:"title"`
}

func (s *server) CreateBlog(ctx context.Context, req *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	fmt.Println("Create blog request")
	blog := req.GetBlog()

	newBlog, err := collection.InsertOne(context.Background(), blogItem{
		AuthorID: blog.GetAuthorId(),
		Content:  blog.GetContent(),
		Title:    blog.GetTitle(),
	})
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Failed to instert blog: %v", err),
		)
	}

	bid, ok := newBlog.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to blog id: %v", err),
		)
	}

	resp := &pb.CreateBlogResponse{
		Blog: &pb.Blog{
			Id:       bid.Hex(),
			AuthorId: blog.GetAuthorId(),
			Content:  blog.GetContent(),
			Title:    blog.GetTitle(),
		},
	}

	return resp, nil
}

func (s *server) ReadBlog(ctx context.Context, req *pb.ReadBlogRequest) (*pb.ReadBlogResponse, error) {
	fmt.Println("Create blog request")

	bid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse blog id: %v", err),
		)
	}

	blog := &blogItem{}
	filter := bson.D{primitive.E{Key: "_id", Value: bid}}
	result := collection.FindOne(context.Background(), filter)
	if err := result.Decode(blog); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find a blog: %v", err),
		)
	}

	resp := &pb.ReadBlogResponse{
		Blog: &pb.Blog{
			Id:       req.GetBlogId(),
			AuthorId: blog.AuthorID,
			Content:  blog.Content,
			Title:    blog.Title,
		},
	}

	return resp, nil
}

func (s *server) UpdateBlog(ctx context.Context, req *pb.UpdateBlogRequest) (*pb.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")

	bid, err := primitive.ObjectIDFromHex(req.GetBlog().GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse blog id: %v", err),
		)
	}

	filter := bson.D{primitive.E{Key: "_id", Value: bid}}
	updateFields := bson.M{
		"$set": bson.M{
			"author_id": req.GetBlog().GetAuthorId(),
			"content":   req.GetBlog().GetContent(),
			"title":     req.GetBlog().GetTitle(),
		},
	}

	if _, err := collection.UpdateOne(context.Background(), filter, updateFields); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Failed to update a blog: %v", err),
		)
	}

	blog := &blogItem{}
	findFilter := bson.D{primitive.E{Key: "_id", Value: bid}}
	b := collection.FindOne(context.Background(), findFilter)
	if err := b.Decode(blog); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find a blog: %v", err),
		)
	}

	resp := &pb.UpdateBlogResponse{
		Blog: &pb.Blog{
			Id:       req.GetBlog().GetId(),
			AuthorId: blog.AuthorID,
			Content:  blog.Content,
			Title:    blog.Title,
		},
	}

	return resp, nil
}

func (s *server) DeleteBlog(ctx context.Context, req *pb.DeleteBlogRequest) (*pb.DeleteBlogResponse, error) {
	fmt.Println("DeleteBlog blog request")

	bid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse blog id: %v", err),
		)
	}

	filter := bson.D{primitive.E{Key: "_id", Value: bid}}
	if _, err := collection.DeleteOne(context.Background(), filter); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Failed to delete a blog: %v", err),
		)
	}

	resp := &pb.DeleteBlogResponse{
		BlogId: req.GetBlogId(),
	}

	return resp, nil
}

func (s *server) ListBlog(req *pb.ListBlogRequest, stream pb.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")

	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Could not find blogs: %v", err),
		)
	}

	defer cur.Close(context.Background())

	for cur.Next(context.Background()) {

		data := &blogItem{}
		if err := cur.Decode(data); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)
		}

		if err := stream.Send(&pb.ListBlogResponse{Blog: &pb.Blog{
			Id:       data.ID.Hex(),
			AuthorId: data.AuthorID,
			Content:  data.Content,
			Title:    data.Title,
		}}); err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Failed to send data: %v", err),
			)
		}
	}

	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}

	return nil
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	fmt.Println("Blog Service Started")
	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterBlogServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server")

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for Control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDB Connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
	fmt.Println("End of Program")
}
