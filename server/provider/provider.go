package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"ritchie-server/server"
)

const (
	providerHttp = "HTTP"
	providerS3   = "S3"
)

type Handler struct {
	sec    server.Constraints
	bToken string
	org    string
	path   string
	repo   server.Repository
}


func NewProviderHandler(sec server.Constraints, bToken, org, path string, repo server.Repository) server.ProviderHandler {
	return Handler{
		sec:    sec,
		bToken: bToken,
		org:    org,
		path:   path,
		repo:   repo,
	}
}

func (hp Handler) TreeAllow() (server.Tree, error) {
	rTree, err := treeRemote(hp.path, hp.repo)
	if err != nil {
		return rTree, err
	}
	roles, err := hp.sec.ListRealmRoles(hp.bToken, hp.org)
	if err != nil {
		return rTree, err
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
	if hp.repo.ReplaceRepoUrl != "" {
		for _, c := range ft.Commands {
			if c.Formula != nil {
				if c.Formula.RepoUrl != "" {
					c.Formula.RepoUrl = hp.repo.ReplaceRepoUrl
				}
			}
		}
	}
	return ft, nil
}

func (hp Handler) FilesFormulasAllow() ([]byte, error) {
	tr, err := hp.TreeAllow()
	if err != nil {
		return nil, err
	}
	roles, err := hp.sec.ListRealmRoles(hp.bToken, hp.org)
	if err != nil {
		return nil, err
	}

	rfind := make(map[string]interface{})
	for _, r := range roles {
		rfind[strings.ToUpper(r.(string))] = r
	}
	p := strings.Replace(hp.path, "/formulas/", "", 1)
	s := strings.Split(p, "/")
	key := strings.ReplaceAll(p, "/"+s[len(s)-1], "")
	for _, c := range tr.Commands {
		if c.Formula != nil {
			if c.Formula.Path == key {
				return hp.bufProvider()
			}
		}
	}
	return nil, nil
}

func (hp Handler) bufProvider() ([]byte, error) {
	switch hp.repo.Provider.Type {
	case providerHttp:
		url := fmt.Sprintf("%s%s", hp.repo.Provider.Remote, hp.path)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	case providerS3:
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(hp.repo.Provider.Region)},
		)
		if err != nil {
			return nil, err
		}
		buf := &aws.WriteAtBuffer{}
		downloader := s3manager.NewDownloader(sess)
		s3obj := s3.GetObjectInput{
			Bucket: aws.String(hp.repo.Provider.Bucket),
			Key:    aws.String(hp.path),
		}
		_, err = downloader.Download(buf,
			&s3obj)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	default:
		return nil, errors.New(fmt.Sprintf("provider %s, not valid. Verify our repo config. Repo name: %s", hp.repo.Provider.Type, hp.repo.Name))
	}
}

func treeRemote(tPath string, repo server.Repository) (server.Tree, error) {
	switch repo.Provider.Type {
	case providerHttp:
		return loadTreeFileHttp(tPath, repo)
	case providerS3:
		return loadTreeFileS3(tPath, repo)
	default:
		return server.Tree{}, errors.New(fmt.Sprintf("provider %s, not valid. Verify our repo config. Repo name: %s", repo.Provider.Type, repo.Name))
	}
}

func loadTreeFileHttp(path string, repo server.Repository) (server.Tree, error) {
	url := fmt.Sprintf("%s%s", repo.Provider.Remote, path)
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

func loadTreeFileS3(path string, repo server.Repository) (server.Tree, error) {
	var response server.Tree
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(repo.Provider.Region)},
	)
	if err != nil {
		return response, err
	}
	buf := &aws.WriteAtBuffer{}
	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(buf,
		&s3.GetObjectInput{
			Bucket: aws.String(repo.Provider.Bucket),
			Key:    aws.String(path),
		})
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(buf.Bytes(), &response); err != nil {
		return response, err
	}
	return response, nil
}
