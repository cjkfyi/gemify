package api

import (
	"fmt"
	"net"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "gemify/api/gen"
)

type server struct {
	pb.UnimplementedGemifyAPIServer
}

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
	pb.RegisterGemifyAPIServer(grpcSvr, &server{})

	// Return on success
	return grpcSvr, lis, addr, nil
}

func (s *server) SendMessage(
	input *pb.Message,
	stream pb.GemifyAPI_SendMessageServer,
) error {
	// Model initialization
	gen := gemini.GenerativeModel("gemini-pro")
	gen.SetMaxOutputTokens(10000) // Tmp flow

	cs := gen.StartChat()
	iter := cs.SendMessageStream(stream.Context(), genai.Text(input.Content))

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break // End of GenAI response stream
		}
		if err != nil {
			return err // Handle errors appropriately
		}

		botReply := &pb.Message{
			// Or modify how you extract content
			Content: printResponse(resp),
			IsUser:  false,
		}

		if err := stream.Send(botReply); err != nil {
			return err // Error sending to the client
		}
	}
	return nil
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
