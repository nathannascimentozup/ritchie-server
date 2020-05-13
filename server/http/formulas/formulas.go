package formulas

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"ritchie-server/server"
	"ritchie-server/server/security"
	"ritchie-server/server/tm"
)

type Handler struct {
		Config server.Config
}


const (
	repoNameHeader = "x-repo-name"
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
	if repos == nil {
		log.Println("No repository config found")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	repoName := r.Header.Get(repoNameHeader)
	bt := r.Header.Get(authorizationHeader)
	sec := security.NewAuthorization(lh.Config)
	allow, repo, err := tm.FormulaAllow(sec, r.URL.Path, bt, repoName, org, repos)
	if err != nil {
		log.Printf("error try allow access: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !allow {
		log.Printf("Not allow access path: %s", r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	u, _ := url.Parse(repo.Remote)
	proxy := httputil.NewSingleHostReverseProxy(u)
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host
	proxy.ServeHTTP(w, r)
}