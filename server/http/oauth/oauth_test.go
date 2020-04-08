package oauth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		c server.Config
		url string
		org    string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name:   "success",
			fields: fields{
				c:       mock.DummyConfig(),
				url:     "/oauth",
				org:     "zup",
			},
			out:    func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-type", "application/json")
				}
			}(),
		},
		{
			name:   "conf not found",
			fields: fields{
				c:       mock.DummyConfig(),
				url:     "/oauth",
				org:     "not",
			},
			out:    func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Header().Set("Content-type", "text/plain; charset=utf-8")
				}
			}(),
		},
		{
			name:   "url not found",
			fields: fields{
				c:       mock.DummyConfig(),
				url:     "/not",
				org:     "zup",
			},
			out:    func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Header().Set("Content-type", "text/plain; charset=utf-8")
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewConfigHandler(tt.fields.c)
			r, _ := http.NewRequest(http.MethodGet, tt.fields.url, bytes.NewReader([]byte{}))
			r.Header.Add(server.OrganizationHeader, tt.fields.org)

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}

			if g.Header().Get("Content-Type") != w.Header().Get("Content-Type") {
				t.Errorf("Wrong content type. Got %v want %v", g.Header().Get("Content-Type"), w.Header().Get("Content-Type"))
			}
		})
	}
}