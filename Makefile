help:  ## Display this help message
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m  %-30s\033[0m %s\n", $$1, $$2}'

###############################################################################
###                                Docker                                 	###
###############################################################################

ARCH := $(shell arch | sed 's/x86_64/amd64/')
FOUNDRY_VERSION := $(shell cat scripts/foundry_version.txt)

.PHONY: build-docker
build-docker: ensure-go-releaser ## Builds the docker images using local arch.
	@git config --local core.abbrev 7
	@sed "s/amd64/$(ARCH)/g" .goreleaser-snapshot.yaml > .goreleaser-local.yaml
	@FOUNDRY_VERSION=$(FOUNDRY_VERSION) goreleaser release -f .goreleaser-local.yaml --snapshot --clean --skip=archive
	@rm .goreleaser-local.yaml
	@scripts/halovisor/build.sh

###############################################################################
###                                Contracts                                 ###
###############################################################################

.PHONY: contracts-gen
contract-bindings: ## Generate golang contract bindings.
	make -C ./contracts bindings

###############################################################################
###                                Utils                                 	###
###############################################################################

.PHONY: setup
setup: install-go-tools install-pre-commit
	@git config --local core.abbrev 7
	@git config --local commit.gpgsign true

.PHONY: install-cli
install-cli: ## Install the omni cli to $GOPATH/bin/omni.
	@go install github.com/omni-network/omni/cli/cmd/omni || echo "❌go install failed"
	@which omni || echo '❌ `which omni` failed, fix go environment: "export PATH=$$PATH:$$(go env GOPATH)/bin" # Or see https://go.dev/doc/gopath_code'

.PHONY: ensure-go-releaser
ensure-go-releaser: ## Installs the go-releaser tool.
	@which goreleaser > /dev/null || echo "go-releaser not installed, see https://goreleaser.com/install/"

.PHONY: install-pre-commit
install-pre-commit: ## Installs the pre-commit tool as the git pre-commit hook for this repo.
	@which pre-commit > /dev/null || echo "pre-commit not installed, see https://pre-commit.com/#install"
	@pre-commit install --install-hooks

.PHONY: install-go-tools
install-go-tools: ## Installs the go-dev-tools, like buf.
	@go generate scripts/tools.go

.PHONY: install-foundry
install-foundry: ## Installs the foundry tool.
	@FORCE=true ./scripts/install_foundry.sh

.PHONY: lint
lint: ## Runs linters via pre-commit.
	@pre-commit run -v --all-files

.PHONY: bufgen
bufgen: ## Generates protobufs using buf generate.
	@./scripts/buf_generate.sh

.PHONY: fix-golden
fix-golden: ## Fixes golden test fixtures.
	@./scripts/fix_golden_tests.sh


.PHONY: staging-addrs
staging-addrs: ## Prints staging address json.
	@go run ./scripts/stagingaddrs

###############################################################################
###                                Testing                                 	###
###############################################################################

.PHONY: devnet-deploy
devnet-deploy: ## Deploys devnet (MANIFEST=devnet1 by default)
	@echo "Creating a docker-compose devnet in ./e2e/run/$(if $(MANIFEST),$(MANIFEST),devnet1)"
	@go run github.com/omni-network/omni/e2e -f e2e/manifests/$(if $(MANIFEST),$(MANIFEST),devnet1).toml deploy

.PHONY: devnet-clean
devnet-clean: ## Deletes devnet containers (MANIFEST=devnet1 by default)
	@echo "Stopping the devnet in ./e2e/run/$(if $(MANIFEST),$(MANIFEST),devnet1)"
	@go run github.com/omni-network/omni/e2e -f e2e/manifests/$(if $(MANIFEST),$(MANIFEST),devnet1).toml clean

.PHONY: e2e-ci
e2e-ci: ## Runs all e2e CI tests
	@go install github.com/omni-network/omni/e2e
	@cd e2e && ./run-multiple.sh manifests/devnet1.toml manifests/devnet2.toml manifests/fuzzyhead.toml manifests/ci.toml

.PHONY: e2e-run
e2e-run: ## Run specific e2e manifest (MANIFEST=single, MANIFEST=devnet1, etc). Note container remain running after the test.
	@if [ -z "$(MANIFEST)" ]; then echo "⚠️ Please specify a manifest: MANIFEST=devnet1 make e2e-run" && exit 1; fi
	@echo "Using MANIFEST=$(MANIFEST)"
	@go run github.com/omni-network/omni/e2e -f e2e/manifests/$(MANIFEST).toml

.PHONY: e2e-logs
e2e-logs: ## Print the docker logs of previously ran e2e manifest (devnet1, etc).
	@if [ -z "$(MANIFEST)" ]; then echo "⚠️  Please specify a manifest: MANIFEST=devnet1 make e2e-logs" && exit 1; fi
	@echo "Using MANIFEST=$(MANIFEST)"
	@go run github.com/omni-network/omni/e2e -f e2e/manifests/$(MANIFEST).toml logs

.PHONY: e2e-clean
e2e-clean: ## Deletes all running containers from previously ran e2e.
	@if [ -z "$(MANIFEST)" ]; then echo "⚠️  Please specify a manifest: MANIFEST=devnet1 make e2e-clean" && exit 1; fi
	@echo "Using MANIFEST=$(MANIFEST)"
	@go run github.com/omni-network/omni/e2e -f e2e/manifests/$(MANIFEST).toml clean

.PHONY: unittest-run
unittest-run:
	go test -timeout=5m -race ./...
