REGISTRY ?= ghcr.io
USERNAME ?= talos-systems
SHA ?= $(shell git describe --match=none --always --abbrev=8 --dirty)
TAG ?= $(shell git describe --tag --always --dirty)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)
NAME := cluster-api-talos-controller

ARTIFACTS := _out

TOOLS ?= ghcr.io/talos-systems/tools:v0.5.0
PKGS ?= v0.5.0
TALOS_VERSION ?= v0.12.0
K8S_VERSION ?= 1.21.4

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

.PHONY: release
release: manifests container ## Create the release YAML. The build result will be ouput to the specified local destination.
	@$(MAKE) local-$@ DEST=./$(ARTIFACTS) PLATFORM=linux/amd64

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

# Make `make test` behave just like `go test` regarding relative paths.
test:
	@$(MAKE) local-integration-test DEST=./internal/integration PLATFORM=linux/amd64
	cd internal/integration && KUBECONFIG=../../kubeconfig ./integration.test -test.v

talosctl:
	curl -Lo talosctl https://github.com/talos-systems/talos/releases/download/$(TALOS_VERSION)/talosctl-$(shell uname -s | tr "[:upper:]" "[:lower:]")-amd64
	chmod +x ./talosctl

env-up: talosctl
	./talosctl cluster create \
		--name=cabpt-env \
		--kubernetes-version=$(K8S_VERSION) \
		--mtu=1450 \
		--memory=2048 \
		--cpus=2.0 \
		--crashdump
	./talosctl -n 10.5.0.2 kubeconfig -f kubeconfig

env-down: talosctl
	./talosctl cluster destroy \
		--name=cabpt-env
