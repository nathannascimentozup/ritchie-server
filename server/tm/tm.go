package tm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"ritchie-server/server"
)

func TreeRemoteAllow(sec server.Constraints, bToken, org, repoName, tPath string, repos []server.Repository) (server.Tree, server.Repository, error) {
	rTree, repo, err := treeRemote(org, repoName, tPath, repos)
	if err != nil {
		return rTree, repo, err
	}

	roles, err := sec.ListRealmRoles(bToken, org)
	if err != nil {
		return rTree, repo, err
	}
	rfind := make(map[string]interface{})
	for _, r := range roles {
		rfind[strings.ToUpper(r.(string))] = r
	}
	ft := server.Tree{}
	ft.Version = rTree.Version
	for _, c := range rTree.Commands {
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
	if repo.ReplaceRepoUrl != "" {
		for _, c := range ft.Commands {
			if c.Formula != nil {
				if c.Formula.RepoUrl != "" {
					c.Formula.RepoUrl = repo.ReplaceRepoUrl
				}
			}
		}
	}
	return ft, repo, nil
}

func FormulaAllow(sec server.Constraints, path, token, repoName, org string, repos []server.Repository) (bool, server.Repository, error) {
	tr, repo, err := TreeRemoteAllow(sec, token, org, repoName, path, repos)
	if err != nil {
		return false, repo, err
	}

	roles, err := sec.ListRealmRoles(token, org)
	if err != nil {
		return false, repo, err
	}

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
							return true, repo, nil
						}
					}
					return false, repo, nil
				} else {
					return true, repo, nil
				}
			}
		}
	}
	return false, repo, nil
}

func treeRemote(org, repoName, tPath string, repos []server.Repository) (server.Tree, server.Repository, error) {
	var tree server.Tree
	var repository server.Repository
	for _, r := range repos {
		if r.Name == repoName {
			repository = r
			break
		}
	}
	if repository.Name == "" {
		return tree, repository, fmt.Errorf("No repo for org %s with name %s\n", org, repoName)
	}

	tURL := fmt.Sprintf("%s%s", repository.Remote, tPath)
	t, err := loadTreeFile(tURL)
	if err != nil {
		return tree, repository, err
	}
	return t, repository, nil
}

func loadTreeFile(url string) (server.Tree, error) {
	var response server.Tree
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
		return response, err
	}
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return response, err
	}
	return response, nil
}
