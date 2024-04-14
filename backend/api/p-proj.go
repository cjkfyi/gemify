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
			Command: "execCreateProject",
			Data:    data,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	proj, err := store.CreateProject(&project)
	if err != nil {
		switch err.Error() {
		case "`name` param is required",
			"`desc` param is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "`name` cannot exceed 160 chars",
			"`desc` cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execCreateProject",
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
			Command: "execCreateProject",
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
		case "invalid projID parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidProjID,
				"message": err.Error(),
			}
			res := Response{
				Command: "execGetProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execGetProject",
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
			Command: "execGetProject",
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
		case "failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execListProjects",
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
			Command: "execListProjects",
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
			Command: "execUpdateProject",
			Data:    errData,
			Status:  "error",
		}
		redFlag(w, http.StatusBadRequest, res)
		return
	}

	proj, err := store.UpdateProject(projID, updatedData)
	if err != nil {
		switch err.Error() {
		case "invalid projID parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidProjID,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return

		case "name cannot exceed 160 chars",
			"desc cannot exceed 260 chars":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return

		case "failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execUpdateProject",
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
			Command: "execUpdateProject",
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
		case "invalid projID parameter":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidProjID,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "failed ds op",
			"failed to open meta ds":
			data := map[string]interface{}{
				"code":    store.ERR_Internal,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteProject",
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
			Command: "execDeleteProject",
			Data:    data,
			Status:  "success",
		}
		greenLight(w, http.StatusOK, res)
	}
}
