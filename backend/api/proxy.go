package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gemify/api/gen"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func genConvoID() string {

	// TODO: implement google's UUID pkg

	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "error-generating-id"
	}
	return hex.EncodeToString(randomBytes)
}

func newConvo(w http.ResponseWriter, r *http.Request) {

	convoID := genConvoID()

	data := ConvoListData{
		ID:           convoID,
		Title:        "The meaning of life",
		LastModified: time.Now(),
	}

	err := chat.SaveNewConvo(convoID, data)
	if err != nil {
		http.Error(w, "Error saving new conversation", http.StatusInternalServerError)
		return
	}

	// Build the Res
	res := Response{
		Command: "execNewConvo",
		Status:  "success",
		Data: map[string]interface{}{
			"convoID": convoID,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func newMessage(w http.ResponseWriter, r *http.Request) {

	// Pull the conversationID from URL
	_ = chi.URLParam(r, "id")

	// 1. Upgrade to WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during WebSocket upgrade:", err) // Log the error
		return
	}
	defer ws.Close()

	_, val, err := ws.ReadMessage()
	if err != nil {
		fmt.Println("Error reading message:", err)
	}

	// Establish gRPC conn
	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%v", *gRPC_port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return
	}
	// Initiate gRPC Streaming
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := pb.NewGemifyAPIClient(conn)
	stream, err := client.SendMessage(ctx, &pb.Message{
		Content: string(val),
		IsUser:  true,
	})
	if err != nil {
		return
	}

	for {
		grpcResponse, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		jsonData, err := json.Marshal(grpcResponse)
		if err != nil {
			break
		}

		if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			break
		}
	}
}

func getConvos(w http.ResponseWriter, r *http.Request) {
	ListArr, err := chat.GetConvoList()
	if err != nil {
		http.Error(w,
			"Error retrieving conversation list",
			http.StatusInternalServerError,
		)
		return
	}

	// Build the Res struct
	res := Response{
		Command: "convoList",
		Status:  "success",
		Data: map[string]interface{}{
			"conversations": ListArr,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func SetupProxy() (*http.Server, string, error) {
	// Proxy address
	var port = *prox_port
	var host = "127.0.0.1"
	var addr = fmt.Sprintf(
		"%v:%v",
		host,
		port,
	)

	// Chi router instance
	r := chi.NewRouter()

	// Chi proxy routes
	r.Post("/chat", newConvo)
	r.Get("/chat/list", getConvos)
	r.Get("/ws/chat/{id}", newMessage)

	// Construct the server
	proxySvr := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Return on success
	return &proxySvr, addr, nil
}
