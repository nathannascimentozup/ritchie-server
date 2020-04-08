package server

import (
	"github.com/hashicorp/vault/api"
	"net/http"
)

const (
	OrganizationHeader = "x-org"
	ContextHeader      = "x-ctx"
	FileConfig         = "FILE_CONFIG"
)

type (
	Command struct {
		Usage   string   `json:"usage"`
		Help    string   `json:"help"`
		Formula *Formula `json:"formula,omitempty"`
		Parent  string   `json:"parent"`
	}

	Repository struct {
		Name     string `json:"name"`
		Priority int    `json:"priority"`
		TreePath string `json:"treePath"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	Formula struct {
		Path    string `json:"path"`
		Bin     string `json:"bin"`
		Config  string `json:"config"`
		RepoUrl string `json:"repoUrl"`
	}

	Credential struct {
		Service    string                 `json:"service"`
		Username   string                 `json:"username"`
		Credential map[string]interface{} `json:"credential"`
	}
	PermitMatcher struct {
		Pattern string   `yaml:"pattern"`
		Methods []string `yaml:"methods"`
	}

	DenyMatcher struct {
		Pattern      string              `yaml:"pattern"`
		RoleMappings map[string][]string `yaml:"roles"`
	}

	SecurityConstraints struct {
		Constraints       []DenyMatcher   `yaml:"constraints"`
		PublicConstraints []PermitMatcher `yaml:"publicConstraints"`
	}

	KeycloakConfig struct {
		Url          string `json:"url"`
		Realm        string `json:"realm"`
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
	}

	OauthConfig struct {
		Url      string `json:"url"`
		ClientId string `json:"clientId"`
	}
	CredentialConfig struct {
		Field string `json:"field"`
		Type  string `json:"type"`
	}

	ConfigFile struct {
		KeycloakConfig   *KeycloakConfig               `json:"keycloak"`
		OauthConfig      *OauthConfig                  `json:"oauth"`
		CredentialConfig map[string][]CredentialConfig `json:"credentials"`
		CliVersionConfig CliVersionConfig              `json:"cliVersionPath"`
		RepositoryConfig []Repository                  `json:"repositories"`
	}

	CliVersionConfig struct {
		Url      string `json:"url"`
		Provider string `json:"provider"`
		Version  string `json:"cliversion"`
	}

	Configuration struct {
		conf         Config
		vaultManager VaultManager
	}

	HealthEndpoints struct {
		KeycloakURL string
		VaultURL    string
	}

	WildcardPattern struct {
		str     string
		pattern string
	}

	SliceInterface struct {
		slice interface{}
	}

	CreateUser struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
	}
)

type Constraints interface {
	AuthorizationPath(bearerToken, path, method, org string) (bool, error)
	ValidatePublicConstraints(path, method string) bool
}

type ConfigHealth interface {
	ReadHealthConfigs() map[string]HealthEndpoints
}

type ConfigCredential interface {
	ReadCredentialConfigs(org string) (map[string][]CredentialConfig, error)
}

type ConfigKeycloak interface {
	ReadKeycloakConfigs(org string) (*KeycloakConfig, error)
}

type ConfigOauth interface {
	ReadOauthConfig(org string) (*OauthConfig, error)
}

type ConfigCliVersion interface {
	ReadCliVersionConfigs(org string) (CliVersionConfig, error)
}

type ConfigRepository interface {
	ReadRepositoryConfig(org string) ([]Repository, error)
}

type ConfigSecurityConstraints interface {
	ReadSecurityConstraints() SecurityConstraints
}

type Config interface {
	ConfigHealth
	ConfigCredential
	ConfigKeycloak
	ConfigOauth
	ConfigCliVersion
	ConfigRepository
	ConfigSecurityConstraints
}

type VaultManager interface {
	Write(key string, data map[string]interface{}) error
	Read(key string) (map[string]interface{}, error)
	List(key string) ([]interface{}, error)
	Delete(key string) error
	Start(*api.Client)
}

type KeycloakManager interface {
	CreateUser(user CreateUser, org string) (string, error)
	DeleteUser(org, email string) error
	Login(org, user, password string) (string, int, error)
}

type VaultConfig interface {
	Start() (*api.Client, error)
}

type DefaultHandler interface {
	Handler() http.HandlerFunc
}

type CredentialHandler interface {
	HandlerAdmin() http.HandlerFunc
	HandlerMe() http.HandlerFunc
}

type MiddlewareHandler interface {
	Filter(next http.Handler) http.Handler
}

type Configurator interface {
	LoadLoginHandler() DefaultHandler
	LoadConfigHandler() DefaultHandler
	LoadUserHandler() DefaultHandler
	LoadCredentialConfigHandler() DefaultHandler
	LoadConfigHealth() DefaultHandler
	LoadOauthHandler() DefaultHandler
	LoadUsageLoggerHandler() DefaultHandler
	LoadCliVersionHandler() DefaultHandler
	LoadRepositoryHandler() DefaultHandler
	LoadMiddlewareHandler() MiddlewareHandler
	LoadCredentialHandler() CredentialHandler
}

type WildcardPatternMatcher interface {
	Match() bool
}

type Slicer interface {
	Interface() ([]interface{}, error)
}
