package login

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"ritchie-server/server"
)

type Handler struct {
	KeyCloakManager server.KeycloakManager
}

func NewLoginHandler(k server.KeycloakManager) server.DefaultHandler {
	return Handler{KeyCloakManager: k}
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type response struct {
	Token string `json:"access_token"`
}

func (lh Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var method = r.Method
		if http.MethodPost == method {
			lh.processPost(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (lh Handler) processPost(w http.ResponseWriter, r *http.Request) {
	var l login
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&l)
	if err != nil {
		log.Error("Failed to process Json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	organizationHeader := r.Header.Get(server.OrganizationHeader)
	if validErrs := l.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}
	token, code, err := lh.KeyCloakManager.Login(organizationHeader, l.Username, l.Password)
	if err != nil {
		w.WriteHeader(code)
		return
	}
	objectResponse := response{Token: token}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(objectResponse)
}


func (l login) validate() url.Values {
	errs := url.Values{}
	if l.Username == "" {
		errs.Add("username", "The username field is required!")
	}
	if l.Password == "" {
		errs.Add("password", "The password field is required!")
	}
	return errs
}
