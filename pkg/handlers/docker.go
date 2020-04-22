package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/softplan/tenkai-docker-api/pkg/constraints"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
	"github.com/softplan/tenkai-docker-api/pkg/global"
	"github.com/softplan/tenkai-docker-api/pkg/util"
)

func (appContext *AppContext) listDockerTags(w http.ResponseWriter, r *http.Request) {

	w.Header().Set(global.ContentType, global.JSONContentType)

	var payload model.ListDockerTagsRequest

	if err := util.UnmarshalPayload(r, &payload); err != nil {
		http.Error(w, err.Error(), 501)
		return
	}

	result, err := appContext.DockerServiceAPI.GetDockerTagsWithDate(payload, appContext.Repositories.DockerDAO, &appContext.DockerTagsCache)
	if err != nil {
		http.Error(w, err.Error(), 501)
		return
	}

	data, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (appContext *AppContext) listDockerRepositories(w http.ResponseWriter, r *http.Request) {

	principal := util.GetPrincipal(r)
	if !util.Contains(principal.Roles, constraints.TenkaiAdmin) {
		http.Error(w, errors.New(global.AccessDenied).Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set(global.ContentType, global.JSONContentType)
	result := &model.ListDockerRepositoryResponse{}
	var err error
	if result.Repositories, err = appContext.Repositories.DockerDAO.ListDockerRepos(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (appContext *AppContext) listDockerVariables(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(global.ContentType, global.JSONContentType)

	var payload model.DockerVariablesPayload
	var err error

	if err = util.UnmarshalPayload(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	splHost := strings.Split(payload.ImageName, "/")
	host := splHost[0]

	var dockerRepo model.DockerRepo
	if dockerRepo, err = appContext.Repositories.DockerDAO.
		GetDockerRepositoryByHost(host); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := &model.DockerVariablesResponse{}
	if result, err = appContext.DockerServiceAPI.GetDockerVariables(&dockerRepo,
		payload.ImageName, payload.ImageTag); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(result)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func (appContext *AppContext) newDockerRepository(w http.ResponseWriter, r *http.Request) {

	principal := util.GetPrincipal(r)
	if !util.Contains(principal.Roles, constraints.TenkaiAdmin) {
		http.Error(w, errors.New(global.AccessDenied).Error(), http.StatusUnauthorized)
	}

	w.Header().Set(global.ContentType, global.JSONContentType)

	var payload model.DockerRepo

	if err := util.UnmarshalPayload(r, &payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := appContext.Repositories.DockerDAO.CreateDockerRepo(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (appContext *AppContext) deleteDockerRepository(w http.ResponseWriter, r *http.Request) {

	principal := util.GetPrincipal(r)
	if !util.Contains(principal.Roles, constraints.TenkaiAdmin) {
		http.Error(w, errors.New(global.AccessDenied).Error(), http.StatusUnauthorized)
	}

	vars := mux.Vars(r)
	sl := vars["id"]
	id, _ := strconv.Atoi(sl)
	w.Header().Set(global.ContentType, global.JSONContentType)
	if err := appContext.Repositories.DockerDAO.DeleteDockerRepo(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

}
