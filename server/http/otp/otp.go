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
		if r.URL.Path != "/otp" {
			http.NotFound(w, r)
		}
		oh.processRequest(w, r)
	}
}

func (oh Handler) processRequest(w http.ResponseWriter, r *http.Request) {
	org := r.Header.Get("x-org")
	_, existence := oh.securityProviders.Providers[org]
	if existence == false {
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
		log.Println("Error in Json Encode")
	}
	return
}
