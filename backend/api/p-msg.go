package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gemify/store"
)

func CreateMsg(w http.ResponseWriter, r *http.Request) {

	var i store.Message

	chatID := chi.URLParam(r, "chatID")
	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode msg input",
		}
		res := Response{
			Command: "execCreateMessage",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	val, err := store.CreateMessage(chatID, projID, i.Message, i.IsUser)
	if err != nil {
		switch err.Error() {
		case "failed to open chat ds",
			"failed to marshal msg",
			"failed to find proj with projID",
			"failed to store msg":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateMessage",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"message": val,
		}
		res := Response{
			Command: "execCreateMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func GetMsg(w http.ResponseWriter, r *http.Request) {

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
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execGetMessage",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"message": message,
		}
		res := Response{
			Command: "execGetMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func ListMsgs(w http.ResponseWriter, r *http.Request) {

	chatID := chi.URLParam(r, "chatID")
	projID := chi.URLParam(r, "projID")

	msgArr, err := store.ListMessages(chatID, projID)
	if err != nil {
		switch err.Error() {
		case "failed to open chat ds",
			"failed to pull msg with key",
			"failed to unmarshal msg":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execListMessages",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"messages": msgArr,
		}
		res := Response{
			Command: "execListMessages",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func UpdateMsg(w http.ResponseWriter, r *http.Request) {

	var i store.Message

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")
	msgID := chi.URLParam(r, "msgID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode msg input",
		}
		res := Response{
			Command: "execUpdateMessage",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	msg, err := store.UpdateMessage(projID, chatID, msgID, i)
	if err != nil {
		switch err.Error() {
		case "failed to pull msg with keys",
			"failed to unmarshal msg",
			"failed to marshal msg",
			"failed to store msg":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateMessage",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"message": msg,
		}
		res := Response{
			Command: "execUpdateMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func DeleteMsg(w http.ResponseWriter, r *http.Request) {

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
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteMessage",
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
			Command: "execDeleteMessage",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}
