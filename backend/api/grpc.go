package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "gemify/api/gen"
)

func SetupGRPC() (*grpc.Server, net.Listener, string, error) {
	// gRPC address
	var port = *gRPC_port
	var host = "127.0.0.1"
	var addr = fmt.Sprintf(
		"%v:%v",
		host,
		port,
	)
	// Create a listener for that address
	lis, err := net.Listen("tcp", addr)
	if err != nil { // Wrap error
		return nil, nil, "", fmt.Errorf("failed to listen: %w", err)
	}

	grpcSvr := grpc.NewServer()
	reflection.Register(grpcSvr) // Only TEST env
	pb.RegisterGeminiAPIServer(grpcSvr, &server{})

	// Return on success
	return grpcSvr, lis, addr, nil
}

type server struct {
	pb.UnimplementedGeminiAPIServer
}

func (s *server) SendMessage(ctx context.Context, in *pb.Message) (*pb.Message, error) {

	// Model initialization
	gen := client.GenerativeModel("gemini-pro")
	gen.SetMaxOutputTokens(10000) // Tmp flow

	// Bursts of output
	var chunks []string

	cs := gen.StartChat()
	// Initialize the Chat Session
	iter := cs.SendMessageStream(
		ctx, // Send the msg
		genai.Text(in.Content),
	)
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		} // Gather up all of the chunks outputted
		chunks = append(chunks, printResponse(resp))
	}

	// Prettify the chunks into a str
	reply := fmt.Sprint(strings.Join(chunks, " "))

	// Construct the msg
	botReply := &pb.Message{
		Content: reply,
		IsUser:  false,
	}

	err := store.SaveConversation(in.Content, []byte(reply))
	if err != nil {
		log.Println("Error saving conversation:", err)
	}

	return botReply, nil
}

func printResponse(resp *genai.GenerateContentResponse) string {
	var output []string
	// List the different parts, construct
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				output = append(output, fmt.Sprint(part))
			}
		}
	} // Return a single string, tmp
	return strings.Join(output, " ")
}
