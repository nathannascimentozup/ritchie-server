package tm

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"

	"ritchie-server/server"
	"ritchie-server/server/mock"
)

func TestTreeRemoteAllow(t *testing.T) {
	type args struct {
		sec    server.Constraints
		bToken string
		org    string
		tPath  string
		repo   server.Repository
	}
	tests := []struct {
		name   string
		in     args
		out    server.Tree
		outErr bool
	}{
		{
			name: "response tree with allow commands",
			in: args{
				sec: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				bToken: "",
				org:    "",
				tPath:  "/tree/tree.json",
				repo:   mock.DummyRepo(),
			},
			out:    treeRoleUser(),
			outErr: false,
		},
		{
			name: "list roles error",
			in: args{
				sec: mock.AuthorizationMock{
					B: false,
					E: errors.New("error"),
					R: nil,
				},
				bToken: "",
				org:    "",
				tPath:  "/tree/tree.json",
				repo:   mock.DummyRepo(),
			},
			out:    server.Tree{},
			outErr: true,
		},
		{
			name: "tree not found",
			in: args{
				sec: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				bToken: "",
				org:    "",
				tPath:  "/tree/tree-notfound.json",
				repo:   mock.DummyRepo(),
			},
			out:    server.Tree{},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TreeRemoteAllow(tt.in.sec, tt.in.bToken, tt.in.org, tt.in.tPath, tt.in.repo)
			if (err != nil) != tt.outErr {
				t.Errorf("TreeRemoteAllow() error = %v, outErr %v", err, tt.outErr)
				return
			}
			commands := make(map[string]*server.Command)
			for _, c := range got.Commands {
				commands[c.Parent+c.Usage] = &c
			}
			for _, c := range tt.out.Commands {
				if commands[c.Parent+c.Usage] == nil {
					t.Errorf("Commands receive in tree error gotT = %v, outT %v", got, tt.out)
				}
			}
		})
	}
}

func TestFormulaAllow(t *testing.T) {
	type args struct {
		sec   server.Constraints
		fPath string
		token string
		org   string
		repo  server.Repository
	}
	tests := []struct {
		name   string
		in     args
		out    bool
		outErr bool
	}{
		{
			name:   "allow",
			in:     args{
				sec: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				fPath: "/formulas/aws/terraform/config.json",
				token: "",
				org:   "",
				repo:  mock.DummyRepo(),
			},
			out:    true,
			outErr: false,
		},
		{
			name:   "not allow",
			in:     args{
				sec: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"NO_RULE"},
				},
				fPath: "/formulas/aws/terraform/config.json",
				token: "",
				org:   "",
				repo:  mock.DummyRepo(),
			},
			out:    false,
			outErr: false,
		},
		{
			name:   "load roles error",
			in:     args{
				sec: mock.AuthorizationMock{
					B: true,
					E: errors.New("error"),
					R: []string{},
				},
				fPath: "/formulas/aws/terraform/config.json",
				token: "",
				org:   "",
				repo:  mock.DummyRepo(),
			},
			out:    false,
			outErr: true,
		},
		{
			name:   "allow formula without role",
			in:     args{
				sec: mock.AuthorizationMock{
					B: true,
					E: nil,
					R: []string{"USER"},
				},
				fPath: "/formulas/scaffold/coffee-go/config.json",
				token: "",
				org:   "",
				repo:  mock.DummyRepo(),
			},
			out:    true,
			outErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormulaAllow(tt.in.sec, tt.in.fPath, tt.in.token, tt.in.org, tt.in.repo)
			if (err != nil) != tt.outErr {
				t.Errorf("FormulaAllow() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if got != tt.out {
				t.Errorf("FormulaAllow() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func TestFindRepo(t *testing.T) {
	type args struct {
		repos    []server.Repository
		repoName string
	}
	tests := []struct {
		name   string
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
			out: mock.DummyRepoList()[0],
			outErr: false,
		},
		{
			name: "repo not found",
			in: args{
				repos:    mock.DummyRepoList(),
				repoName: "notfound",
			},
			out: server.Repository{},
			outErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindRepo(tt.in.repos, tt.in.repoName)
			if (err != nil) != tt.outErr {
				t.Errorf("FindRepo() error = %v, outErr %v", err, tt.outErr)
				return
			}
			if !reflect.DeepEqual(got, tt.out) {
				t.Errorf("FindRepo() got = %v, out %v", got, tt.out)
			}
		})
	}
}

func treeRoleUser() server.Tree {
	js := `{
        "commands": [
          {
            "usage": "aws",
            "help": "Apply Aws objects",
            "parent": "root",
            "roles" : ["USER"]
          },
          {
            "usage": "apply",
            "help": "Apply Aws objects",
            "parent": "root_aws",
            "roles" : ["USER"]
          },
          {
            "usage": "terraform",
            "help": "Apply Aws terraform objects",
            "formula": {
              "path": "aws/terraform",
              "bin": "terraform-cli-${so}",
              "bundle": "${so}.zip",
              "repoUrl": "https://commons-repo.ritchiecli.io/formulas"
            },
            "parent": "root_aws_apply",
            "roles" : ["USER"]
          },
          {
            "usage": "scaffold",
            "help": "Manipulate scaffold objects",
            "parent": "root"
          },
          {
            "usage": "generate",
            "help": "Generates a scaffold by some template",
            "parent": "root_scaffold"
          },
          {
            "usage": "coffee-go",
            "help": "Generates a project by coffee template in Go",
            "formula": {
              "path": "scaffold/coffee-go",
              "bin": "coffee-go-${so}",
              "bundle": "${so}.zip",
              "repoUrl": "https://commons-repo.ritchiecli.io/formulas"
            },
            "parent": "root_scaffold_generate"
          }
        ],
        "version": "1.0.0"
      }`
	var s server.Tree
	if err := json.Unmarshal([]byte(js), &s); err != nil {
		log.Fatal("error Unmarshal tree")
	}
	return s
}

