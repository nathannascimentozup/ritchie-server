package user

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		k       server.KeycloakManager
		v       server.VaultManager
		method  string
		org     string
		payload string
	}
	tests := []struct {
		name   string
		fields fields
		out    http.HandlerFunc
	}{
		{
			name: "create user",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				method:  http.MethodPost,
				org:     "zup",
				payload: `{"email" : "teste@test.com","firstName" : "test","lastName" : "test","organization" : "test","password" : "test","username" : "test"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusCreated)
				}
			}(),
		},
		{
			name: "method not found",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodPatch,
				org:     "zup",
				payload: `{"email" : "teste@test.com","firstName" : "test","lastName" : "test","organization" : "test","password" : "test","username" : "test"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Header().Set("Content-type", "text/plain; charset=utf-8")
				}
			}(),
		},
		{
			name: "create invalid payload",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodPost,
				org:     "zup",
				payload: `12312431`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "create empty fields payload",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodPost,
				org:     "zup",
				payload: `{}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set("Content-type", "application/json")
				}
			}(),
		},
		{
			name: "failed create user",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   errors.New("error"),
				},
				method:  http.MethodPost,
				org:     "zup",
				payload: `{"email" : "teste@test.com","firstName" : "test","lastName" : "test","organization" : "test","password" : "test","username" : "test"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "delete user",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodDelete,
				org:     "zup",
				payload: `{"email" : "teste@test.com", "username" : "test"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}
			}(),
		},
		{
			name: "delete invalid payload",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodDelete,
				org:     "zup",
				payload: `12312431`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "delete empty fields payload",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: nil,
					Keys:    nil,
				},
				method:  http.MethodDelete,
				org:     "zup",
				payload: `{}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusBadRequest)
					w.Header().Set("Content-type", "application/json")
				}
			}(),
		},
		{
			name: "failed delete user",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   errors.New("error"),
				},
				method:  http.MethodDelete,
				org:     "zup",
				payload: `{"email" : "teste@test.com","username" : "test"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "failed delete vault",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     errors.New("error"),
					ErrList: nil,
					Keys:    []interface{}{"first", "second"},
				},
				method:  http.MethodDelete,
				org:     "zup",
				payload: `{"email" : "teste@test.com", "username" : "test"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
		{
			name: "failed delete list vault",
			fields: fields{
				k: mock.KeycloakMock{
					Token: "123",
					Code:  0,
					Err:   nil,
				},
				v: mock.VaultMock{
					Err:     nil,
					ErrList: errors.New("error"),
					Keys:    []interface{}{"first", "second"},
				},
				method:  http.MethodDelete,
				org:     "zup",
				payload: `{"email" : "teste@test.com", "username" : "test"}`,
			},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewUserHandler(tt.fields.k, tt.fields.v)
			var b []byte
			if len(tt.fields.payload) > 0 {
				b = append(b, []byte(tt.fields.payload)...)
			}
			r, _ := http.NewRequest(tt.fields.method, "/test", bytes.NewReader(b))
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
