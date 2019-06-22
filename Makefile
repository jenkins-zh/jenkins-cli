NAME := jcli
CGO_ENABLED = 0
GO := go
BUILD_TARGET = build
BUILDFLAGS = 
COVERED_MAIN_SRC_FILE=./main

darwin: ## Build for OSX
	CGO_ENABLED=$(CGO_ENABLED) GOOS=darwin GOARCH=amd64 $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/darwin/$(NAME) $(MAIN_SRC_FILE)
	chmod +x bin/darwin/$(NAME)

build: $(GO_DEPENDENCIES) ## Build jx binary for current OS
	CGO_ENABLED=$(CGO_ENABLED) $(GO) $(BUILD_TARGET) $(BUILDFLAGS) -o bin/$(NAME) $(MAIN_SRC_FILE)

release: clean darwin
	mkdir release
	cd ./bin/darwin; tar -zcvf ../../release/jcli-darwin-amd64.tar.gz jcli

clean: ## Clean the generated artifacts
	rm -rf bin release