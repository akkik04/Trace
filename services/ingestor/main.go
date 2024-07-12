package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/akkik04/Trace/services/proto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

// server struct to implement gRPC server.
type server struct {
	proto.LogServiceServer
}

// sendToS3 function to send log data to S3 bucket.
func sendToS3(logData *proto.LogEntry) error {

	// load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("AWS_S3_BUCKET_NAME")

	// initialize a session with explicit credentials.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	// create S3 service client.
	svc := s3.New(sess)

	// convert logData to string.
	logDataStr := logData.String()

	// prepare the input for the S3 'PutObject' call.
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fmt.Sprintf("logs/%d.log", time.Now().Unix())),
		Body:   strings.NewReader(logDataStr),
	}

	// perfom the S3 PutObject call.
	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	// log success message.
	fmt.Println("Successfully uploaded log to S3")
	return nil
}

// sendLog function to implement the gRPC server.
func (s server) SendLog(ctx context.Context, in *proto.LogEntry) (*proto.LogResponse, error) {

	fmt.Printf("Received from client: %v\n", in)

	// call function to send log data to S3
	err := sendToS3(in)
	if err != nil {
		fmt.Printf("Error sending log to S3: %v\n", err)
		return nil, err
	}
	return &proto.LogResponse{Message: "Successfully read log"}, nil
}

// main.
func main() {

	// create a listener on port 8082 for gRPC server.
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// create a new gRPC server.
	grpcServer := grpc.NewServer()
	proto.RegisterLogServiceServer(grpcServer, server{})
	fmt.Println("gRPC server started on port 8082")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
