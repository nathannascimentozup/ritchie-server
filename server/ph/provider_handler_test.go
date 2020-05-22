package ph

import (
	"reflect"
	"testing"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestHandler_FindRepo(t *testing.T) {
	type fields struct {
		authorization server.Constraints
	}
	type args struct {
		repos         []server.Repository
		authorization server.Constraints
		repoName      string
	}
	tests := []struct {
		name   string
		fields fields
		in     args
		out    server.Repository
		outErr bool
	}{
		{
			name: "find commons",
			in: args{
				repos:    mock.DummyRepoList(),
				repoName: "commons",
			},
			out:    mock.DummyRepoList()[0],
			outErr: false,
		},
		{
			name: "repo not found",
			in: args{
				repos:    mock.DummyRepoList(),
				repoName: "notfound",
			},
			out:    server.Repository{},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hp := Handler{
				authorization: tt.fields.authorization,
			}
			got, err := hp.FindRepo(tt.in.repos, tt.in.repoName)
			if (err != nil) != tt.outErr {
				t.Errorf("FindRepo() error = %v, wantErr %v", err, tt.outErr)
				return
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("FindRepo() got = %v, want %v", got, tt.out)
			}
		})
	}
}