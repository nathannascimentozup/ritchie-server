package keycloak

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		c      server.Config
		org    string
		method string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name: "success keycloak configuration",
			fields: fields{
				c:      mock.DummyConfig(),
				org:    "zup",
				method: http.MethodGet,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-type", "application/json")
					json.NewEncoder(w).Encode(keycloakConfigWant())
				}
			}(),
		},
		{
			name: "notfound  keycloak configuration",
			fields: fields{
				c:      mock.DummyConfig(),
				org:    "notfound",
				method: http.MethodGet,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-type", "text/plain; charset=utf-8")
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
		{
			name: "method not found",
			fields: fields{
				c:      mock.DummyConfig(),
				org:    "notfound",
				method: http.MethodPost,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-type", "text/plain; charset=utf-8")
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewKeycloakHandler(tt.fields.c)
			r, _ := http.NewRequest(tt.fields.method, "/test", bytes.NewReader([]byte{}))
			r.Header.Add(server.OrganizationHeader, tt.fields.org)

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v out %v", g.Code, w.Code)
			}

			if g.Header().Get("Content-Type") != w.Header().Get("Content-Type") {
				t.Errorf("Wrong content type. Got %v out %v", g.Header().Get("Content-Type"), w.Header().Get("Content-Type"))
			}

			if len(g.Body.String()) > 0 {
				var got, want server.KeycloakConfig
				json.Unmarshal(g.Body.Bytes(), &got)
				json.Unmarshal(w.Body.Bytes(), &want)
				if got != want {
					t.Errorf("Wrong version. Got %v out %v", got, want)
				}
			}
		})
	}
}

func keycloakConfigWant() server.KeycloakConfig {
	conf, _ := mock.DummyConfig().ReadKeycloakConfigs("zup")
	return *conf
}
