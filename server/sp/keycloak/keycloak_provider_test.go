package keycloak

import (
	"reflect"
	"testing"

	"ritchie-server/server"
)

func Test_keycloakConfig_Login(t *testing.T) {
	type fields struct {
		config map[string]string
	}
	type args struct {
		username string
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		outUser   server.User
		outError  server.LoginError
	}{
		{
			name:     "login success",
			fields:   fields{config: map[string]string{
				"type" : "keycloak",
				"url": "http://localhost:8080",
				"realm": "ritchie",
				"clientId": "user-login",
				"clientSecret": "user-login",
			}},
			args:     args{
				username: "user",
				password: "admin",
			},
			outUser:  keycloakUser{
				roles:    []string{"admin", "offline_access", "uma_authorization", "user"},
				userInfo: server.UserInfo{
					Name:     "user user",
					Username: "user",
					Email:    "user@user.com",
				},
			},
			outError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeycloakProvider(tt.fields.config)
			gotUser, gotError := k.Login(tt.args.username, tt.args.password)
			if !reflect.DeepEqual(gotUser.UserInfo(), tt.outUser.UserInfo()) {
				t.Errorf("Login() gotUser.UserInfo() = %v, want %v", gotUser.UserInfo(), tt.outUser.UserInfo())
			}
			roles := make(map[string]string)
			for _, c := range gotUser.Roles() {
				roles[c] = c
			}
			for _, c := range tt.outUser.Roles() {
				if roles[c] == "" {
					t.Errorf("Error roles gotUser.Roles() = %v, want %v", gotUser.Roles(), tt.outUser.Roles())
				}
			}
			if !reflect.DeepEqual(gotError, tt.outError) {
				t.Errorf("Login() gotError = %v, want %v", gotError, tt.outError)
			}
		})
	}
}