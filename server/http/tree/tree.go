package tree

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"ritchie-server/server"
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
	}

	formula struct {
		Path       string   `json:"path"`
		Bin        string   `json:"bin,omitempty"`
		BinWindows string   `json:"binWindows,omitempty"`
		BinDarwin  string   `json:"binDarwin,omitempty"`
		BinLinux   string   `json:"binLinux,omitempty"`
		Config     string   `json:"config,omitempty"`
		RepoUrl    string   `json:"repoUrl"`
		Roles      []string `json:"roles,omitempty"`
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
	internalTreeUrl := repository.InternalUrl + r.URL.Path
	t, err := loadTreeFile(internalTreeUrl)

	//TODO: Remover os que nao possui roles

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		fmt.Sprintln("Error in Json Encode ")
		return
	}
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
