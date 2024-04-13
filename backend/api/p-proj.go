package api

import (
	"encoding/json"
	"errors"
	"gemify/store"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateProj(w http.ResponseWriter, r *http.Request) {

	var project store.Project

	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		data := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode proj input",
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
		case "name param is required",
			"desc param is required":
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
		case "name cannot exceed 160 chars",
			"desc cannot exceed 260 chars":
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
		case "failed to open meta ds",
			"failed to mk the proj dir",
			"failed to marshal proj",
			"failed to store proj in meta ds":
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

func GetProj(w http.ResponseWriter, r *http.Request) {

	// TODO:

	projID := chi.URLParam(r, "projID")
	if projID == "" {
		err := errors.New("projID param is required")
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
	}

	// finish...

	project, err := store.GetProject(projID)
	if err != nil {
		switch err.Error() {
		case "proj returned nil":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidProjID,
				"message": "param is invalid",
			}
			res := Response{
				Command: "execGetProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "projID param is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execGetProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "failed to open meta ds",
			"failed to find proj with key",
			"failed to unmarshal proj":
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

func ListProjs(w http.ResponseWriter, r *http.Request) {

	projects, err := store.ListProjects()
	if err != nil {
		switch err.Error() {
		case "failed to open meta ds",
			"failed to get pair in meta ds",
			"failed to unmarshal proj":
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

func UpdateProj(w http.ResponseWriter, r *http.Request) {

	var updatedData store.Project

	projID := chi.URLParam(r, "projID")

	err := json.NewDecoder(r.Body).Decode(&updatedData)
	if err != nil {
		errData := map[string]interface{}{
			"code":    store.ERR_Decode,
			"message": "failed to decode updated proj input",
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
		case "proj returned nil":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidProjID,
				"message": "param is invalid",
			}
			res := Response{
				Command: "execUpdateProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "projID param is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
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
		case "failed to open meta ds",
			"failed to marshal updated proj",
			"failed to scan meta ds for key",
			"failed to delete old proj entry",
			"failed to store new proj entry":
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
	}

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

//

func DeleteProj(w http.ResponseWriter, r *http.Request) {

	projID := chi.URLParam(r, "projID")

	err := store.DeleteProject(projID)
	if err != nil {
		switch err.Error() {
		case "projID param is required":
			data := map[string]interface{}{
				"code":    store.ERR_MissingInput,
				"message": err.Error(),
			}
			res := Response{
				Command: "execDeleteProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "failed to find proj with projID":
			data := map[string]interface{}{
				"code":    store.ERR_InvalidProjID,
				"message": "param is invalid",
			}
			res := Response{
				Command: "execDeleteProject",
				Data:    data,
				Status:  "error",
			}
			redFlag(w, http.StatusBadRequest, res)
			return
		case "failed to open meta ds",
			"failed to delete proj entry":
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
