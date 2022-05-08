package grpc

import (
	"context"
	"log"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/pkg/message"
	pb "backend/proto"
)

type Server struct {
	pb.UnimplementedMessageStoreServer
	messages []message.Message
}

func NewServer() Server {
	return Server{}
}

func (s *Server) StoreMessage(_ context.Context, msg *pb.Message) (*emptypb.Empty, error) {
	chatMsg := message.Message{
		ClientId:   msg.ClientId,
		ClientName: msg.ClientName,
		Time:       msg.Time.AsTime(),
		Content:    msg.Content,
	}
	log.Printf("storing message - %v", chatMsg)
	s.messages = append(s.messages, chatMsg)
	return &emptypb.Empty{}, nil
}

func (s *Server) RetrieveMessages(context.Context, *emptypb.Empty) (*pb.Messages, error) {
	log.Printf("retrieving stored messages - length: %d", len(s.messages))
	messages := make([]*pb.Message, 0)
	for _, msg := range s.messages {
		messages = append(messages, &pb.Message{
			ClientId:   msg.ClientId,
			ClientName: msg.ClientName,
			Time:       timestamppb.New(msg.Time),
			Content:    msg.Content,
		})
	}
	return &pb.Messages{StoredMessages: messages}, nil
}
