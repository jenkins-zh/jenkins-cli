NAME := jcli
CGO_ENABLED = 0
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
# CHANGE_LOG := $(shell echo -n "$(shell hub release show $(shell hub release --include-drafts -L 1))" | base64)
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BUILDFLAGS = -ldflags "-X github.com/jenkins-zh/jenkins-cli/app.version=$(VERSION) -X github.com/jenkins-zh/jenkins-cli/app.commit=$(COMMIT)"
COVERED_MAIN_SRC_FILE=./main

gen-mock:
	go get github.com/golang/mock/gomock@v1.4.0
	go get github.com/golang/mock/mockgen@v1.4.0
	mockgen -destination ./mock/mhttp/roundtripper.go -package mhttp net/http RoundTripper

init: gen-mock

darwin: gen-data
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)

linux: gen-data
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)

win: gen-data
	go get github.com/inconshreveable/mousetrap
	go get github.com/mattn/go-isatty
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=386 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/windows/$(NAME).exe $(MAIN_SRC_FILE)

build-all: darwin linux win

release: clean build-all
	mkdir release
	cd ./bin/darwin; upx jcli; tar -zcvf ../../release/jcli-darwin-amd64.tar.gz jcli; cd ../../release/; shasum -a 256 jcli-darwin-amd64.tar.gz > jcli-darwin-amd64.txt
	cd ./bin/linux; upx jcli; tar -zcvf ../../release/jcli-linux-amd64.tar.gz jcli; cd ../../release/; shasum -a 256 jcli-linux-amd64.tar.gz > jcli-linux-amd64.txt
	cd ./bin/windows; upx jcli.exe; tar -zcvf ../../release/jcli-windows-386.tar.gz jcli.exe; cd ../../release/; shasum -a 256 jcli-windows-386.tar.gz > jcli-windows-386.txt

clean: ## Clean the generated artifacts
	rm -rf bin release
	rm -rf coverage.out
	rm -rf app/cmd/test-app.xml
	rm -rf app/test-app.xml
	rm -rf util/test-utils.xml

copy: darwin
	sudo cp bin/darwin/$(NAME) $(shell which jcli)

copy-linux: linux
	cp bin/linux/$(NAME) /usr/local/bin/jcli

tools: i18n-tools
	go get -u golang.org/x/lint/golint

i18n-tools:
	go get -u github.com/gosexy/gettext/go-xgettext
# 	go get -u github.com/go-bindata/go-bindata/...
# 	go get -u github.com/kevinburke/go-bindata/...

go-bindata-download-linux:
	mkdir -p bin
	curl -L https://github.com/kevinburke/go-bindata/releases/download/v3.11.0/go-bindata-linux-amd64 -o bin/go-bindata
	chmod u+x bin/go-bindata

go-bindata-download-darwin:
	mkdir -p bin
	curl -L https://github.com/kevinburke/go-bindata/releases/download/v3.11.0/go-bindata-darwin-amd64 -o bin/go-bindata
	chmod u+x bin/go-bindata

go-bindata-download-windows:
	mkdir -p bin
	curl -L https://github.com/kevinburke/go-bindata/releases/download/v3.11.0/go-bindata-windows-amd64 -o bin/go-bindata
	chmod u+x bin/go-bindata

verify:
	go vet ./...
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

test: clean gen-data verify fmt
	mkdir -p bin
	go vet ./...
	go test ./... -v -coverprofile coverage.out

test-slow-latest:
	JENKINS_VERSION=2.190.3 make test-slow

test-slow:
	go test ./... -v -coverprofile coverage.out

dep:
	go get github.com/AlecAivazis/survey/v2
	go get github.com/gosuri/uiprogress
	go get github.com/spf13/cobra
	go get github.com/spf13/viper
	go get gopkg.in/yaml.v2
	go get github.com/Pallinder/go-randomdata

JCLI_FILES="app/cmd/*.go"
gettext:
	go-xgettext -k=i18n.T "${JCLI_FILES}" > app/i18n/jcli.pot

gen-data:
	cd app/i18n && go-bindata -o bindata.go -pkg i18n jcli/zh_CN/LC_MESSAGES/

image:
	docker build . -t jenkinszh/jcli

image-win:
	docker build . -t jenkinszh/jcli:win -f Dockerfile-win

image-darwin:
	docker build . -t jenkinszh/jcli:darwin -f Dockerfile-darwin

image-dev:
	docker build . -t jenkinszh/jcli:dev -f Docker-dev