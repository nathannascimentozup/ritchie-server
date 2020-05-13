package security

import (
	"github.com/Nerzal/gocloak"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestAuthorization_AuthorizationPath(t *testing.T) {
	type fields struct {
		Config server.Config
	}
	type args struct {
		bearerToken string
		path        string
		method      string
		org         string
	}
	tests := []struct {
		name    string
		fields  fields
		in      args
		out    bool
		outErr bool
	}{
		{
			name:    "empty org",
			fields:  fields{},
			in:      args{},
			out:    false,
			outErr: true,
		},
		{
			name: "empty token",
			in: args{
				bearerToken: "",
				org:         "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "invalid token",
			in: args{
				bearerToken: "invalid",
				org:         "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "empty jwt",
			in: args{
				bearerToken: "Bearer ",
				org:         "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "keycloak config not found",
			fields: fields{
				Config: mock.DummyConfig(),
			},
			in: args{
				bearerToken: "Bearer " + generateAccessTokenAdmin(mock.DummyConfig()),
				org:         "notfound",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "Error decode token",
			fields: fields{
				Config: mock.DummyConfig(),
			},
			in: args{
				bearerToken: "Bearer  invalid",
				org:         "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "Validate constrains",
			fields: fields{
				/*Config: dummyConfig(dummyKeycloakConfig(), server.SecurityConstraints{
					Constraints: []server.DenyMatcher{{
						Pattern:      "/test",
						RoleMappings: map[string][]string{"user": {"POST", "GET"}},
					}},
				}),*/
				Config:mock.DummyConfig(),
			},
			in: args{
				bearerToken: "Bearer " + generateAccessTokenAdmin(mock.DummyConfig()),
				org:         "zup",
				method:      "GET",
				path:        "/validate",
			},
			out:    true,
			outErr: false,
		},
		{
			name: "Invalid constrains",
			fields: fields{
				Config:mock.DummyConfig(),
			},
			in: args{
				bearerToken: "Bearer " + generateAccessTokenAdmin(mock.DummyConfig()),
				org:         "zup",
				method:      "GET",
				path:        "/test",
			},
			out:    false,
			outErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAuthorization(tt.fields.Config)
			got, err := auth.AuthorizationPath(tt.in.bearerToken, tt.in.path, tt.in.method, tt.in.org)
			if (err != nil) != tt.outErr {
				t.Errorf("AuthorizationPath() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if got != tt.out {
				t.Errorf("AuthorizationPath() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func TestAuthorization_ValidatePublicConstraints(t *testing.T) {
	type fields struct {
		Config              server.Config
	}
	type args struct {
		path   string
		method string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		out   bool
	}{
		{
			name: "Validate constrains",
			fields: fields{
				Config: mock.DummyConfig(),
			},
			args: args{
				method: "GET",
				path:   "/public",
			},
			out: true,
		},
		{
			name: "Invalid constrains",
			fields: fields{
				Config: mock.DummyConfig(),
			},
			args: args{
				method: "GET",
				path:   "/invalid",
			},
			out: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := NewAuthorization(tt.fields.Config)
			if got := auth.ValidatePublicConstraints(tt.args.path, tt.args.method); got != tt.out {
				t.Errorf("ValidatePublicConstraints() = %v, out %v", got, tt.out)
			}
		})
	}
}

func generateAccessTokenAdmin(configs server.Config) string {
	kc, _ := configs.ReadKeycloakConfigs("zup")
	client := gocloak.NewClient(kc.Url)
	token, _ := client.Login(kc.ClientId, kc.ClientSecret, kc.Realm, "user", "admin")
	return token.AccessToken
}
