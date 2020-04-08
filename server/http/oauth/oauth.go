package oauth

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ritchie-server/server"
)

type Handler struct {
	Config server.Config
}

func NewConfigHandler(config server.Config) server.DefaultHandler {
	return Handler{Config: config}
}

func (ch Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/oauth" {
			http.NotFound(w, r)
		} else {
			orgHeader := r.Header.Get(server.OrganizationHeader)

			oc, err := ch.Config.ReadOauthConfig(orgHeader)
			if err != nil {
				log.Error("OauthConfig not found for organization ", orgHeader)
				http.NotFound(w, r)
				return
			}
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(oc)
		}
	}
}
