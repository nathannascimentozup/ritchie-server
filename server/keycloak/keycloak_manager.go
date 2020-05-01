package keycloak

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak"
	"net/http"
	"ritchie-server/server"
	"strconv"
	"strings"
)

const getUsersUrl = "%s/auth/admin/realms/%s/users?email=%s"

type Manager struct {
	config server.Config
}

func NewKeycloakManager(c server.Config) server.KeycloakManager {
	return Manager{config: c}
}

func (m Manager) CreateUser(user server.CreateUser, org string) (string, error)  {
	kc, err := m.config.ReadKeycloakConfigs(org)
	if err != nil {
		return "", err
	}
	c := gocloak.NewClient(kc.Url)
	jwt, err := m.loginClient(c, *kc)
	if err != nil {
		return "", err
	}
	u := buildKeyCloakUser(user)

	id, err := m.create(c, u, jwt.AccessToken, kc.Realm)
	if err != nil {
		return "", err
	}

	go m.actions(*id, jwt.AccessToken, kc.Url, kc.Realm)

	return *id, nil
}

func (m Manager) Login(org, user, password string) (string, int, error) {
	kc, err := m.config.ReadKeycloakConfigs(org)
	if err != nil {
		return "", 404, err
	}
	client := gocloak.NewClient(kc.Url)
	token, err := client.Login(kc.ClientId, kc.ClientSecret, kc.Realm, user, password)
	if err != nil {
		code := strings.Split(err.Error(), " ")[0]
		codeInt, errConverter := strconv.ParseInt(code, 10, 64)
		if errConverter != nil {
			return "", 500, err
		}
		return "", int(codeInt), err
	}
	return token.AccessToken, 0, nil
}

func (m Manager) DeleteUser(org, email string) error {
	kc, err := m.config.ReadKeycloakConfigs(org)
	if err != nil {
		return err
	}
	c := gocloak.NewClient(kc.Url)
	jwt, err := m.loginClient(c, *kc)
	if err != nil {
		return err
	}
	user, err := m.getUser(jwt.AccessToken, email, kc.Url, kc.Realm)
	if err != nil {
		return err
	}
	err = c.DeleteUser(jwt.AccessToken, kc.Realm, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m Manager) loginClient(c gocloak.GoCloak, kc server.KeycloakConfig) (*gocloak.JWT, error) {
	jwt, err := c.LoginClient(kc.ClientId, kc.ClientSecret, kc.Realm)
	if err != nil {
		return nil, err
	}

	return jwt, nil
}

func (m Manager) create(c gocloak.GoCloak, user gocloak.User, accessToken, realm string) (*string, error) {
	id, err := c.CreateUser(accessToken, realm, user)
	if err != nil {
		return nil, err
	}
	return id, nil
}


func (m Manager) actions(id, accessToken, keycloakUrl, realm string) {
	actions := []byte(`["UPDATE_PASSWORD", "VERIFY_EMAIL"]`)
	keyCloakEmailUrl := fmt.Sprintf("%v/auth/admin/realms/%v/users/%v/execute-actions-email", keycloakUrl, realm, id)

	req, _ := http.NewRequest(http.MethodPut, keyCloakEmailUrl, bytes.NewBuffer(actions))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	_, err := client.Do(req)
	if err != nil {
		fmt.Sprintln("Error in client.Do ")
		return
	}
}

func buildKeyCloakUser(u server.CreateUser) gocloak.User {
	return gocloak.User{
		ID:                         "",
		CreatedTimestamp:           0,
		Username:                   u.Username,
		Enabled:                    true,
		Totp:                       false,
		EmailVerified:              false,
		FirstName:                  u.FirstName,
		LastName:                   u.LastName,
		Email:                      u.Email,
		FederationLink:             "",
		Attributes:                 nil,
		DisableableCredentialTypes: nil,
		RequiredActions:            nil,
		Access:                     nil,
	}
}

func (m Manager) getUser(accessToken, email, keycloakUrl, realm string) (*gocloak.User, error) {
	url := fmt.Sprintf(getUsersUrl, keycloakUrl, realm, email)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, _ := http.DefaultClient.Do(req)

	var users []*gocloak.User
	defer res.Body.Close()

	_ = json.NewDecoder(res.Body).Decode(&users)
	if users == nil || len(users) <= 0 {
		return nil, errors.New("user not found")
	}

	return users[0], nil
}
