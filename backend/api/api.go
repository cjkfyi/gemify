package api

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

type Response struct {
	Command string                 `json:"command"`
	Data    map[string]interface{} `json:"data"`
	Status  string                 `json:"status"`
}

var (
	isTest    = flag.Bool("test_env", true, "Test environment")
	prox_port = flag.Int("prox_port", 8080, "Proxy server port")
	gRPC_port = flag.Int("grpc_port", 50051, "gRPC server port")
	host_addr = flag.String("host_addr", "127.0.0.1", "host address")

	gemini *genai.Client

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func getRole(isUser bool) string {
	if isUser {
		return "user"
	} else {
		return "model"
	}
}

func init() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("err loading .env file")
	}

	API_KEY := os.Getenv("API_KEY")
	if API_KEY == "" {
		log.Fatal("API_KEY was not found in .env")
	}

	ctx := context.Background()
	gemini, err = genai.NewClient(
		ctx,
		option.WithAPIKey(API_KEY),
	)
	if err != nil {
		log.Fatal(err)
	}
}
