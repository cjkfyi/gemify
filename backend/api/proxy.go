package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "gemify/api/gen"
	"gemify/models"
	"gemify/store"
)

func SetupProxy() (*http.Server, string, error) {

	// Our wise tree
	var addr = fmt.Sprintf(
		"%v:%v",
		*host_addr,
		*prox_port,
	)
	//
	r := chi.NewRouter()
	////
	////
	r.Get("/projs", ListProjectsHandler)
	r.Post("/proj", CreateProjectHandler)
	r.Get("/proj/{projID}", GetProjectHandler)
	r.Put("/proj/{projID}", UpdateProjectHandler)
	r.Delete("/proj/{projID}", DeleteProjectHandler)
	////
	r.Post("/proj/{projID}/chat", CreateChatHandler)
	r.Get("/proj/{projID}/chats", ListChatsHandler)
	r.Get("/proj/{projID}/chat/{chatID}", GetChatHandler)
	r.Put("/proj/{projID}/chat/{chatID}", UpdateChatHandler)
	r.Delete("/proj/{projID}/chat/{chatID}", DeleteChatHandler)
	////
	r.Get("/proj/{projID}/chat/{chatID}/msg", NewMessageHandler)
	r.Post("/proj/{projID}/chat/{chatID}/msg", CreateMessageHandler)
	r.Get("/proj/{projID}/chat/{chatID}/history", ListMessagesHandler)
	r.Get("/proj/{projID}/chat/{chatID}/msg/{msgID}", GetMessageHandler)
	r.Put("/proj/{projID}/chat/{chatID}/msg/{msgID}", UpdateMessageHandler)
	r.Delete("/proj/{projID}/chat/{chatID}/msg/{msgID}", DeleteMessageHandler)
	////
	////
	proxySvr := http.Server{
		Addr:    addr,
		Handler: r,
	}
	return &proxySvr, addr, nil
}

//
// Helpers

func redFlag(
	w http.ResponseWriter,
	statusCode int,
	res models.Response,
) {
	errJSON, _ := json.Marshal(res)
	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	w.WriteHeader(statusCode)
	w.Write(errJSON)
}

func greenLight(
	w http.ResponseWriter,
	statusCode int,
	res models.Response,
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
// Project Handlers

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {

	var project models.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		data := map[string]interface{}{
			"code":    models.ErrorCodeDecode,
			"message": "failed to decode proj input",
		}
		res := models.Response{
			Command: "execCreateProject",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
	}

	proj, err := store.CreateProject(&project)
	if err != nil {
		switch err.Error() {
		case "name param is required",
			"desc param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "name cannot exceed 160 chars",
			"desc cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInvalidInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to open meta ds",
			"failed to mk the proj dir",
			"failed to marshal proj",
			"failed to store proj in meta ds":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"project": proj,
		}
		res := models.Response{
			Command: "execCreateProject",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusCreated, res)
	}
}

func GetProjectHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")

	project, err := store.GetProject(projID)
	if err != nil {
		switch err.Error() {
		case "projID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execGetProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to open meta ds",
			"failed to find proj with key",
			"failed to unmarshal proj":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execGetProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"project": project,
		}
		errorResponse := models.Response{
			Command: "execGetProject",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, errorResponse)
	}
}

