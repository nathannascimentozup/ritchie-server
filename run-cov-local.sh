#!/bin/sh

export VAULT_ADDR=http://localhost:8200
export VAULT_AUTHENTICATION=APPROLE
export VAULT_ROLE_ID=$(cat /tmp/vault/role-id.txt)
export VAULT_SECRET_ID=$(cat /tmp/vault/secret-id.txt)
export FILE_CONFIG="$(pwd)/server/resources/file_config_local.json"
export KEYCLOAK_URL=http://localhost:8080

mkdir -p bin
go test -v -coverprofile=bin/cov.out $1
rm testdata/file_config_test.json
go tool cover -html=bin/cov.out
