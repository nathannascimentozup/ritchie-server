package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"ritchie-server/server"
)

type Handler struct {
	Config           server.Config
}

type response struct {
	Name        string `json:"name"`
	Priority    int    `json:"priority"`
	TreePath    string `json:"treePath"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
}

func NewConfigHandler(config server.Config) server.DefaultHandler {
	return Handler{Config: config}
}

func (lh Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var method = r.Method
		if http.MethodGet == method {
			lh.processGet(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (lh Handler) processGet(w http.ResponseWriter, r *http.Request) {

	organizationHeader := r.Header.Get(server.OrganizationHeader)
	repositoryConfigs, err := lh.Config.ReadRepositoryConfig(organizationHeader)

	if err != nil {
		log.Printf("Error while processing %v's repository configuration: %v", organizationHeader, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if repositoryConfigs == nil {
		log.Println("No repository config found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var resp []response
	for _, r:= range repositoryConfigs {
		res := response{
			Name:     r.Name,
			Priority: r.Priority,
			TreePath: r.ServerUrl + r.TreePath,
			Username: r.Username,
			Password: r.Password,
		}
		resp = append(resp, res)
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Sprintln("Error in Json Encode ")
		return
	}
}
