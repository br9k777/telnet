GOPATH?=$(go env GOPATH)

.PHONY: setup
setup: ## Install all the build and lint dependencies
# 	go get -u github.com/alecthomas/gometalinter
	[ -r ${GOPATH}/bin/golangci-lint ] && rm ${GOPATH}/bin/golangci-lint
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(GOPATH)/bin v1.21.0	
	go get -u golang.org/x/tools/cmd/cover
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/gorilla/mux
	go get -u github.com/spf13/viper
	go get -u go.uber.org/zap
	go get -u github.com/davecgh/go-spew/spew
# 	gometalinter --install --update
# 	@$(MAKE) dep

.PHONY: test
test: test1 test2 test3
test1: ## Run all the tests		
	cd pkg/storage/ && go test -count=1 -run TestInsertEvents
test2:
	cd pkg/storage/ && go test -count=1 -run TestDeleteEvents
test3:	
	cd pkg/storage/ && go test -count=1 -run TestEditEvents
test_cover: 
	echo 'mode: atomic' > /tmp/coverage.txt && cd pkg/storage/ && \
	go test -covermode=atomic -coverprofile=/tmp/coverage.txt -v -race -timeout=30s ./... && go tool cover -html=/tmp/coverage.txt

# .PHONY: cover
# cover: test ## Run all the tests and opens the coverage report
# 	go tool cover -html=coverage.txt

.PHONY: fmt
fmt: ## Run goimports on all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: lint
lint: ## Run all the linters
	golangci-lint run ./...

# .PHONY: lint
# lint: ## Run all the linters
# 	gometalinter --vendor --disable-all \
# 		--enable=deadcode \
# 		--enable=ineffassign \
# 		--enable=gosimple \
# 		--enable=staticcheck \
# 		--enable=gofmt \
# 		--enable=goimports \
# 		--enable=misspell \
# 		--enable=errcheck \
# 		--enable=vet \
# 		--enable=vetshadow \
# 		--deadline=10m \
# 		./...

.PHONY: build
build: ## Build a version
	go build -o /tmp/lesson_telnet ./cmd/
# 	go build -v ./...

.PHONY: clean
clean: ## Remove temporary files
	go clean

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
