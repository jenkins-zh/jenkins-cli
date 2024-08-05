NAME := jcli
CGO_ENABLED = 0
BUILD_GOOS=$(shell go env GOOS)
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
BIN_PATH:=$(shell rm -rf jcli && which jcli)
# CHANGE_LOG := $(shell echo -n "$(shell hub release show $(shell hub release --include-drafts -L 1))" | base64)
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BUILDFLAGS = -ldflags "-X github.com/linuxsuren/cobra-extension/version.version=$(VERSION) \
	-X github.com/linuxsuren/cobra-extension/version.commit=$(COMMIT) \
	-X github.com/linuxsuren/cobra-extension/version.date=$(shell date +'%Y-%m-%d')"
MAIN_SRC_FILE = main.go
PATH := $(PATH):$(PWD)/bin

.PHONY: build

build: pre-build
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=$(BUILD_GOOS) GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/$(BUILD_GOOS)/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/$(BUILD_GOOS)/$(NAME)
	rm -rf $(NAME) && ln -s bin/$(BUILD_GOOS)/$(NAME) $(NAME)

darwin: pre-build
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)
	rm -rf $(NAME) && ln -s bin/darwin/$(NAME) $(NAME)

linux: pre-build
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)
	rm -rf $(NAME)
	ln -s bin/linux/$(NAME) $(NAME)

win: pre-build
	go get github.com/inconshreveable/mousetrap
	go get github.com/mattn/go-isatty
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=386 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/windows/$(NAME).exe $(MAIN_SRC_FILE)

build-all: darwin linux win

init: gen-mock
gen-mock:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
	mockgen -destination ./mock/mhttp/roundtripper.go -package mhttp net/http RoundTripper

release: build-all
	mkdir -p release
	cd ./bin/darwin; upx $(NAME); tar -zcvf ../../release/$(NAME)-darwin-amd64.tar.gz $(NAME); cd ../../release/; shasum -a 256 $(NAME)-darwin-amd64.tar.gz > $(NAME)-darwin-amd64.txt
	cd ./bin/linux; upx $(NAME); tar -zcvf ../../release/$(NAME)-linux-amd64.tar.gz $(NAME); cd ../../release/; shasum -a 256 $(NAME)-linux-amd64.tar.gz > $(NAME)-linux-amd64.txt
	cd ./bin/windows; upx $(NAME).exe; tar -zcvf ../../release/$(NAME)-windows-386.tar.gz $(NAME).exe; cd ../../release/; shasum -a 256 $(NAME)-windows-386.tar.gz > $(NAME)-windows-386.txt

clean: ## Clean the generated artifacts
	rm -rf bin release
	rm -rf coverage.out
	rm -rf app/cmd/test-app.xml
	rm -rf app/test-app.xml
	rm -rf util/test-utils.xml

copy: build
	sudo cp bin/$(BUILD_GOOS)/$(NAME) $(BIN_PATH)

get-golint:
	go get -u golang.org/x/lint/golint

tools: i18n-tools get-golint

i18n-tools:
	go get -u github.com/gosexy/gettext/go-xgettext
# 	go get -u github.com/go-bindata/go-bindata/...
# 	go get -u github.com/kevinburke/go-bindata/...

go-bindata-download-linux:
	mkdir -p bin
	curl -L https://github.com/kevinburke/go-bindata/releases/download/v3.11.0/go-bindata-linux-amd64 -o bin/go-bindata
	chmod u+x bin/go-bindata

gen-data-linux: go-bindata-download-linux
	cd app/i18n && ../../bin/go-bindata -o bindata.go -pkg i18n jcli/zh_CN/LC_MESSAGES/

go-bindata-download-darwin:
	mkdir -p bin
	curl -L https://github.com/kevinburke/go-bindata/releases/download/v3.11.0/go-bindata-darwin-amd64 -o bin/go-bindata
	chmod u+x bin/go-bindata

gen-data-darwin: go-bindata-download-darwin
	cd app/i18n && ../../bin/go-bindata -o bindata.go -pkg i18n jcli/zh_CN/LC_MESSAGES/

verify: dep tools lint

pre-build:
	export GO111MODULE=on
	export GOPROXY=https://goproxy.io
	go mod tidy

vet:
	go vet ./...

lint: vet
	golint -set_exit_status app/cmd/...
	golint -set_exit_status app/helper/...
	golint -set_exit_status app/i18n/i18n.go
	golint -set_exit_status app/.
	golint -set_exit_status client/...
	golint -set_exit_status util/...

fmt:
	go fmt ./util/...
	go fmt ./client/...
	go fmt ./app/...
	gofmt -s -w .

test-slow:
#	JENKINS_VERSION=2.190.3 go test ./e2e/... -v -count=1 -parallel 1
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestBashCompletion$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestZshCompletion$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestPowerShellCompletion$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestListComputers$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestConfigList$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestConfigGenerate$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestConfigList$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestShowCurrentConfig$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestCrumb$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestDoc$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestSearchPlugins$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestListPlugins$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestCheckUpdateCenter$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestInstallPlugin$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestDownloadPlugin$
	go test github.com/jenkins-zh/jenkins-cli/e2e -v -test.run ^TestListQueue$

test:
	mkdir -p bin
	go test ./util ./app/health ./app/i18n ./app/cmd/common -v -count=1 -coverprofile coverage.out
#	go test ./util -v -count=1
#	go test ./client -v -count=1 -coverprofile coverage.out
#	go test ./app -v -count=1
#	go test ./app/health -v -count=1
#	go test ./app/helper -v -count=1
#	go test ./app/i18n -v -count=1
#	go test ./app/cmd -v -count=1

test-release:
	goreleaser release --rm-dist --snapshot --skip-publish

dep:
	go get github.com/AlecAivazis/survey/v2
	go get github.com/spf13/cobra
	go get github.com/spf13/viper
	go get gopkg.in/yaml.v2
	go get github.com/Pallinder/go-randomdata
	go install github.com/gosuri/uiprogress

JCLI_FILES="app/cmd/*.go"
gettext:
	go-xgettext -k=i18n.T "${JCLI_FILES}" > app/i18n/jcli.pot

gen-data:
	cd app/i18n && go-bindata -o bindata.go -pkg i18n jcli/zh_CN/LC_MESSAGES/

image:
	docker build . -t jenkinszh/$(NAME)

setup-env-centos:
	yum install make golang -y
