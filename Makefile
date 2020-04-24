REGISTRY = $(DOCKER_REGISTRY)
RELEASE = $(RELEASE_VERSION)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=ritchie-server
CMD_PATH=./server/cmd/server/main.go

GIT_REMOTE=https://$(GIT_USERNAME):$(GIT_PASSWORD)@github.com/ZupIT/ritchie-server

# Docker
DOCKERCMD=docker
DOCKERBUILD=${DOCKERCMD} build
DOCKERPUSH=${DOCKERCMD} push
DOCKERTAG=${DOCKERCMD} tag

all: test build
build:
	GOOS=linux GOARCH=amd64 ${GOBUILD} -o ./${BINARY_NAME} -v ${CMD_PATH}
	cp $(BINARY_NAME) server
	$(DOCKERBUILD) -t "${REGISTRY}/${BINARY_NAME}:${RELEASE}" ./server
	# $(DOCKERTAG) "${REGISTRY}/${BINARY_NAME}:${RELEASE}" "${REGISTRY}/${BINARY_NAME}:latest"
	rm go.sum ritchie-server

build-local-mac:
	GOOS=darwin GOARCH=amd64 ${GOBUILD} -o ./${BINARY_NAME} -v ${CMD_PATH}

build-local:
	${GOBUILD} -o ./${BINARY_NAME} -v ${CMD_PATH}

publish:
	LOGIN_CMD="$(shell aws ecr get-login --region ${DOCKER_AWS_REGION} --no-include-email | sed 's|https://||')"
	${LOGIN_CMD}
	${DOCKERPUSH} "${REGISTRY}/${BINARY_NAME}:${RELEASE}"
	# ${DOCKERPUSH} "${REGISTRY}/${BINARY_NAME}:latest"

test:
	DOCKER_REGISTRY_BUILDER=${REGISTRY} docker-compose -f docker-compose-ci.yml run server

test-local:
	docker-compose up -d
	./run-tests.sh
	docker-compose down

release:
	git config --global user.email "$(GIT_EMAIL)"
	git config --global user.name "$(GIT_USER)"
	git add .
	git commit --allow-empty -m "release"
	git push $(GIT_REMOTE) HEAD:release-$(RELEASE_VERSION)
	git tag -a $(RELEASE_VERSION) -m "release"
	git push $(GIT_REMOTE) $(RELEASE_VERSION)
	curl --user $(GIT_USERNAME):$(GIT_PASSWORD) -X POST https://api.github.com/repos/ZupIT/ritchie-server/pulls -H 'Content-Type: application/json' -d '{ "title": "Release $(RELEASE_VERSION) merge", "body": "Release $(RELEASE_VERSION) merge with master", "head": "release-$(RELEASE_VERSION)", "base": "master" }'

build-circle:
	mkdir bin
	GOOS=linux GOARCH=amd64 ${GOBUILD} -o ./bin/${BINARY_NAME} -v ${CMD_PATH}

build-container-circle:
	cp bin/$(BINARY_NAME) server
	$(DOCKERBUILD) -t "${REGISTRY}/${BINARY_NAME}:${RELEASE}" ./server
