package hello

import (
	"net/http"
	"ritchie-server/server"
)

type Handler struct {
}

func NewHelloHandler() server.DefaultHandler {
	return Handler{}
}

func (hh Handler) Handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
		}
	})
}
