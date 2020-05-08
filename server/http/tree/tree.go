package tree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"ritchie-server/server"
	"ritchie-server/server/security"
	"strings"
	"time"
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
	//TODO: Validar se encontrou - Talvez um metodo
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

	internalTreeUrl := repository.ProxyTo + r.URL.Path
	t, err := loadTreeFile(internalTreeUrl)
	//TODO: Validar se encontrou

	authorizationToken := r.Header.Get(authorizationHeader)
	sec := security.NewAuthorization(lh.Config)
	roles, err := sec.ListRealmRoles(authorizationToken, org)

	finalTree := finalTreeFile(roles, t)
	if repository.ReplaceRepoUrl != "" {
		for _, c := range finalTree.Commands {
			if c.Formula != nil {
				if c.Formula.RepoUrl != "" {
					c.Formula.RepoUrl = repository.ReplaceRepoUrl
				}
			}
		}
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(finalTree)
	if err != nil {
		fmt.Sprintln("Error in Json Encode ")
		return
	}
}

func finalTreeFile(roles []interface{}, actualTree tree) tree {
	rfind := make(map[string]interface{})
	for _, r := range roles {
		rfind[strings.ToUpper(r.(string))] = r
	}
	ft := tree{}
	ft.Version = actualTree.Version
	for _, c := range actualTree.Commands {
		if len(c.Roles) > 0 {
			for _, r := range c.Roles {
				if rfind[strings.ToUpper(r)] != nil {
					ft.Commands = append(ft.Commands, c)
				}
			}
		} else {
			ft.Commands = append(ft.Commands, c)
		}
	}
	return ft
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
