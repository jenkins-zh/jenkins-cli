NAME := jcli
CGO_ENABLED = 0
GO := go
BUILD_TARGET = build
COMMIT := $(shell git rev-parse --short HEAD)
VERSION := dev-$(shell git describe --tags $(shell git rev-list --tags --max-count=1))
BUILDFLAGS = -ldflags "-X github.com/jenkins-zh/jenkins-cli/app.version=$(VERSION) -X github.com/jenkins-zh/jenkins-cli/app.commit=$(COMMIT)"
COVERED_MAIN_SRC_FILE=./main

darwin: ## Build for OSX
	GO111MODULE=on CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)

linux: ## Build for linux
	CGO_ENABLED=$(CGO_ENABLED) GOOS=linux GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/linux/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/linux/$(NAME)

win: ## Build for windows
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

dep:
	go get github.com/AlecAivazis/survey/v2
	go get github.com/gosuri/uiprogress
	go get github.com/spf13/cobra
	go get github.com/spf13/viper
	go get gopkg.in/yaml.v2
	go get github.com/Pallinder/go-randomdata
