#!/bin/sh
cd /home/application
ls
#./wait-for-it.sh "stubby4j:8882" && echo "stubby4j is up"
#./wait-for-it.sh "vault:8200" && echo "vault is up"
#./wait-for-it.sh "keycloak:8080" && echo "keycloak is up"
#
#./create-vault-approle.sh . http://vault:8200
#
#export VAULT_ADDR=http://vault:8200
#export VAULT_AUTHENTICATION=APPROLE
#export VAULT_ROLE_ID=$(cat /tmp/vault/role-id.txt)
#export VAULT_SECRET_ID=$(cat /tmp/vault/secret-id.txt)
#export FILE_CONFIG="$(pwd)/server/resources/file_config_ci.json"
#export KEYCLOAK_URL=http://keycloak:8080
#export OAUTH_URL=http://keycloak:8080/auth/realms/ritchie
#export CLI_VERSION_URL=http://stubby4j:8882/s3-version-mock
#
#mkdir -p bin
#go test -v -coverprofile=bin/cov.out `go list ./... | grep -v vendor/`
#testStatus=$?
#if [ $testStatus -ne 0 ]; then
#    echo "Tests failed"
#    exit 1
#fi
#
#go tool cover -func=bin/cov.out
#rm -rf testdata/file_config_test.json

