REGISTRY ?= ghcr.io
USERNAME ?= siderolabs
SHA ?= $(shell git describe --match=none --always --abbrev=8 --dirty)
TAG ?= $(shell git describe --tag --always --dirty)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)
NAME := cluster-api-talos-controller

ARTIFACTS := _out
TEST_RUN ?= ./...

TOOLS ?= ghcr.io/siderolabs/tools:v1.4.0-1-g955aabc
PKGS ?= v1.4.1-5-ga333a84
TALOS_VERSION ?= v1.4.0
K8S_VERSION ?= 1.27.1

CONTROLLER_GEN_VERSION ?= v0.11.3
CONVERSION_GEN_VERSION ?= v0.26.0

BUILD := docker buildx build
PLATFORM ?= linux/amd64
PROGRESS ?= auto
PUSH ?= false
COMMON_ARGS := --file=Dockerfile
COMMON_ARGS += --progress=$(PROGRESS)
COMMON_ARGS += --platform=$(PLATFORM)
COMMON_ARGS += --build-arg=REGISTRY_AND_USERNAME=$(REGISTRY_AND_USERNAME)
COMMON_ARGS += --build-arg=NAME=$(NAME)
COMMON_ARGS += --build-arg=TAG=$(TAG)
COMMON_ARGS += --build-arg=PKGS=$(PKGS)
COMMON_ARGS += --build-arg=TOOLS=$(TOOLS)
COMMON_ARGS += --build-arg=CONTROLLER_GEN_VERSION=$(CONTROLLER_GEN_VERSION)
COMMON_ARGS += --build-arg=CONVERSION_GEN_VERSION=$(CONVERSION_GEN_VERSION)
COMMON_ARGS += --build-arg=TALOS_VERSION=$(TALOS_VERSION)
COMMON_ARGS += --build-arg=ARTIFACTS=$(ARTIFACTS)

all: manifests container

.PHONY: help
help: ## This help menu.
	@grep -E '^[a-zA-Z%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

target-%: ## Builds the specified target defined in the Dockerfile. The build result will remain only in the build cache.
	@$(BUILD) \
		--target=$* \
		$(COMMON_ARGS) \
		$(TARGET_ARGS) .

local-%: ## Builds the specified target defined in the Dockerfile using the local output type. The build result will be output to the specified local destination.
	@$(MAKE) target-$* TARGET_ARGS="--output=type=local,dest=$(DEST) $(TARGET_ARGS)"

docker-%: ## Builds the specified target defined in the Dockerfile using the docker output type. The build result will be loaded into docker.
	@$(MAKE) target-$* TARGET_ARGS="--tag $(REGISTRY_AND_USERNAME)/$(NAME):$(TAG) $(TARGET_ARGS)"

define RELEASEYAML
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: $(NAMESPACE)
commonLabels:
  app: $(NAME)
bases:
  - crd
  - rbac
  - manager
endef

export RELEASEYAML
.PHONY: init
init: ## Initialize the project.
	@mkdir tmp \
	&& cd tmp \
	&& kubebuilder init --repo $(REGISTRY_AND_USERNAME)/$(NAME) --domain $(DOMAIN) \
	&& rm -rf Dockerfile Makefile .gitignore bin hack \
	&& mv ./* ../ \
	&& cd .. \
	&& rm -rf tmp \
	&& echo "$$RELEASEYAML" > ./config/kustomization.yaml

.PHONY: generate
generate: ## Generate source code.
	@$(MAKE) local-$@ DEST=./ PLATFORM=linux/amd64

.PHONY: container
container: generate ## Build the container image.
	@$(MAKE) docker-$@ TARGET_ARGS="--push=$(PUSH)"

.PHONY: manifests
manifests: ## Generate manifests (e.g. CRD, RBAC, etc.).
	@$(MAKE) local-$@ DEST=./ PLATFORM=linux/amd64

.PHONY: release-notes
release-notes: ## Create the release notes.
	@mkdir -p $(ARTIFACTS)
	ARTIFACTS=$(ARTIFACTS) ./hack/release.sh $@ $(ARTIFACTS)/RELEASE_NOTES.md $(TAG)

.PHONY: release-manifests
release-manifests:
	@$(MAKE) local-$@ DEST=./$(ARTIFACTS) PLATFORM=linux/amd64

.PHONY: release
release: manifests container release-notes release-manifests ## Create the release YAML. The build result will be ouput to the specified local destination.

.PHONY: deploy
deploy: manifests ## Deploy to a cluster. This is for testing purposes only.
	kubectl apply -k config/default

.PHONY: destroy
destroy: ## Remove from a cluster. This is for testing purposes only.
	kubectl delete -k config/default

.PHONY: install
install: manifests ## Install CRDs into a cluster.
	kubectl apply -k config/crd

.PHONY: uninstall
uninstall: manifests ## Uninstall CRDs from a cluster.
	kubectl delete -k config/crd

.PHONY: run
run: install ## Run the controller locally. This is for testing purposes only.
	@$(MAKE) docker-container TARGET_ARGS="--load"
	@docker run --rm -it --net host -v $(PWD):/src -v $(KUBECONFIG):/root/.kube/config -e KUBECONFIG=/root/.kube/config $(REGISTRY_AND_USERNAME)/$(NAME):$(TAG)

.PHONY: clean
clean:
	@rm -rf $(ARTIFACTS)

check-dirty: ## Verifies that source tree is not dirty
	@if test -n "`git status --porcelain`"; then echo "Source tree is dirty"; git status; exit 1 ; fi

conformance:  ## Performs policy checks against the commit and source code.
	docker run --rm -it -v $(PWD):/src -w /src ghcr.io/siderolabs/conform:latest enforce

# Make `make test` behave just like `go test` regarding relative paths.
test:  ## Run tests.
	@$(MAKE) local-integration-test DEST=./internal/integration PLATFORM=linux/amd64
	cd internal/integration && KUBECONFIG=../../kubeconfig ./integration.test -test.v -test.coverprofile=../../coverage.txt -test.run $(TEST_RUN)

coverage:  ## Upload coverage data to codecov.io.
	/usr/local/bin/codecov -f coverage.txt -X fix

talosctl:
	curl -Lo talosctl https://github.com/siderolabs/talos/releases/download/$(TALOS_VERSION)/talosctl-$(shell uname -s | tr "[:upper:]" "[:lower:]")-amd64
	chmod +x ./talosctl

env-up: talosctl  ## Start development environment.
	./talosctl cluster create \
		--talosconfig=talosconfig \
		--name=cabpt-env \
		--kubernetes-version=$(K8S_VERSION) \
		--mtu=1450 \
		--skip-kubeconfig \
		--crashdump
	./talosctl kubeconfig kubeconfig \
		--talosconfig=talosconfig \
		--nodes=10.5.0.2 \
		--force

env-down: talosctl ## Stop development environment.
	./talosctl cluster destroy \
		--talosconfig=talosconfig \
		--name=cabpt-env
	rm -f talosconfig kubeconfig
