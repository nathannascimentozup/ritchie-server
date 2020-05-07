package security

import (
	"fmt"
	"github.com/Nerzal/gocloak/v4"
	"github.com/dgrijalva/jwt-go"
	"ritchie-server/server"
	"ritchie-server/server/wpm"
	"strings"
)

const (
	bearer = "Bearer "
)

type Authorization struct {
	Config server.Config
}

func NewAuthorization(c server.Config) server.Constraints {
	return Authorization{Config: c}
}

func (auth Authorization) AuthorizationPath(bearerToken, path, method, org string) (bool, error) {
	if org == "" {
		return false, fmt.Errorf("x-org not received. ")
	}
	roles, err := auth.ListRealmRoles(bearerToken, org)
	if err != nil {
		return false, err
	}
	return auth.validateConstraints(path, method, roles), nil
}

func (auth Authorization) ValidatePublicConstraints(path, method string) bool {

	sc := auth.Config.ReadSecurityConstraints()

	for _, pc := range sc.PublicConstraints {
		if wpm.NewWildcardPattern(path, pc.Pattern).Match() {
			for _, m := range pc.Methods {
				if method == m {
					return true
				}
			}
		}
	}
	return false
}

func (auth Authorization) validateConstraints(path, method string, roles []interface{}) bool {

	sc := auth.Config.ReadSecurityConstraints()

	for _, pc := range sc.Constraints {
		if wpm.NewWildcardPattern(path, pc.Pattern).Match() {
			for _, role := range roles {
				rm := pc.RoleMappings[role.(string)]
					for _, m := range rm {
						if method == m {
							return true
						}
					}
			}
		}
	}
	return false
}

func (auth Authorization) ListRealmRoles(bearerToken, org string) ([]interface{}, error) {
	if "" == bearerToken {
		return nil, fmt.Errorf("Bearer Token is empty ")
	}
	if !strings.Contains(bearerToken, bearer) {
		return nil, fmt.Errorf("Bearer Token is not valid ")
	}
	jwtString := strings.Replace(bearerToken, bearer, "", -1)
	if "" == jwtString {
		return nil, fmt.Errorf("Bearer Token result is empty ")
	}
	kConfig, err := auth.Config.ReadKeycloakConfigs(org)
	if err != nil {
		return nil, err
	}
	client := gocloak.NewClient(kConfig.Url)
	claims := jwt.MapClaims{}
	_, err = client.DecodeAccessTokenCustomClaims(jwtString, kConfig.Realm, claims)
	if err != nil {
		return nil, err
	}
	realmAccess := claims["realm_access"].(map[string]interface{})
	roles := realmAccess["roles"].([]interface{})
	return roles, nil
}
