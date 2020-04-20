NAME := jcli
CGO_ENABLED = 0
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
# CHANGE_LOG := $(shell echo -n "$(shell hub release show $(shell hub release --include-drafts -L 1))" | base64)
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BUILDFLAGS = -ldflags "-X github.com/jenkins-zh/jenkins-cli/app.version=$(VERSION) -X github.com/jenkins-zh/jenkins-cli/app.commit=$(COMMIT)"
COVERED_MAIN_SRC_FILE=./main
PATH  := $(PATH):$(PWD)/bin

gen-mock:
	go get github.com/golang/mock/gomock
	go install github.com/golang/mock/mockgen
	mockgen -destination ./mock/mhttp/roundtripper.go -package mhttp net/http RoundTripper

init: gen-mock

darwin:
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)
	rm -rf jcli && ln -s bin/darwin/$(NAME) jcli

linux:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)

win:
	go get github.com/inconshreveable/mousetrap
	go get github.com/mattn/go-isatty
	CGO_ENABLED=$(CGO_ENABLED) GOOS=windows GOARCH=386 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/windows/$(NAME).exe $(MAIN_SRC_FILE)

build-all: darwin linux win

release: build-all
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

gen-data-linux: go-bindata-download-linux
	cd app/i18n && ../../bin/go-bindata -o bindata.go -pkg i18n jcli/zh_CN/LC_MESSAGES/

go-bindata-download-darwin:
	mkdir -p bin
	curl -L https://github.com/kevinburke/go-bindata/releases/download/v3.11.0/go-bindata-darwin-amd64 -o bin/go-bindata
	chmod u+x bin/go-bindata

gen-data-darwin: go-bindata-download-darwin
	cd app/i18n && ../../bin/go-bindata -o bindata.go -pkg i18n jcli/zh_CN/LC_MESSAGES/

verify: dep tools lint


lint:
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
	gofmt -s -w .

test:
	mkdir -p bin
	go test ./util -v -count=1
	go test ./client -v -count=1 -coverprofile coverage.out
	go test ./app -v -count=1
	go test ./app/health -v -count=1
	go test ./app/helper -v -count=1
	go test ./app/i18n -v -count=1
	go test ./app/cmd -v -count=1

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
	docker build . -t jenkinszh/jcli

image-win:
	docker build . -t jenkinszh/jcli:win -f Dockerfile-win

image-darwin:
	docker build . -t jenkinszh/jcli:darwin -f Dockerfile-darwin

image-dev:
	docker build . -t jenkinszh/jcli:dev -f Docker-dev
