package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"gemify/store"
)

func CreateProj(
	w http.ResponseWriter,
	r *http.Request,
) {
	var project store.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode input",
		}
		res := Response{
			Command: "execCreateProj",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	proj, err := store.CreateProject(&project)
	if err != nil {
		switch err.Error() {
		case
			"`name` field is required",
			"`desc` field is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateProj",
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
				Command: "execCreateProj",
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
				Command: "execCreateProj",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"project": proj,
		}
		res := Response{
			Command: "execCreateProj",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusCreated, res)
	}
}

//

func GetProj(
	w http.ResponseWriter,
	r *http.Request,
) {
	projID := chi.URLParam(r, "projID")

	project, err := store.GetProject(projID)
	if err != nil {
		switch err.Error() {
		case
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execGetProj",
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
				Command: "execGetProj",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"project": project,
		}
		res := Response{
			Command: "execGetProj",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func ListProjs(
	w http.ResponseWriter,
	r *http.Request,
) {
	projects, err := store.ListProjects()
	if err != nil {
		switch err.Error() {
		case
			"failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execListProjs",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"projects": projects,
		}
		res := Response{
			Command: "execListProjs",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func UpdateProj(
	w http.ResponseWriter,
	r *http.Request,
) {
	var updatedData store.Project

	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		errData := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode input",
		}
		res := Response{
			Command: "execUpdateProj",
			Data:    errData,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	proj, err := store.UpdateProject(projID, updatedData)
	if err != nil {
		switch err.Error() {
		case
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateProj",
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
				Command: "execUpdateProj",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusUnprocessableEntity, res)
			return

		case
			"failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": "oops, something uhh...",
			}
			res := Response{
				Command: "execUpdateProj",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusInternalServerError, res)
			return
		}
	} else {
		data := map[string]interface{}{
			"project": proj,
		}
		res := Response{
			Command: "execUpdateProj",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}

//

func DeleteProj(
	w http.ResponseWriter,
	r *http.Request,
) {
	projID := chi.URLParam(r, "projID")

	err := store.DeleteProject(projID)
	if err != nil {
		switch err.Error() {
		case
			"invalid `projID` parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidParam,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteProj",
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
				Command: "execDeleteProj",
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
			Command: "execDeleteProj",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}
