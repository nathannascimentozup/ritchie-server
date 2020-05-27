package security

import (
	"fmt"

	"ritchie-server/server"
	"ritchie-server/server/wpm"
)

type Authorization struct {
	config server.Config
}

func NewAuthorization(c server.Config) server.Constraints {
	return Authorization{
		config: c,
	}
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

	sc := auth.config.ReadSecurityConstraints()

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

func (auth Authorization) validateConstraints(path, method string, roles []string) bool {

	sc := auth.config.ReadSecurityConstraints()

	for _, pc := range sc.Constraints {
		if wpm.NewWildcardPattern(path, pc.Pattern).Match() {
			for _, role := range roles {
				rm := pc.RoleMappings[role]
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

//TODO: Alterar esse metodo para devolver o userInfo
func (auth Authorization) ListRealmRoles(token, org string) ([]string, error) {
	return []string{"role"}, nil
}
