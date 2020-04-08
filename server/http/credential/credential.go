package credential

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"

	"ritchie-server/server"
)

const (
	credentialVaultPath = "/%s/%s/%s"
	authorizationHeader = "Authorization"
	bearer              = "Bearer "
)

type user struct {
	username   string
	email      string
	name       string
	familyName string
}

type Handler struct {
	v server.VaultManager
	c server.Config
}

func NewCredentialHandler(v server.VaultManager, c server.Config) server.CredentialHandler {
	return Handler{v: v, c: c}
}

func (h Handler) defaultValidate(c server.Credential, org string) map[string]string {
	errs := make(map[string]string)

	credentials, _ := h.c.ReadCredentialConfigs(org)

	found := false

	for k := range credentials {
		if k == c.Service {
			found = true
			break
		}
	}

	if !found {
		errs["credential"] = "The credential is not valid."
	}
	if c.Service == "" {
		errs["service"] = "The service field is required."
	}
	if len(c.Credential) == 0 {
		errs["credential"] = "The credential field must not be empty."
	}
	return errs
}

func (h Handler) findCredential(org, ctx string, c server.Credential) (map[string]interface{}, error) {
	path := fmt.Sprintf(credentialVaultPath, ctxResolver(org, ctx), c.Username, c.Service)
	credential, err := h.v.Read(path)
	if err != nil {
		return nil, err
	}

	return credential, nil
}

func (h Handler) createCredential(org, ctx string, c server.Credential) error {
	path := fmt.Sprintf(credentialVaultPath, ctxResolver(org, ctx), c.Username, c.Service)
	if err := h.v.Write(path, c.Credential); err != nil {
		return err
	}
	return nil
}

func ctxResolver(org, ctx string) string {
	if ctx != "" {
		return fmt.Sprintf("%s_%s", org, ctx)
	}

	return org
}

func org(r *http.Request) string {
	return r.Header.Get(server.OrganizationHeader)
}

func ctx(r *http.Request) string {
	return r.Header.Get(server.ContextHeader)
}

func loadUser(r http.Request) user {
	authorizationToken := r.Header.Get(authorizationHeader)
	jwtString := strings.Replace(authorizationToken, bearer, "", -1)
	token, _ := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	claims := token.Claims.(jwt.MapClaims)
	name := claims["given_name"].(string)
	familyName := claims["family_name"].(string)
	username := claims["preferred_username"].(string)
	email := claims["email"].(string)
	return user{
		username:   username,
		email:      email,
		name:       name,
		familyName: familyName,
	}
}
