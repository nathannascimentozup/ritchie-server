package credential

import (
	"encoding/json"
	"net/http"

	"ritchie-server/server"

	log "github.com/sirupsen/logrus"
)

func (h Handler) HandleAdmin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.processAdminPost(w, r)
		default:
			http.NotFound(w, r)
		}
	}
}

func (h Handler) processAdminPost(w http.ResponseWriter, r *http.Request) {
	org := org(r)
	ctx := ctx(r)
	var c server.Credential
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		log.Error("Failed to process request ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	validationError := h.defaultValidate(c, org)
	adminValidate(c, validationError)

	if len(validationError) > 0 {
		err := map[string]interface{}{"validationError": validationError}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	err := h.createCredential(org, ctx, c)
	if err != nil {
		log.Error("Failed to create credential ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func adminValidate(c server.Credential, validationError map[string]string) {
	if c.Username == "" {
		validationError["username"] = "The username field is required."
	}
}
