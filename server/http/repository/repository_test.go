package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/config"
	"ritchie-server/server/mock"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		Config server.Config
		org    string
		method string
	}
	tests := []struct {
		name   string
		fields fields
		want   http.HandlerFunc
	}{
		{
			name: "success",
			fields: fields{
				Config: mock.DummyConfig(),
				org:    "zup",
				method: http.MethodGet,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-type", "application/json")
					err := json.NewEncoder(w).Encode(repositoryConfigWant())
					if err != nil {
						fmt.Sprintln("Error in Encode Json ")
						return
					}
				}
			}(),
		},
		{
			name: "not found",
			fields: fields{
				Config: mock.DummyConfig(),
				org:    "zup",
				method: http.MethodPost,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
		{
			name: "not found org",
			fields: fields{
				Config: mock.DummyConfig(),
				org:    "notfound",
				method: http.MethodGet,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
		{
			name: "nil config",
			fields: fields{
				Config: config.Configuration{
					Configs: map[string]*server.ConfigFile{
						"empty": {
							RepositoryConfig: nil,
						}},
				},
				org:    "empty",
				method: http.MethodGet,
			},
			want: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := NewConfigHandler(tt.fields.Config)

			r, _ := http.NewRequest(tt.fields.method, "/usage", bytes.NewReader([]byte{}))

			r.Header.Add(server.OrganizationHeader, tt.fields.org)
			r.Header.Add("Content-Type", "application/json")

			w := httptest.NewRecorder()

			tt.want.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			mu.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}

			if g.Code == http.StatusOK {
				if !reflect.DeepEqual(g.Body, w.Body) {
					t.Errorf("Handler returned wrong body: got %v \n want %v", g.Body, w.Body)
				}
			}
		})
	}
}

func repositoryConfigWant() []response {
	var resp []response
	conf, _ := mock.DummyConfig().ReadRepositoryConfig("zup")
	for _,r := range conf {
		r := response{
			Name:     r.Name,
			Priority: r.Priority,
			TreePath: r.ServerUrl + r.TreePath,
			Username: r.Username,
			Password: r.Password,
		}
		resp = append(resp, r)
	}
	return resp
}
