package otp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		securityProviders server.SecurityProviders
		org               string
		method            string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name: "success",
			fields: fields{
				securityProviders: server.SecurityProviders{
					Providers: map[string]server.SecurityManager{
						"zup": mock.SecurityManagerMock{
							O: true,
						},
					},
				},
				method: http.MethodGet,
				org:    "zup",
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}
			}(),
		},
		{
			name: "organization not found",
			fields: fields{
				securityProviders: server.SecurityProviders{
					Providers: map[string]server.SecurityManager{
						"zup": mock.SecurityManagerMock{
							O: true,
						},
					},
				},
				method: http.MethodGet,
				org:    "mock",
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
			oh := Handler{
				securityProviders: tt.fields.securityProviders,
			}
			r, _ := http.NewRequest(tt.fields.method, "/otp", nil)
			r.Header.Add(server.OrganizationHeader, tt.fields.org)

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			oh.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v\", g.Code, w.Code", g.Code, w.Code)
			}
		})
	}
}
