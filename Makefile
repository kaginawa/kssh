.DEFAULT_GOAL := help

.PHONY: setup
setup: ## Resolve dependencies using Go Modules
	go mod download

.PHONY: clean
clean: ## Remove build artifact directory
	-rm -rfv build

.PHONY: test
test: ## Tests all code
	go test -cover -race ./...

.PHONY: lint
lint: ## Runs static code analysis
	command -v golint >/dev/null 2>&1 || { go get golang.org/x/lint/golint; }
	golint -set_exit_status ./...

.PHONY: run
run: ## Run agent without build artifact generation
	go run . -d

.PHONY: build
build: ## Build executable binaries for local execution
	go build -ldflags "-s -w" -o build/kssh .

.PHONY: build-all
build-all: build ## Build executable binaries for all supported OSs and architectures
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.ver=`git describe --tags`" -o build/kssh.exe .
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.ver=`git describe --tags`" -o build/kssh.macos-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X main.ver=`git describe --tags`" -o build/kssh.macos-arm64 .
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.ver=`git describe --tags`" -o build/kssh.linux-x64 .
	GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "-s -w -X main.ver=`git describe --tags`" -o build/kssh.linux-arm6 .
	GOOS=linux GOARCH=arm GOARM=7 go build -ldflags "-s -w -X main.ver=`git describe --tags`" -o build/kssh.linux-arm7 .
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.ver=`git describe --tags`" -o build/kssh.linux-arm8 .

.PHONY: count-go
count-go: ## Count number of lines of all go codes
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
