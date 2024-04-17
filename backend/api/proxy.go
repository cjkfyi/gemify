package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gemify/api/gen"
	"gemify/store"
)

//
// Helpers

func redFlag(
	w http.ResponseWriter,
	statusCode int,
	res Response,
) {
	data, _ := json.Marshal(res)
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(statusCode)
	w.Write(data)
}

func greenLight(
	w http.ResponseWriter,
	statusCode int,
	res Response,
) {
	data, _ := json.Marshal(res)
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(statusCode)
	w.Write(data)
}

//

func SetupProxy() (*http.Server, string, error) {

	var addr = fmt.Sprintf(
		"%v:%v",
		*host_addr,
		*prox_port,
	)

	r := chi.NewRouter()

	//
	// Project API

	r.Post(
		"/p",
		CreateProj,
	)
	r.Get(
		"/projects",
		ListProjs,
	)
	r.Get(
		"/p/{projID}",
		GetProj,
	)
	r.Put(
		"/p/{projID}",
		UpdateProj,
	)
	r.Delete(
		"/p/{projID}",
		DeleteProj,
	)

	//
	// Chat API

	r.Post(
		"/p/{projID}",
		CreateChat,
	)
	r.Get(
		"/p/{projID}/chats",
		ListChats,
	)
	r.Get(
		"/p/{projID}/c/{chatID}",
		GetChat,
	)
	r.Put(
		"/p/{projID}/c/{chatID}",
		UpdateChat,
	)
	r.Delete(
		"/p/{projID}/c/{chatID}",
		DeleteChat,
	)

	//
	// Message API

	r.Post(
		"/p/{projID}/c/{chatID}",
		CreateMsg,
	)
	r.Get(
		"/p/{projID}/c/{chatID}/s",
		StreamMsg,
	)
	r.Get(
		"/p/{projID}/c/{chatID}/history",
		ListMsgs,
	)
	r.Get(
		"/p/{projID}/c/{chatID}/m/{msgID}",
		GetMsg,
	)
	r.Put(
		"/p/{projID}/c/{chatID}/m/{msgID}",
		UpdateMsg,
	)
	r.Delete(
		"/p/{projID}/c/{chatID}/m/{msgID}",
		DeleteMsg,
	)

	//

	proxySvr := http.Server{
		Addr:    addr,
		Handler: r,
	}
	return &proxySvr, addr, nil
}

//
// Experimental Handler

func StreamMsg(
	w http.ResponseWriter,
	r *http.Request,
) {

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		dataRes := map[string]interface{}{
			"code": store.ERR_Internal,
			"msg":  "failed to upgrade ws: " + err.Error(),
		}
		errorMessage := map[string]interface{}{
			"cmd":    "execNewMessage",
			"data":   dataRes,
			"status": "error",
		}
		jsonData, _ := json.Marshal(errorMessage)
		ws.WriteMessage(websocket.TextMessage, jsonData)
		ws.Close()
		return
	}
	defer ws.Close()

	_, msg, err := ws.ReadMessage()
	if err != nil {
		dataRes := map[string]interface{}{
			"code": store.ERR_Internal,
			"msg":  err.Error(),
		}
		errorMessage := map[string]interface{}{
			"cmd":    "execNewMessage",
			"data":   dataRes,
			"status": "error",
		}
		jsonData, _ := json.Marshal(errorMessage)
		ws.WriteMessage(websocket.TextMessage, jsonData)
		ws.Close()
		return
	}

	_, err = store.CreateMessage(projID, chatID, string(msg), true)
	if err != nil {
		dataRes := map[string]interface{}{
			"code": store.ERR_Internal,
			"msg":  err.Error(),
		}
		errorMessage := map[string]interface{}{
			"cmd":    "execNewMessage",
			"data":   dataRes,
			"status": "error",
		}
		jsonData, _ := json.Marshal(errorMessage)
		ws.WriteMessage(websocket.TextMessage, jsonData)
		ws.Close()
		return
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("127.0.0.1:%v", *gRPC_port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		dataRes := map[string]interface{}{
			"code": store.ERR_Internal,
			"msg":  err.Error(),
		}
		errorMessage := map[string]interface{}{
			"cmd":    "execNewMessage",
			"data":   dataRes,
			"status": "error",
		}
		jsonData, _ := json.Marshal(errorMessage)
		ws.WriteMessage(websocket.TextMessage, jsonData)
		ws.Close()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := pb.NewGemifyClient(conn)
	stream, err := client.SendMessage(ctx, &pb.Message{
		Content: string(msg),
		ChatID:  chatID,
		ProjID:  projID,
	})
	if err != nil {
		dataRes := map[string]interface{}{
			"code": store.ERR_Internal,
			"msg":  err.Error(),
		}
		errorMessage := map[string]interface{}{
			"cmd":    "execNewMessage",
			"data":   dataRes,
			"status": "error",
		}
		jsonData, _ := json.Marshal(errorMessage)
		ws.WriteMessage(websocket.TextMessage, jsonData)
		ws.Close()
		return
	}

	var buffer strings.Builder

	for {
		grpcResponse, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			dataRes := map[string]interface{}{
				"code": store.ERR_Internal,
				"msg":  "failed to receive gRPC res: " + err.Error(),
			}
			errorMessage := map[string]interface{}{
				"cmd":    "execNewMessage",
				"data":   dataRes,
				"status": "error",
			}
			jsonData, _ := json.Marshal(errorMessage)
			ws.WriteMessage(websocket.TextMessage, jsonData)
			ws.Close()
			break
		}

		jsonData, err := json.Marshal(grpcResponse)
		if err != nil {
			dataRes := map[string]interface{}{
				"code": store.ERR_Internal,
				"msg":  "failed to unmarshal gRPC res: " + err.Error(),
			}
			errorMessage := map[string]interface{}{
				"cmd":    "execNewMessage",
				"data":   dataRes,
				"status": "error",
			}
			jsonData, _ := json.Marshal(errorMessage)
			ws.WriteMessage(websocket.TextMessage, jsonData)
			ws.Close()
			break
		}

		err = ws.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			dataRes := map[string]interface{}{
				"code": store.ERR_Internal,
				"msg":  "failed to write ws msg: " + err.Error(),
			}
			errorMessage := map[string]interface{}{
				"cmd":    "execNewMessage",
				"data":   dataRes,
				"status": "error",
			}
			jsonData, _ := json.Marshal(errorMessage)
			ws.WriteMessage(websocket.TextMessage, jsonData)
			ws.Close()
			break
		}

		buffer.WriteString(grpcResponse.Content)
	}

	res := buffer.String()

	_, err = store.CreateMessage(projID, chatID, res, false)
	if err != nil {
		dataRes := map[string]interface{}{
			"code": store.ERR_Internal,
			"msg":  "failed to create new msg: " + err.Error(),
		}
		errorMessage := map[string]interface{}{
			"cmd":    "execNewMessage",
			"data":   dataRes,
			"status": "error",
		}
		jsonData, _ := json.Marshal(errorMessage)
		ws.WriteMessage(websocket.TextMessage, jsonData)
		ws.Close()
		return
	}
}
