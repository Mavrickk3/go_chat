package websocket

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"backend/pkg/message"
	pb "backend/proto"
)

const defaultTimeout = 5

type Pool struct {
	Register      chan *Client
	Unregister    chan *Client
	Clients       map[*Client]bool
	ClientCounter int32
	Broadcast     chan message.Message
	MessageStore  pb.MessageStoreClient
	GRPCTimeout   time.Duration
}

func NewPool(msc pb.MessageStoreClient, grpcTimeout int) *Pool {
	if msc == nil {
		log.Println("message store client connection was not initialised - messages will not be saved and retrieved")
	}
	if grpcTimeout <= 0 {
		log.Println("invalid timeout value - using default timeout")
		grpcTimeout = defaultTimeout
	}

	return &Pool{
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		Clients:      make(map[*Client]bool),
		Broadcast:    make(chan message.Message),
		MessageStore: msc,
		GRPCTimeout:  time.Duration(grpcTimeout),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case newClient := <-pool.Register:
			newClient.ID = pool.ClientCounter
			newClient.Name = fmt.Sprintf("User-%d", newClient.ID)
			pool.ClientCounter += 1
			pool.Clients[newClient] = true
			pool.RetrieveMessages(newClient)
			msg := message.Message{
				ClientId:   newClient.ID,
				ClientName: newClient.Name,
				Time:       time.Now(),
				Content:    fmt.Sprintf("%s joined to the chat!", newClient.Name),
			}
			pool.StoreMessage(msg)
			pool.NotifyClients(msg)
			break
		case unregisteredClient := <-pool.Unregister:
			delete(pool.Clients, unregisteredClient)
			msg := message.Message{
				ClientId:   unregisteredClient.ID,
				ClientName: unregisteredClient.Name,
				Time:       time.Now(),
				Content:    fmt.Sprintf("%s left the chat!", unregisteredClient.Name),
			}
			pool.StoreMessage(msg)
			pool.NotifyClients(msg)
			break
		case msg := <-pool.Broadcast:
			pool.StoreMessage(msg)
			pool.NotifyClients(msg)
		}
	}
}

func (pool *Pool) RetrieveMessages(chatClient *Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*pool.GRPCTimeout)
	defer cancel()

	messages, err := pool.MessageStore.RetrieveMessages(ctx, &emptypb.Empty{})
	if err != nil {
		log.Printf("could not retrieve stored messages: %v", err)
		return
	}

	for _, msg := range messages.GetStoredMessages() {
		if err = chatClient.Conn.WriteJSON(message.Message{
			ClientId:   msg.GetClientId(),
			ClientName: msg.GetClientName(),
			Time:       msg.GetTime().AsTime(),
			Content:    msg.GetContent(),
		}); err != nil {
			log.Printf("could not send message: %v", err)
		}
	}
}

func (pool *Pool) StoreMessage(m message.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*pool.GRPCTimeout)
	defer cancel()

	_, err := pool.MessageStore.StoreMessage(ctx, &pb.Message{
		ClientId:   m.ClientId,
		ClientName: m.ClientName,
		Time:       timestamppb.New(m.Time),
		Content:    m.Content,
	})
	if err != nil {
		log.Printf("could not store message: %v", err)
	}
}

func (pool *Pool) NotifyClients(m message.Message) {
	for client, _ := range pool.Clients {
		if err := client.Conn.WriteJSON(m); err != nil {
			log.Printf("error notifying %s(%d) - %s", client.Name, client.ID, err)
		}
	}
}
