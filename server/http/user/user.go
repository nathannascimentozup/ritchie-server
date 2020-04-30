package user

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"ritchie-server/server"
)

type Handler struct {
	keycloakManager server.KeycloakManager
	vaultManager    server.VaultManager
}

func NewUserHandler(k server.KeycloakManager, v server.VaultManager) server.DefaultHandler {
	return Handler{keycloakManager: k, vaultManager: v}
}

type (
	deleteUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
)

func (uh Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			uh.processPost(w, r)
		case http.MethodDelete:
			uh.processDelete(w, r)
		default:
			http.NotFound(w, r)
		}
	}
}

func (uh Handler) processPost(w http.ResponseWriter, r *http.Request) {
	var u server.CreateUser
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error("Failed to process Json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	if validErrs := validate(u); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	organizationHeader := r.Header.Get(server.OrganizationHeader)

	_, err = uh.keycloakManager.CreateUser(u, organizationHeader)
	if err != nil {
		log.Error("Failed to create user ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (uh Handler) processDelete(w http.ResponseWriter, r *http.Request) {
	var deleteUser deleteUser
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&deleteUser)
	if err != nil {
		log.Error("Failed to process Json ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	if validErrs := deleteUser.validate(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	org := r.Header.Get(server.OrganizationHeader)
	err = uh.keycloakManager.DeleteUser(org, deleteUser.Email)
	if err != nil {
		log.Error("Failed to delete user ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}
	err = uh.deleteVaultCredentials(org, deleteUser)
	if err != nil {
		log.Error("Failed to delete user ", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (uh Handler) deleteVaultCredentials(org string, deleteUser deleteUser) error {
	credentialKey := fmt.Sprintf("%s/%s", org, deleteUser.Username)
	keys, err := uh.vaultManager.List(credentialKey)
	if err != nil {
		return err
	}

		for i := range keys {
			credentialPath := fmt.Sprintf("%s/%s", credentialKey, keys[i])
			err := uh.vaultManager.Delete(credentialPath)
			if err != nil {
				return err
			}
		}

	return nil
}

func validate(u server.CreateUser) url.Values {
	errs := url.Values{}
	if u.Username == "" {
		errs.Add("username", "The username field is required!")
	}
	if u.Password == "" {
		errs.Add("password", "The password field is required!")
	}
	if u.Email == "" {
		errs.Add("email", "The email field is required!")
	}
	if u.FirstName == "" {
		errs.Add("firstName", "The firstName field is required!")
	}
	return errs
}

func (u deleteUser) validate() url.Values {
	errs := url.Values{}
	if u.Username == "" {
		errs.Add("username", "The username field is required!")
	}
	if u.Email == "" {
		errs.Add("email", "The email field is required!")
	}
	return errs
}
