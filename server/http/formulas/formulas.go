package formulas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"ritchie-server/server"
	"ritchie-server/server/security"
)

type (
	Handler struct {
		Config server.Config
	}

	tree struct {
		Commands []command `json:"commands"`
		Version  string    `json:"version"`
	}

	command struct {
		Usage   string   `json:"usage"`
		Help    string   `json:"help"`
		Formula *formula `json:"formula,omitempty"`
		Parent  string   `json:"parent"`
		Roles   []string `json:"roles,omitempty"`
	}

	formula struct {
		Path       string   `json:"path"`
		Bin        string   `json:"bin,omitempty"`
		BinWindows string   `json:"binWindows,omitempty"`
		BinDarwin  string   `json:"binDarwin,omitempty"`
		BinLinux   string   `json:"binLinux,omitempty"`
		Config     string   `json:"config,omitempty"`
		RepoUrl    string   `json:"repoUrl"`
	}
)

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
	repositoryConfigs, err := lh.Config.ReadRepositoryConfig(org)
	if err != nil {
		log.Printf("Error while processing %v's repository configuration: %v", org, err)
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
	t, err := loadTreeFile(repository.ProxyTo + repository.TreePath)
	//TODO: Tratar erro e ver se veio tree para repo
	bt := r.Header.Get(authorizationHeader)
	if !lh.allow(t, r.URL.Path, bt, org) {
		log.Printf("Not allow access path: %s", r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	u, _ := url.Parse(repository.ProxyTo)
	proxy := httputil.NewSingleHostReverseProxy(u)
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host
	proxy.ServeHTTP(w, r)
}

func (lh Handler) allow(tr tree, path, token, org string) bool {
	sec := security.NewAuthorization(lh.Config)
	roles, _ := sec.ListRealmRoles(token, org)
	rfind := make(map[string]interface{})
	for _, r := range roles {
		rfind[strings.ToUpper(r.(string))] = r
	}
	p := strings.Replace(path, "/formulas/", "", 1)
	s := strings.Split(p, "/")
	key := strings.ReplaceAll(p, "/" + s[len(s) -1], "")
	for _, c := range tr.Commands {
		if c.Formula != nil {
			if c.Formula.Path == key {
				if len(c.Roles) > 0 {
					for _, r := range c.Roles {
						if rfind[strings.ToUpper(r)] != nil {
							return true
						}
					}
					return false
				} else {
					return true
				}
			}
		}
	}
	return false
}

func loadTreeFile(url string) (tree, error) {
	var response tree
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return response, err
	}

	hc := &http.Client{Timeout: 5 * time.Second}
	resp, err := hc.Do(req)
	if err != nil {
		return response, err
	}

	if resp.StatusCode != 200 {
		return response, fmt.Errorf("%d - failed to get index for %s\n", resp.StatusCode, url)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(bodyBytes, &response)
	return response, nil
}