func ListProjectsHandler(w http.ResponseWriter, r *http.Request) {

	projects, err := store.ListProjects()
	if err != nil {
		switch err.Error() {
		case "failed to open meta ds",
			"failed to get pair in meta ds",
			"failed to unmarshal proj":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execListProjects",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"projects": projects,
		}
		errorResponse := models.Response{
			Command: "execListProjects",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, errorResponse)
	}
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {

	var updatedData models.Project

	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		errData := map[string]interface{}{
			"code":    models.ErrorCodeDecode,
			"message": "failed to decode updated proj input",
		}
		res := models.Response{
			Command: "execUpdateProject",
			Data:    errData,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
	}

	proj, err := store.UpdateProject(projID, updatedData)
	if err != nil {
		switch err.Error() {
		case "projID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execUpdateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "name cannot exceed 160 chars",
			"desc cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInvalidInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execUpdateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to open meta ds",
			"failed to marshal updated proj",
			"failed to scan meta ds for key",
			"failed to delete old proj entry",
			"failed to store new proj entry":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execUpdateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	}

	data := map[string]interface{}{
		"project": proj,
	}
	res := models.Response{
		Command: "execUpdateProject",
		Data:    data,
		Status:  "success",
	}
	greenLight(w, http.StatusOK, res)
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")

	err := store.DeleteProject(projID)
	if err != nil {
		switch err.Error() {
		case "projID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execDeleteProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to find proj with projID":
			data := map[string]interface{}{
				"code":    models.ErrorCodeWrongKey,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execDeleteProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusNotFound, res)
		case "failed to open meta ds",
			"failed to delete proj entry":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execDeleteProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
		return
	} else {
		data := map[string]interface{}{
			"deleted": true,
		}
		errorResponse := models.Response{
			Command: "execDeleteProject",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, errorResponse)
	}
}

//
// Chat Handlers

func CreateChatHandler(w http.ResponseWriter, r *http.Request) {

	var i *models.Chat

	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    models.ErrorCodeDecode,
			"message": "failed to decode chat input",
		}
		res := models.Response{
			Command: "execCreateChat",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
	}

	i.ProjID = projID

	chat, err := store.CreateChat(i)
	if err != nil {
		switch err.Error() {
		case "projID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "name param is required",
			"desc param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "name cannot exceed 160 chars",
			"desc cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInvalidInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to store proj in meta ds",
			"failed to store new chat entity",
			"failed to delete old chat entity",
			"failed to open chat ds",
			"failed to open meta ds",
			"failed to mk the proj dir",
			"failed to marshal proj",
			"failed to unmarshal proj",
			"failed to find proj with key",
			"failed to find proj with projID":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"chat": chat,
		}
		res := models.Response{
			Command: "execCreateChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

func GetChatHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	chat, err := store.GetChat(projID, chatID)
	if err != nil {
		switch err.Error() {
		case "projID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execGetChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to open meta ds",
			"failed to find proj with projID",
			"failed to find chat with chatID",
			"failed to unmarshal proj":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execGetChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"chat": chat,
		}
		errorResponse := models.Response{
			Command: "execGetChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, errorResponse)
	}
}

func ListChatsHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")

	chats, err := store.ListChats(projID)
	if err != nil {
		switch err.Error() {
		case "projID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execListChats",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to open meta ds",
			"failed to find proj with key",
			"failed to unmarshal proj":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execListChats",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"chats": chats,
		}
		errorResponse := models.Response{
			Command: "execListChats",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, errorResponse)
	}
}

func UpdateChatHandler(w http.ResponseWriter, r *http.Request) {

	var i models.Chat

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    models.ErrorCodeDecode,
			"message": "failed to decode chat input",
		}
		res := models.Response{
			Command: "execUpdateChat",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
	}

	updated, err := store.UpdateChat(projID, chatID, i)
	if err != nil {
		switch err.Error() {
		case "projID param is required",
			"chatID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execUpdateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "name cannot exceed 160 chars",
			"desc cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInvalidInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execUpdateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to open meta ds",
			"chat not found with chatID",
			"failed to marshal proj",
			"failed to find proj with projID",
			"failed to delete old chat entity",
			"failed to store new chat entity",
			"failed to find proj with key",
			"failed to unmarshal proj":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execUpdateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"chat": updated,
		}
		res := models.Response{
			Command: "execUpdateChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

func DeleteChatHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	err := store.DeleteChat(projID, chatID)
	if err != nil {
		switch err.Error() {
		case "projID param is required",
			"chatID param is required":
			data := map[string]interface{}{
				"code":    models.ErrorCodeMissingInput,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execDeleteChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
		case "failed to open meta ds",
			"failed to find proj with projID",
			"failed to unmarshal proj",
			"failed to find chat with chatID",
			"failed to marshal proj",
			"failed to store new proj":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execDeleteChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"deleted": true,
		}
		errorResponse := models.Response{
			Command: "execDeleteChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, errorResponse)
	}
}

//
// Message Handlers

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {

	var i *models.Message

	chatID := chi.URLParam(r, "chatID")
	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(i)
	if err != nil {
		data := map[string]interface{}{
			"code":    models.ErrorCodeDecode,
			"message": "failed to decode msg input",
		}
		res := models.Response{
			Command: "execCreateMessage",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
	}

	val, err := store.CreateMessage(chatID, projID, i)
	if err != nil {
		switch err.Error() {
		case "failed to open chat ds",
			"failed to marshal msg",
			"failed to store msg":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execCreateMessage",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"message": val,
		}
		res := models.Response{
			Command: "execCreateMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")
	msgID := chi.URLParam(r, "msgID")

	message, err := store.GetMessage(projID, chatID, msgID)
	if err != nil {
		switch err.Error() {
		case "failed to open chat ds",
			"failed to pull msg with key",
			"failed to unmarshal msg":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execGetMessage",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"message": message,
		}
		res := models.Response{
			Command: "execGetMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

func ListMessagesHandler(w http.ResponseWriter, r *http.Request) {

	chatID := chi.URLParam(r, "chatID")
	projID := chi.URLParam(r, "projID")

	msgArr, err := store.ListMessages(chatID, projID)
	if err != nil {
		switch err.Error() {
		case "failed to open chat ds",
			"failed to pull msg with key",
			"failed to unmarshal msg":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execListMessages",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"messages": msgArr,
		}
		res := models.Response{
			Command: "execListMessages",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {

	var i models.Message

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")
	msgID := chi.URLParam(r, "msgID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    models.ErrorCodeDecode,
			"message": "failed to decode msg input",
		}
		res := models.Response{
			Command: "execUpdateMessage",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
	}

	msg, err := store.UpdateMessage(projID, chatID, msgID, i)
	if err != nil {
		switch err.Error() {
		case "failed to pull msg with keys",
			"failed to unmarshal msg",
			"failed to marshal msg",
			"failed to store msg":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execUpdateMessage",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"message": msg,
		}
		res := models.Response{
			Command: "execUpdateMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")
	msgID := chi.URLParam(r, "msgID")

	err := store.DeleteMessage(projID, chatID, msgID)
	if err != nil {
		switch err.Error() {
		case "failed to pull msg with key",
			"failed to store updated msg",
			"failed to marshal msg",
			"failed to unmarshal msg":
			data := map[string]interface{}{
				"code":    models.ErrorCodeInternal,
				"message": err.Error(),
			}
			res := models.Response{
				Command: "execDeleteMessage",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
		}
	} else {
		data := map[string]interface{}{
			"deleted": true,
		}
		res := models.Response{
			Command: "execDeleteMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//
// Experimental Handlers

func NewMessageHandler(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		dataRes := map[string]interface{}{
			"code": models.ErrorCodeInternal,
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
			"code": models.ErrorCodeInternal,
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

	userStamp := int(time.Now().UnixNano())

	usrMsg := &models.Message{
		ID:           store.GenID(),
		ChatID:       chatID,
		ProjID:       projID,
		IsUser:       true,
		Message:      string(msg),
		LastModified: userStamp,
		FirstCreated: userStamp,
	}
	_, err = store.CreateMessage(chatID, projID, usrMsg)
	if err != nil {
		dataRes := map[string]interface{}{
			"code": models.ErrorCodeInternal,
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
			"code": models.ErrorCodeInternal,
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
			"code": models.ErrorCodeInternal,
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
				"code": models.ErrorCodeInternal,
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
				"code": models.ErrorCodeInternal,
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
				"code": models.ErrorCodeInternal,
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

	response := buffer.String()

	gemStamp := int(time.Now().UnixNano())
	gemMsg := &models.Message{
		ID:           store.GenID(),
		ChatID:       chatID,
		ProjID:       projID,
		IsUser:       false,
		Message:      response,
		LastModified: gemStamp,
		FirstCreated: gemStamp,
	}
	_, err = store.CreateMessage(chatID, projID, gemMsg)
	if err != nil {
		dataRes := map[string]interface{}{
			"code": models.ErrorCodeInternal,
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
