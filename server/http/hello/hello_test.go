package hello

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_Handler(t *testing.T) {
	type fields struct {
		Path   string
	}
	tests := []struct {
		name string
		fields fields
		out  http.HandlerFunc

	}{
		{
			name: "ok",
			fields: fields{Path: "/"},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
				}
			}(),
		},
		{
			name: "not found",
			fields: fields{Path: "/notfound"},
			out: func() http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := NewHelloHandler()

			r, _ := http.NewRequest(http.MethodGet, tt.fields.Path, bytes.NewReader([]byte{}))

			w := httptest.NewRecorder()

			tt.out.ServeHTTP(w, r)

			g := httptest.NewRecorder()

			h.Handler().ServeHTTP(g, r)

			if g.Code != w.Code {
				t.Errorf("Handler returned wrong status code: got %v want %v", g.Code, w.Code)
			}
		})
	}
}

func TestNewHelloHandler(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "correct type",
			want: "http.HandlerFunc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{}
			if got := reflect.TypeOf(h.Handler()); fmt.Sprint(got) != tt.want {
				t.Errorf("Handler() = %v, want %v", got, tt.want)
			}
		})
	}
}