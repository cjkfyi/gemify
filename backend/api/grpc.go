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
	"gemify/store"
)

type server struct {
	pb.UnimplementedGemifyServer
}

func SetupGRPC() (
	*grpc.Server,
	net.Listener,
	string, error,
) {

	var addr = fmt.Sprintf(
		"%v:%v",
		*host_addr,
		*gRPC_port,
	)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, "",
			fmt.Errorf("failed to listen: %w", err)
	}

	grpcSvr := grpc.NewServer()
	if *isTest {
		reflection.Register(grpcSvr)
	}
	pb.RegisterGemifyServer(grpcSvr, &server{})

	return grpcSvr, lis, addr, nil
}

func (s *server) SendMessage(
	i *pb.Message,
	stream pb.Gemify_SendMessageServer,
) error {

	var convo []*genai.Content

	hist, err := store.ListMessages(
		i.ProjID, i.ChatID,
	)
	if err != nil {
		return err
	}

	model := gemini.GenerativeModel(
		"gemini-1.5-pro-latest",
	)
	cs := model.StartChat()

	for _, el := range hist {
		ex := &genai.Content{
			Parts: []genai.Part{
				genai.Text(el.Message),
			},
			Role: getRole(el.IsUser),
		}
		convo = append(convo, ex)
	}
	cs.History = convo

	iter := cs.SendMessageStream(
		stream.Context(),
		genai.Text(i.Content),
	)

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		botReply := &pb.Message{
			Content: printRes(resp),
		}

		if err := stream.Send(botReply); err != nil {
			return err
		}
	}
	return nil
}

func printRes(
	res *genai.GenerateContentResponse,
) string {
	var output []string
	for _, cand := range res.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				output = append(output, fmt.Sprint(part))
			}
		}
	}
	return strings.Join(output, " ")
}
