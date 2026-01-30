.PHONY: all
all: format test lint

.PHONY: format
format:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

GOLANGCI_LINT_VERSION := v2.8.0
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b tools $(GOLANGCI_LINT_VERSION))
	./tools/golangci-lint run
