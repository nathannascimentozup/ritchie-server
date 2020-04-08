package credential

import (
	"encoding/json"
	"net/http"
	"ritchie-server/server"
)

type ConfigHandler struct {
	Config server.Config
}

func NewConfigHandler(config server.Config) server.DefaultHandler {
	return ConfigHandler{Config: config}
}

func (c ConfigHandler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			c.processGet(w, r)
		}
	}
}

func (c ConfigHandler) processGet(w http.ResponseWriter, r *http.Request) {
	organizationHeader := r.Header.Get(server.OrganizationHeader)

	credentialConfigs, err := c.Config.ReadCredentialConfigs(organizationHeader)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-type", "application/json")
	_ = json.NewEncoder(w).Encode(credentialConfigs)
}
