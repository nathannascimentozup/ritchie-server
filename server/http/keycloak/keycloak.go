package keycloak

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ritchie-server/server"
)

type Handler struct {
	Config server.Config
}

func NewKeycloakHandler(config server.Config) server.DefaultHandler {
	return Handler{Config: config}
}

func (lh Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var method = r.Method
		if http.MethodGet == method {
			lh.processGet(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (lh Handler) processGet(w http.ResponseWriter, r *http.Request) {
	organizationHeader := r.Header.Get(server.OrganizationHeader)

	keycloakConfigs, err := lh.Config.ReadKeycloakConfigs(organizationHeader)

	if err != nil {
		log.Error("KeycloakConfig not found for organization ", organizationHeader)
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(*keycloakConfigs)

}
