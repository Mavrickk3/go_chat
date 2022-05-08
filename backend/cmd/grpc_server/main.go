package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	chatGRPC "backend/pkg/grpc"
	pb "backend/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	chatGrpcServer := chatGRPC.NewServer()
	pb.RegisterMessageStoreServer(s, &chatGrpcServer)
	log.Printf("server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
