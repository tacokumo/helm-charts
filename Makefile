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

.PHONY: submodule copy-crds

submodule:
	git submodule init
	git submodule update --recursive

copy-crds:
	@mkdir -p charts/portal-controller-kubernetes/templates/crds
	cp portal-controller-kubernetes/config/crd/bases/tacokumo.github.io_applications.yaml \
		charts/portal-controller-kubernetes/templates/crds/crd-application.yaml
	cp portal-controller-kubernetes/config/crd/bases/tacokumo.github.io_portals.yaml \
		charts/portal-controller-kubernetes/templates/crds/crd-portal.yaml
	cp portal-controller-kubernetes/config/crd/bases/tacokumo.github.io_releases.yaml \
		charts/portal-controller-kubernetes/templates/crds/crd-release.yaml
