package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gemify/store"
)

func CreateMsg(
	w http.ResponseWriter,
	r *http.Request,
) {
	var i store.Message

	chatID := chi.URLParam(r, "chatID")
	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode input",
		}
		res := Response{
			Command: "execCreateMsg",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	val, err := store.CreateMessage(chatID, projID, i.Message, i.IsUser)
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
				Command: "execCreateMsg",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"`message` field is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateMsg",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusUnprocessableEntity, res)
			return
		case
			"failed ds op",
			"failed to open chat ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execCreateMsg",
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
			Command: "execCreateMsg",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func GetMsg(
	w http.ResponseWriter,
	r *http.Request,
) {
	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")
	msgID := chi.URLParam(r, "msgID")

	message, err := store.GetMessage(projID, chatID, msgID)
	if err != nil {
		switch err.Error() {
		case
			"invalid `msgID` parameter",
			"invalid `chatID` parameter",
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execGetMsg",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open chat ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execGetMsg",
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
			Command: "execGetMsg",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func ListMsgs(
	w http.ResponseWriter,
	r *http.Request,
) {
	chatID := chi.URLParam(r, "chatID")
	projID := chi.URLParam(r, "projID")

	msgArr, err := store.ListMessages(chatID, projID)
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
				Command: "execListMsgs",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open chat ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execListMsgs",
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
			Command: "execListMsgs",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func UpdateMsg(
	w http.ResponseWriter,
	r *http.Request,
) {
	var i store.Message

	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")
	msgID := chi.URLParam(r, "msgID")

	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode input",
		}
		res := Response{
			Command: "execUpdateMsg",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	msg, err := store.UpdateMessage(projID, chatID, msgID, i)
	if err != nil {
		switch err.Error() {
		case
			"invalid `chatID` parameter",
			"invalid `projID` parameter",
			"invalid `msgID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateMsg",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open chat ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execUpdateMsg",
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
			Command: "execUpdateMsg",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func DeleteMsg(
	w http.ResponseWriter,
	r *http.Request,
) {
	projID := chi.URLParam(r, "projID")
	chatID := chi.URLParam(r, "chatID")
	msgID := chi.URLParam(r, "msgID")

	err := store.DeleteMessage(projID, chatID, msgID)
	if err != nil {
		switch err.Error() {
		case
			"invalid `projID` parameter",
			"invalid `chatID` parameter",
			"invalid `msgID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteMsg",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case
			"failed ds op",
			"failed to open chat ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execDeleteMsg",
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
			Command: "execDeleteMsg",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}
