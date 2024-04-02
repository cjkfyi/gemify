package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

//

// Temporary conversation state management
// type Conversation struct {
// 	gRPCStream pb.GemifyAPI_SendMessageServer
// }

// var conversations sync.Map

//

func genConvoID() string {

	// TODO: implement google's UUID pkg

	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "error-generating-id"
	}
	return hex.EncodeToString(randomBytes)
}

//

func newConvo(w http.ResponseWriter, r *http.Request) {

	// tmp solution
	convoID := genConvoID()

	data := ConvoListData{
		ID:           convoID,
		Title:        "",
		LastModified: time.Now(),
	}

	err := chat.SaveNewConvo(convoID, data)
	if err != nil {
		// Handle error ... you might return an error response here
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

	// Simple echo logic for testing
	// for {
	// 	key, val, err := ws.ReadMessage()
	// 	if err != nil {
	// 		fmt.Println("Error reading message:", err)
	// 		break
	// 	}

	// 	fmt.Printf("Received: %s\n", val)

	// 	err = ws.WriteMessage(key, val)
	// 	if err != nil {
	// 		fmt.Println("Error sending message:", err)
	// 		break
	// 	}
	// }

	_, val, err := ws.ReadMessage()
	if err != nil {
		fmt.Println("Error reading message:", err)
	}

	fmt.Printf("Received: %s\n", string(val))

	//

	// Establish gRPC conn
	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%v", *gRPC_port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		// connection error
		fmt.Println(err)
	}
	defer conn.Close()

	// Initiate gRPC Streaming
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // For cleanup later

	client := pb.NewGemifyAPIClient(conn)
	stream, err := client.SendMessage(ctx,
		&pb.Message{
			Content: string(val),
			IsUser:  true,
		},
	)
	if err != nil {
		return
		// ... handle gRPC setup error ...
	}
	// defer stream.CloseSend()
	// Close the stream from the client-side once done

	// Translation Loop
	for {
		grpcResponse, err := stream.Recv()
		if err == io.EOF { // Stream ended
			break
		}
		if err != nil {
			// Handle gRPC stream errors
			log.Println("Error receiving gRPC response:", err)
			// You may want to send an error message over the WebSocket here
			break
		}

		// Translate grpcResponse into WebSocket-friendly format
		jsonData, err := json.Marshal(grpcResponse)
		if err != nil {
			return
			// ... handle JSON encoding errors ...
		}

		// Send over WebSocket
		if err := ws.WriteMessage(websocket.TextMessage, jsonData); err != nil {
			// ... handle WebSocket send errors ...
			return

		}
	}
}

func getConvos(w http.ResponseWriter, r *http.Request) {
	ListArr, err := chat.GetConvoList()
	if err != nil {
		// return an error response here
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
	r.Get("/chat/list", getConvos)
	r.Post("/chat", newConvo)
	r.Get("/ws/chat/{id}", newMessage)
	// first new msg, updates the title summarized.

	// r.Get("/chat/list/s", getShortConvoList)

	// r.Post("/chat/{convoId}/messages", sendMessageHandler)
	// r.Get("/ws/chat/{convoId}", wsConversationHandler)

	// Construct the server
	proxySvr := http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Return on success
	return &proxySvr, addr, nil
}
