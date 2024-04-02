package api

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var (
	prox_port = flag.Int("prox_port", 8000, "Proxy server port")
	gRPC_port = flag.Int("grpc_port", 50051, "gRPC server port")
	gemini    *genai.Client
	chat      *Chat
)

func init() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil { // Check if file exists
		log.Fatal("Error loading .env file")
	}

	API_KEY := os.Getenv("API_KEY")
	if API_KEY == "" { // Check if API_KEY is present
		log.Fatal("API_KEY not found in environment")
	}

	chatStore, err := InitDataStores()
	if err != nil {
		log.Fatal(err)
	}
	chat = chatStore

	ctx := context.Background()
	gemini, err = genai.NewClient(ctx, option.WithAPIKey(API_KEY))
	if err != nil {
		log.Fatal(err)
	}
}
