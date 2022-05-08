package main

import (
	"log"
	"net/http"

	gws "github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"backend/pkg/websocket"
	pb "backend/proto"
)

var upgrader = gws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading to websocket protocol - %v", err)
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- client
	client.Read()
}

func main() {
	conn, err := grpc.Dial("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("error creating connection to grpc server: %v", err)
	}
	defer conn.Close()

	pool := websocket.NewPool(pb.NewMessageStoreClient(conn), 1)
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(pool, w, r)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("error serving websocket chat server - %v", err)
	}
}
