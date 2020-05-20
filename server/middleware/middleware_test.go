package middleware

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"ritchie-server/server"
	"ritchie-server/server/mock"

	"testing"
)

func TestHandler_Filter(t *testing.T) {

	type fields struct {
		SecurityConstraints server.SecurityConstraints
		Authorization       server.Constraints
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name   string
		path   string
		fields fields
		in     args
		out   http.Handler
	}{
		{
			name: "public path",
			path: "/test",
			fields: fields{SecurityConstraints: server.SecurityConstraints{
				PublicConstraints: []server.PermitMatcher{{
					Pattern: "/test",
					Methods: []string{"GET"},
				}},
			},
				Authorization: mock.AuthorizationMock{
					B:  true,
					E: nil,
				},
			},
			in: args{next: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/test" {
						http.NotFound(w, r)
					}
				}
			}(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/test" {
						http.NotFound(w, r)
					}
				}
			}(),
		},
		{
			name: "authorized",
			path: "/test",
			fields: fields{
				SecurityConstraints: server.SecurityConstraints{
					Constraints: []server.DenyMatcher{{
						Pattern:      "/test",
						RoleMappings: map[string][]string{"user": {"POST", "GET"}},
					}},
					PublicConstraints: []server.PermitMatcher{},
				},
				Authorization: mock.AuthorizationMock{
					B:  true,
					E: nil,
				},
			},
			in: args{next: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/test" {
						http.NotFound(w, r)
					}
				}
			}(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusOK)
				}
			}(),
		},
		{
			name: "not authorized",
			path: "/test",
			fields: fields{
				SecurityConstraints: server.SecurityConstraints{
					Constraints: []server.DenyMatcher{{
						Pattern:      "/test",
						RoleMappings: map[string][]string{"user": {"POST", "GET"}},
					}},
					PublicConstraints: []server.PermitMatcher{},
				},
				Authorization: mock.AuthorizationMock{
					B:  false,
					E: nil,
				},
			},
			in: args{next: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/test" {
						http.NotFound(w, r)
					}
				}
			}(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusForbidden)
				}
			}(),
		},
		{
			name: "authorization error",
			path: "/test",
			fields: fields{
				SecurityConstraints: server.SecurityConstraints{
					Constraints: []server.DenyMatcher{{
						Pattern:      "/test",
						RoleMappings: map[string][]string{"user": {"POST", "GET"}},
					}},
					PublicConstraints: []server.PermitMatcher{},
				},
				Authorization: mock.AuthorizationMock{
					B:  false,
					E: errors.New("error"),
				},
			},
			in: args{next: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path != "/test" {
						http.NotFound(w, r)
					}
				}
			}(),
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusUnauthorized)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := NewMiddlewareHandler(tt.fields.Authorization)

			r, _ := http.NewRequest(http.MethodGet, tt.path, bytes.NewReader([]byte{}))

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			mh.Filter(tt.in.next).ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}
		})
	}
}
