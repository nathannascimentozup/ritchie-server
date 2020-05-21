package provider

import (
	"github.com/aws/aws-sdk-go/aws"

	"ritchie-server/server"
)

type Handler struct {
	sec server.Constraints
	bToken string
	org string
	path string
	repo server.Repository
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
	return server.Tree{}, nil
}

func (hp Handler) FilesFormulasAllow() (*aws.WriteAtBuffer, error) {
	return nil, nil
}
