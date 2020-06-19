package ldap

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jtblin/go-ldap-client"

	"ritchie-server/server"
)

const (
	base               = "base"
	host               = "host"
	serverName         = "serverName"
	port               = "port"
	useSSL             = "useSSL"
	skipTLS            = "skipTLS"
	insecureSkipVerify = "insecureSkipVerify"
	bindDN             = "bindDN"
	bindPassword       = "bindPassword"
	userFilter         = "userFilter"
	groupFilter        = "groupFilter"
	attributeUsername  = "attributeUsername"
	attributeName      = "attributeName"
	attributeEmail     = "attributeEmail"
	ttl                = "ttl"
	otp                = "otp"
)

type ldapError struct {
	code int
	err  error
}

type ldapUser struct {
	roles    []string
	userInfo server.UserInfo
}

type lConfig struct {
	base               string
	host               string
	serverName         string
	port               int
	useSSL             bool
	skipTLS            bool
	insecureSkipVerify bool
	bindDN             string
	bindPassword       string
	userFilter         string
	groupFilter        string
	attributeUsername  string
	attributeName      string
	attributeEmail     string
	ttl                int64
	otp                bool
}

type ldapConfig struct {
	client *ldap.LDAPClient
	config lConfig
}

func NewLdapProvider(config map[string]string) server.SecurityManager {
	cf := loadLConfig(config)
	cl := loadClient(cf)
	return ldapConfig{
		client: cl,
		config: cf,
	}
}

func loadClient(cf lConfig) *ldap.LDAPClient {
	att := []string{cf.attributeName, cf.attributeUsername, cf.attributeEmail}
	return &ldap.LDAPClient{
		Base:               cf.base,
		Host:               cf.host,
		ServerName:         cf.serverName,
		InsecureSkipVerify: cf.insecureSkipVerify,
		Port:               cf.port,
		UseSSL:             cf.useSSL,
		SkipTLS:            cf.skipTLS,
		BindDN:             cf.bindDN,
		BindPassword:       cf.bindPassword,
		UserFilter:         cf.userFilter,
		GroupFilter:        cf.groupFilter,
		Attributes:         att,
	}
}

func loadLConfig(config map[string]string) lConfig {
	p, _ := strconv.Atoi(config[port])
	us, _ := strconv.ParseBool(config[useSSL])
	st, _ := strconv.ParseBool(config[skipTLS])
	isv, _ := strconv.ParseBool(config[insecureSkipVerify])
	ttl, _ := strconv.ParseInt(config[ttl], 10, 64)
	otp, _ := strconv.ParseBool(config[otp])
	return lConfig{
		base:               config[base],
		host:               config[host],
		serverName:         config[serverName],
		port:               p,
		useSSL:             us,
		skipTLS:            st,
		insecureSkipVerify: isv,
		bindDN:             config[bindDN],
		bindPassword:       config[bindPassword],
		userFilter:         config[userFilter],
		groupFilter:        config[groupFilter],
		attributeUsername:  config[attributeUsername],
		attributeName:      config[attributeName],
		attributeEmail:     config[attributeEmail],
		ttl:                ttl,
		otp:                otp,
	}
}

func (l ldapConfig) Otp() bool {
	return l.config.otp
}

func (l ldapConfig) Login(username, password string) (server.User, server.LoginError) {
	defer l.client.Close()
	ok, user, err := l.client.Authenticate(username, password)
	if err != nil {
		return nil, ldapError{
			code: 401,
			err:  err,
		}
	}
	if !ok {
		return nil, ldapError{
			code: 401,
			err:  fmt.Errorf("Authenticating failed for user %s", username),
		}
	}
	groups, err := l.client.GetGroupsOfUser(username)
	if err != nil {
		return nil, ldapError{
			code: 500,
			err:  fmt.Errorf("Error getting groups for user %s", username),
		}
	}
	lu := ldapUser{
		roles: groups,
		userInfo: server.UserInfo{
			Name:     user[l.config.attributeName],
			Username: username,
			Email:    user[l.config.attributeEmail],
		},
	}
	return lu, nil
}

func (l ldapConfig) TTL() int64 {
	ttlF := time.Now().Unix() + l.config.ttl
	return ttlF
}

func (le ldapError) Error() error {
	return le.err
}
func (le ldapError) Code() int {
	return le.code
}

func (u ldapUser) Roles() []string {
	return u.roles
}
func (u ldapUser) UserInfo() server.UserInfo {
	return u.userInfo
}
