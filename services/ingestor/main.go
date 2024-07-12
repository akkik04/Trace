package main

import (
	"context"
	"fmt"
	"net"

	"github.com/akkik04/Trace/services/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.LogServiceServer
}

func (s server) SendLog(ctx context.Context, in *proto.LogEntry) (*proto.LogResponse, error) {
	fmt.Printf("Received: %v\n", in)
	return &proto.LogResponse{Message: "Successfully read log"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterLogServiceServer(grpcServer, server{})
	grpcServer.Serve(listener)
}
