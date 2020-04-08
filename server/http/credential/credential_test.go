package credential

import (
	"github.com/hashicorp/vault/api"
)

type vaultManagerMock struct {
	Error      error
	ReturnMap  map[string]interface{}
	ReturnList []interface{}
}

const (
	bearerTest = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJweWZ0SV9UVWRzZEdadWtKRzd5aXRWdWQxWF9aNXNGREU0ZFJXMmE1" +
		"cGtFIn0.eyJqdGkiOiJmYzI2MGMwMS01ZDg2LTQ3YjEtODRjNS1mMTA2YTBhYTEzMTkiLCJleHAiOjE1ODUyODkwNzYsIm5iZiI6MCwiaWF0Ij" +
		"oxNTg1MjUzMDc2LCJpc3MiOiJodHRwczovL3JpdGNoaWUta2V5Y2xvYWsuaXRpYXdzLmRldi9hdXRoL3JlYWxtcy9yaXRjaGllIiwiYXVkIjpb" +
		"InVzZXItbG9naW4iLCJhY2NvdW50Il0sInN1YiI6IjU0ZGMyZjRlLWQzNDItNDkwYy04ZTkxLWM2NTBiNGZmYWVkNSIsInR5cCI6IkJlYXJlci" +
		"IsImF6cCI6InVzZXItbG9naW4iLCJhdXRoX3RpbWUiOjAsInNlc3Npb25fc3RhdGUiOiI5Y2UzNWQ5Yy01MmIwLTQyZGEtODRhMi05YzRmMDUx" +
		"YTY3ZDAiLCJhY3IiOiIxIiwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiYWRtaW4iLCJ1bWFfYXV0aG9yaXphdG" +
		"lvbiIsInVzZXIiXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291" +
		"bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6ImVtYWlsIHByb2ZpbGUiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibmFtZSI6Ik" +
		"FkbWluIGFkbWluIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiYWRtaW4iLCJnaXZlbl9uYW1lIjoiQWRtaW4iLCJmYW1pbHlfbmFtZSI6ImFkbWlu" +
		"IiwiZW1haWwiOiJhZG1pbkB6dXAuY29tLmJyIn0.LtcR-achyJMhK7ur95ogR3XVvMio_oxK_lkH6gwRt6-SbatktZryMA1Inq4MS7_W88ISZC" +
		"ox_BeLhs--UXYBU_j6DUHVEL3ZdTQxTV8Nw7Q59no5NFG-QPbn6KUnwIw1J7YHx0R6UK1fa6b8tdnj-4BGYbJTssfLn9W5G94h8DIyHUbHVUYo" +
		"6PEorV2gchPUHo4SQPPFacxaG20-nb4YtM5zap019NJKvJmobCwW6C3TspQNj4HH8Ykk74pfWs_lBXQjR773FuTbFA1cFyYTIWhr-nmtl3oQ67" +
		"q0M58C77yghWLVUvVYXeSJEupsakUVWT2cKj-6-zp-FhIT1uBe-g"
)

func (v vaultManagerMock) Write(key string, data map[string]interface{}) error {
	return v.Error
}
func (v vaultManagerMock) Read(key string) (map[string]interface{}, error) {
	return v.ReturnMap, v.Error
}
func (v vaultManagerMock) List(key string) ([]interface{}, error) {
	return v.ReturnList, v.Error
}
func (v vaultManagerMock) Delete(key string) error {
	return v.Error
}

func (v vaultManagerMock) Start(*api.Client) {
}

