package main

import (
	"log"

	"github.com/jtblin/go-ldap-client"
)

func main() {
	client := &ldap.LDAPClient{
		Base:         "dc=example,dc=org",
		Host:         "localhost",
		ServerName:   "ldap.example.org",
		InsecureSkipVerify: false,
		Port:         389,
		UseSSL:       false,
		SkipTLS:      true,
		BindDN:       "cn=admin,dc=example,dc=org",
		BindPassword: "admin",
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
	// It is the responsibility of the caller to close the connection
	defer client.Close()

	ok, user, err := client.Authenticate("user", "user")
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", "user", err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", "user")
	}
	log.Printf("User: %+v", user)

	groups, err := client.GetGroupsOfUser("user")
	if err != nil {
		log.Fatalf("Error getting groups for user %s: %+v", "user", err)
	}
	log.Printf("Groups: %+v", groups)
}
