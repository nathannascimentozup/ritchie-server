package keycloak

import (
	"fmt"
	"github.com/google/uuid"
	"ritchie-server/server"
	"ritchie-server/server/mock"
	"testing"
)

func TestManager_CreateUser(t *testing.T) {
	type fields struct {
		c server.Config
	}
	type args struct {
		user server.CreateUser
		org  string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    bool
		outErr bool
	}{
		{
			name: "success",
			fields: fields{c: mock.DummyConfig()},
			in: args{
				user: randomCreateUser(),
				org:  "zup",
			},
			out:    true,
			outErr: false,
		},
		{
			name: "failed read config",
			fields: fields{c: mock.DummyConfig()},
			in: args{
				user: randomCreateUser(),
				org:  "notfound",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "client login failed",
			fields: fields{c: mock.DummyConfig("", "", "", "failed")},
			in: args{
				user: randomCreateUser(),
				org:  "zup",
			},
			out:    false,
			outErr: true,
		},
		{
			name: "failed create",
			fields: fields{c: mock.DummyConfig()},
			in: args{
				user: server.CreateUser{
					Username:  "user",
					Password:  "admin",
					FirstName: "",
					LastName:  "",
					Email:     "user@user.com",
				},
				org: "zup",
			},
			out:    false,
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewKeycloakManager(tt.fields.c)
			got, err := m.CreateUser(tt.in.user, tt.in.org)
			if (err != nil) != tt.outErr {
				t.Errorf("CreateUser() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if (len(got) > 0) != tt.out {
				t.Errorf("CreateUser() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func TestManager_DeleteUser(t *testing.T) {
	type fields struct {
		c server.Config
	}
	type args struct {
		user server.CreateUser
		org  string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		outErr bool
	}{
		{
			name: "success",
			fields: fields{c: mock.DummyConfig()},
			in: args{
				user: createUserTestForDelete(),
				org:  "zup",
			},
			outErr: false,
		},
		{
			name: "without permission",
			fields: fields{c: mock.DummyConfig(
				"",
				"",
				"user-login-without-permission",
				"c0b1f8f1-c746-4573-9679-b7ecdf2d48cc")},
			in: args{
				user: createUserTestForDelete(),
				org:  "zup",
			},
			outErr: true,
		},
		{
			name: "user not found",
			fields: fields{c: mock.DummyConfig()},
			in: args{
				user: randomCreateUser(),
				org:  "zup",
			},
			outErr: true,
		},
		{
			name: "failed read config",
			fields: fields{c: mock.DummyConfig()},
			in: args{
				user: createUserTestForDelete(),
				org:  "notfound",
			},
			outErr: true,
		},
		{
			name: "client login failed",
			fields: fields{c: mock.DummyConfig("", "", "", "failed")},
			in: args{
				user: randomCreateUser(),
				org:  "zup",
			},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewKeycloakManager(tt.fields.c)
			if err := m.DeleteUser(tt.in.org, tt.in.user.Email); (err != nil) != tt.outErr {
				t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.outErr)
			}
		})
	}
}

func TestManager_Login(t *testing.T) {
	type fields struct {
		c server.Config
	}
	type args struct {
		org      string
		user     string
		password string
	}
	tests := []struct {
		name     string
		fields   fields
		in     args
		outToken bool
		outCod   int
		outErr   bool
	}{
		{
			name: "success",
			fields: fields{c: mock.DummyConfig()},
			in:     args{
				org:      "zup",
				user:     "user",
				password: "admin",
			},
			outToken: true,
			outCod:   0,
			outErr:   false,
		},
		{
			name: "config not found",
			fields: fields{c: mock.DummyConfig()},
			in:     args{
				org:      "not",
				user:     "user",
				password: "admin",
			},
			outToken: false,
			outCod:   404,
			outErr:   true,
		},
		{
			name: "wrong user",
			fields: fields{c: mock.DummyConfig()},
			in:     args{
				org:      "zup",
				user:     "user",
				password: "wrong",
			},
			outToken: false,
			outCod:   401,
			outErr:   true,
		},
		{
			name: "invalid client",
			fields: fields{c: mock.DummyConfig(
				"",
				"",
				"",
				"invalid")},
			in:     args{
				org:      "zup",
				user:     "user",
				password: "admin",
			},
			outToken: false,
			outCod:   400,
			outErr:   true,
		},
		{
			name: "invalid keycloak url",
			fields: fields{c: mock.DummyConfig(
				"invalid",
				"",
				"",
				"")},
			in:     args{
				org:      "zup",
				user:     "user",
				password: "admin",
			},
			outToken: false,
			outCod:   500,
			outErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Manager{
				config: tt.fields.c,
			}
			gotToken, gotCode, err := m.Login(tt.in.org, tt.in.user, tt.in.password)
			if (err != nil) != tt.outErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.outErr)
				return
			}
			if (len(gotToken) > 0) != tt.outToken {
				t.Errorf("Login() got = %v, want %v", gotToken, tt.outToken)
			}
			if gotCode != tt.outCod {
				t.Errorf("Login() got1 = %v, want %v", gotCode, tt.outCod)
			}
		})
	}
}

func randomCreateUser() server.CreateUser {
	return server.CreateUser{
		Username:  uuid.New().String(),
		Password:  uuid.New().String(),
		FirstName: uuid.New().String(),
		LastName:  uuid.New().String(),
		Email:     uuid.New().String() + "@test.com",
	}
}

func createUserTestForDelete() server.CreateUser {
	m := NewKeycloakManager(mock.DummyConfig())
	u := randomCreateUser()
	_, err := m.CreateUser(u, "zup")
	if err != nil {
		fmt.Sprintln("Error in Create User ")
		return server.CreateUser{}
	}
	return u
}
