package tree

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ritchie-server/server"
	"ritchie-server/server/security"
	"ritchie-server/server/tm"
)

type Handler struct {
		Config server.Config
}

const (
	repoNameHeader      = "x-repo-name"
	authorizationHeader = "Authorization"
)

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
	org := r.Header.Get(server.OrganizationHeader)
	repos, err := lh.Config.ReadRepositoryConfig(org)
	if err != nil {
		log.Printf("Error while processing %v's repository configuration: %v", org, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if repos == nil || len(repos) == 0 {
		log.Println("No repository config found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	repoName := r.Header.Get(repoNameHeader)
	at := r.Header.Get(authorizationHeader)
	sec := security.NewAuthorization(lh.Config)
	finalTree, _, err := tm.TreeRemoteAllow(sec, at, org, repoName, r.URL.Path, repos)
	if err != nil {
		log.Printf("error load tree: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(finalTree)
	if err != nil {
		fmt.Sprintln("Error in Json Encode ")
		return
	}
}
