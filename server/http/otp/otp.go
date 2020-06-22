package otp

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"ritchie-server/server"
)

type Handler struct {
	securityProviders server.SecurityProviders
}

type response struct {
	Otp bool `json:"otp"`
}

func NewOtpHandler(sp server.SecurityProviders) server.DefaultHandler {
	return Handler{
		securityProviders: sp,
	}
}

func (oh Handler) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			oh.processRequest(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
}

func (oh Handler) processRequest(w http.ResponseWriter, r *http.Request) {
	org := r.Header.Get(server.OrganizationHeader)
	_, existence := oh.securityProviders.Providers[org]
	if !existence {
		log.Printf("Organization {%s} not found", org)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hasOtp := oh.securityProviders.Providers[org].Otp()
	resp := response{
		Otp: hasOtp,
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(&resp)
	if err != nil {
		log.Error("Error in Json Encode")
	}
}
