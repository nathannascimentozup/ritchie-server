package keycloak

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/viniciussousazup/gocloak/v5"

	"ritchie-server/server"
)

const (
	url           = "url"
	realm         = "realm"
	clientId      = "clientId"
	clientSecret  = "clientSecret"
	ttl           = "ttl"
	otp           = "otp"
	validOtpError = "this realm have enable otp please send a totp value to login"
)

type keycloakConfig struct {
	client gocloak.GoCloak
	config kConfig
}

type kConfig struct {
	url          string
	realm        string
	clientId     string
	clientSecret string
	ttl          int64
	otp          bool
}

type keycloakError struct {
	code int
	err  error
}

type keycloakUser struct {
	roles    []string
	userInfo server.UserInfo
}

func NewKeycloakProvider(config map[string]string) server.SecurityManager {
	ttl, _ := strconv.ParseInt(config[ttl], 10, 64)
	otp, _ := strconv.ParseBool(config[otp])

	kc := kConfig{
		url:          config[url],
		realm:        config[realm],
		clientId:     config[clientId],
		clientSecret: config[clientSecret],
		ttl:          ttl,
		otp:          otp,
	}
	c := gocloak.NewClient(kc.url)
	return keycloakConfig{
		client: c,
		config: kc,
	}
}

func (k keycloakConfig) TTL() int64 {
	ttlF := time.Now().Unix() + k.config.ttl
	return ttlF
}

func (k keycloakConfig) Login(username, password, totp string) (server.User, server.LoginError) {

	if totp == "" && k.config.otp {
		return nil, keycloakError{
			code: 400,
			err:  errors.New(validOtpError),
		}
	}

	token, err := k.client.LoginOtp(k.config.clientId, k.config.clientSecret, k.config.realm, username, password, totp)
	if err != nil {
		code := strings.Split(err.Error(), " ")[0]
		codeInt, errConverter := strconv.ParseInt(code, 10, 64)
		if errConverter != nil {
			return nil, keycloakError{
				code: 500,
				err:  err,
			}
		}
		return nil, keycloakError{
			code: int(codeInt),
			err:  err,
		}
	}
	claims := jwt.MapClaims{}
	_, err = k.client.DecodeAccessTokenCustomClaims(token.AccessToken, k.config.realm, claims)
	if err != nil {
		return nil, keycloakError{
			code: 500,
			err:  err,
		}
	}
	realmAccess := claims["realm_access"].(map[string]interface{})
	ri := realmAccess["roles"].([]interface{})
	var roles []string
	for _, r := range ri {
		roles = append(roles, r.(string))
	}
	name := claims["name"].(string)
	email := claims["email"].(string)
	ku := keycloakUser{
		roles: roles,
		userInfo: server.UserInfo{
			Name:     name,
			Username: username,
			Email:    email,
		},
	}
	return ku, nil
}

func (ke keycloakError) Error() error {
	return ke.err
}
func (ke keycloakError) Code() int {
	return ke.code
}

func (ui keycloakUser) Roles() []string {
	return ui.roles
}
func (ui keycloakUser) UserInfo() server.UserInfo {
	return ui.userInfo
}
