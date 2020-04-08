package credential

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestConfigHandler_Handler(t *testing.T) {
	type fields struct {
		config server.Config
		org    string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name: "success getting config",
			fields: fields{
				config: mock.DummyConfig(),
				org:    "zup",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}
			}(),
		},
		{
			name: "config not found",
			fields: fields{
				config: mock.DummyConfig(),
				org:    "not",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConfigHandler(tt.fields.config)

			r, _ := http.NewRequest(http.MethodGet, "/test", bytes.NewReader([]byte{}))
			r.Header.Add(server.OrganizationHeader, tt.fields.org)

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			c.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v out %v", g.Code, w.Code)
			}
		})
	}
}
