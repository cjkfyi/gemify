package api

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

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
	r.Get("/ws", websocketHandler)
	r.Post("/api/gemini", geminiHandler)

	// Construct the server
	proxySvr := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Return on success
	return &proxySvr, addr, nil
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during WebSocket upgrade:", err) // Log the error
		return
	}
	defer ws.Close()

	// Simple echo logic for testing
	for {
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received: %s\n", message)

		err = ws.WriteMessage(messageType, message)
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}

func geminiHandler(w http.ResponseWriter, r *http.Request) {
	flag.Parse() // gRPC_port

	// Establish gRPC conn
	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%v", *gRPC_port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		// connection error
		fmt.Println(err)
	}
	defer conn.Close()

	client := pb.NewGeminiAPIClient(conn)

	// Extract data from HTTP req
	var requestData map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		// JSON decoding error
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding request body")
		return
	}

	message, ok := requestData["message"].(string)
	if !ok {
		// Missing or wrong message type ...
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Missing or wrong msg type")
		return
	}

	// Construct the gRPC req
	grpcRequest := &pb.Message{
		Content: message,
		IsUser:  true,
	}

	// Send the gRPC req
	grpcResponse, err := client.SendMessage(
		context.Background(),
		grpcRequest,
	)
	if err != nil {
		// gRPC error
		fmt.Println(err)
	}

	// Translate gRPC res to HTTP
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": grpcResponse.GetContent(),
	})
	if err != nil {
		// JSON encoding error
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error encoding response")
	}
}
