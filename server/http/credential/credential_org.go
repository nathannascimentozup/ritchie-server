package credential

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"ritchie-server/server"
)

func (h Handler) HandleOrg() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.processOrgPost(w, r)
		default:
			w.Header().Add("Allow", http.MethodPost)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (h Handler) processOrgPost(w http.ResponseWriter, r *http.Request) {
	org := org(r)
	ctx := ctx(r)

	var orgCred server.Credential
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&orgCred); err != nil {
		log.Error("Failed to process request ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	if err := h.defaultValidate(orgCred, org); len(err) > 0 {
		err := map[string]interface{}{"validationError": err}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	if err := h.createOrgCredential(org, ctx, orgCred); err != nil {
		log.Error("Failed to create org credential ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) createOrgCredential(org server.Org, ctx server.Ctx, c server.Credential) error {
	path := fmt.Sprintf(orgCredentialVaultPath, ctxResolver(org, ctx), c.Service)
	if err := h.v.Write(path, c.Credential); err != nil {
		return err
	}
	return nil
}
