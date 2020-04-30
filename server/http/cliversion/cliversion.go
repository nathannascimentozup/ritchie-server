package cliversion

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"ritchie-server/server"
)

type Handler struct {
	Config           server.Config
}

func NewConfigHandler(config server.Config) server.DefaultHandler {
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

	var client http.Client
	var bodyString string

	organizationHeader := r.Header.Get(server.OrganizationHeader)
	cliVersionConfigs, err := lh.Config.ReadCliVersionConfigs(organizationHeader)

	if err != nil {
		log.Error("Organization not found ", organizationHeader)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if len(cliVersionConfigs.Url) == 0 {
		log.Error("Cli Version File not Configured")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := client.Get(cliVersionConfigs.Url)
	if err != nil || resp.StatusCode == http.StatusNotFound {
		cliVersionConfigs.Version = ""
		log.Error("No cli version found on ", cliVersionConfigs.Url)
		log.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString = string(bodyBytes)
	}
	cliVersionConfigs.Version = bodyString

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(cliVersionConfigs)
	if err != nil {
		fmt.Sprintln("Error in Json Encode ")
		return
	}

}
