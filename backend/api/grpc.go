package api

import (
	"context"
	"flag"
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

func StartGRPCServer() error {
	flag.Parse() // gRPC_port

	// Create the listener for the gRPC server address constructed
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *gRPC_port))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return err
	}

	s := grpc.NewServer()  // Basic initialization
	reflection.Register(s) // Only for test env
	pb.RegisterGeminiAPIServer(s, &server{})

	fmt.Printf("\n🌵  gRPC live at: 127.0.0.1:%v\n\n", *gRPC_port)

	// Serve gRPC or kick upwards an error
	if err := s.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
		return err
	} // goofy
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
