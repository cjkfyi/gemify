package api

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gemify/api/gen"
)

func StartHTTPProxy() error {
	flag.Parse() // prox_port
	r := chi.NewRouter()
	r.Post("/api/gemini", GeminiHandler)
	fmt.Printf("\n🌵  Proxy live at: 127.0.0.1:%v\n", *prox_port)
	return http.ListenAndServe(fmt.Sprintf(":%v", *prox_port), r)
}

func GeminiHandler(w http.ResponseWriter, r *http.Request) {
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
