package formulas

import (
	"log"
	"net/http"

	"ritchie-server/server"
	"ritchie-server/server/provider"
	"ritchie-server/server/tm"
)

type Handler struct {
	Config        server.Config
	Authorization server.Constraints
}


const (
	repoNameHeader = "x-repo-name"
	authorizationHeader = "Authorization"
)

func NewConfigHandler(config server.Config, auth server.Constraints) server.DefaultHandler {
	return Handler{Config: config, Authorization: auth}
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
	if repos == nil {
		log.Println("No repository configDummy found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	repoName := r.Header.Get(repoNameHeader)
	repo, err := tm.FindRepo(repos, repoName)
	if err != nil {
		log.Printf("no repo for org %s, with name %s, error: %v", org, repoName, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bt := r.Header.Get(authorizationHeader)
	ph := provider.NewProviderHandler(lh.Authorization, r.URL.Path, bt, org, repo)
	buf, err := ph.FilesFormulasAllow()
	if err != nil {
		log.Printf("error try allow access: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(buf)
	if err != nil {
		log.Printf("Failed to path: %s, error: %v", r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}