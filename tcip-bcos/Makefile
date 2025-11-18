SDK_VERSION=v2.3.1
VERSION=v2.3.1
DATETIME=$(shell date "+%Y%m%d%H%M%S")
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT = $(shell git log --pretty=format:'%h' -n 1)

LOCALCONF_HOME=chainmaker.org/chainmaker/tcip-bcos/v2/module/conf
GOLDFLAGS += -X "${LOCALCONF_HOME}.CurrentVersion=${VERSION}"
GOLDFLAGS += -X "${LOCALCONF_HOME}.BuildTime=${DATETIME}"
GOLDFLAGS += -X "${LOCALCONF_HOME}.CurrentBranch=${GIT_BRANCH}"
GOLDFLAGS += -X "${LOCALCONF_HOME}.CurrentCommit=${GIT_COMMIT}"

build :
	cd main && go build -ldflags '${GOLDFLAGS}' -o ../tcip-bcos ./
ut:
	sudo rm -rf $GOPATH/pkg/mod/chainmaker.org/chainmaker/tcip-chainmaker
	sudo rm -rf $GOPATH/pkg/mod/cache/download/chainmaker.org/chainmaker/tcip-chainmaker
	go test ./...

release:
	./build_release.sh

docker_build:
	./build_release.sh
	docker build -t tcip-bcos -f ./DOCKER/Dockerfile .
	docker tag tcip-bcos tcip-bcos:${VERSION}

lint:
	golangci-lint run  ./...

gomod:
	go get chainmaker.org/chainmaker/sdk-go/v2@$(SDK_VERSION)
	go get chainmaker.org/chainmaker/tcip-go/v2@$(VERSION)
	go mod tidy