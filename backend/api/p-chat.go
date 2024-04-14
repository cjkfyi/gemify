package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gemify/store"
)

func CreateChat(
	w http.ResponseWriter,
	r *http.Request,
) {
	var i *store.Chat

	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode input",
		}
		res := Response{
			Command: "execCreateChat",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	i.ProjID = projID

	chat, err := store.CreateChat(i)
	if err != nil {
		switch err.Error() {
		case
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"`name` param is required",
			"`desc` param is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"`name` cannot exceed 160 chars",
			"`desc` cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execCreateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"chat": chat,
		}
		res := Response{
			Command: "execCreateChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func GetChat(
	w http.ResponseWriter,
	r *http.Request,
) {
	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	chat, err := store.GetChat(projID, chatID)
	if err != nil {
		switch err.Error() {
		case
			"invalid `chatID` parameter",
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execGetChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execGetChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"chat": chat,
		}
		res := Response{
			Command: "execGetChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func ListChats(
	w http.ResponseWriter,
	r *http.Request,
) {
	projID := chi.URLParam(r, "projID")

	chats, err := store.ListChats(projID)
	if err != nil {
		switch err.Error() {
		case
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execListChats",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execListChats",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"chats": chats,
		}
		res := Response{
			Command: "execListChats",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func UpdateChat(
	w http.ResponseWriter,
	r *http.Request,
) {
	var i store.Chat

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode input",
		}
		res := Response{
			Command: "execUpdateChat",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	updated, err := store.UpdateChat(projID, chatID, i)
	if err != nil {
		switch err.Error() {
		case
			"`projID` param is required",
			"`chatID` param is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"invalid `projID` parameter",
			"invalid `chatID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"name cannot exceed 160 chars",
			"desc cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execUpdateChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"chat": updated,
		}
		res := Response{
			Command: "execUpdateChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func DeleteChat(
	w http.ResponseWriter,
	r *http.Request,
) {
	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")

	err := store.DeleteChat(projID, chatID)
	if err != nil {
		switch err.Error() {
		case
			"`projID` param is required",
			"`chatID` param is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"invalid `chatID` parameter",
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execDeleteChat",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"deleted": true,
		}
		res := Response{
			Command: "execDeleteChat",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}
