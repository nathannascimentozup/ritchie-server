package formulas

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"ritchie-server/server"
)

type (
	Handler struct {
		Config server.Config
	}
)

const (
	repoNameHeader = "x-repo-name"
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

	var repository server.Repository
	repoName := r.Header.Get(repoNameHeader)
	if repoName != "" {
		for _, r := range repositoryConfigs {
			if r.Name == repoName {
				repository = r
				break
			}
		}
	}
	//TODO: Validar se encontrou
	//TODO: Achar commando na tree file e valiar role

	u, _ := url.Parse(repository.ProxyTo)
	proxy := httputil.NewSingleHostReverseProxy(u)

	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host
	proxy.ServeHTTP(w, r)
}

