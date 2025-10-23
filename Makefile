LOCAL_BIN := $(CURDIR)/bin
GOIMPORTS_BIN := $(LOCAL_BIN)/goimports
GOLANGCI_BIN := $(LOCAL_BIN)/golangci-lint
GO_TEST=$(LOCAL_BIN)/gotest
GO_TEST_ARGS=-race -v -tags=integration_test ./...
UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)

ARCH :=

.PHONY: lint
lint:
	echo 'Running linter on files...'
	$(GOLANGCI_BIN) run \
	--config=.golangci.yaml \
	--sort-results \
	--max-issues-per-linter=0 \
	--max-same-issues=0

build:
	go mod tidy
	go build -o ./bin/habit-bot ./cmd/habit-bot/