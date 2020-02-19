REGISTRY ?= docker.io
USERNAME ?= autonomy
SHA ?= $(shell git describe --match=none --always --abbrev=8 --dirty)
TAG ?= $(shell git describe --tag --always --dirty)
BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)
NAME := cluster-api-talos-controller

ARTIFACTS := _out

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
  app: cluster-api-talos-controller
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
	@$(MAKE) local-$@ DEST=./

.PHONY: container
container: generate ## Build the container image.
	@$(MAKE) docker-$@ TARGET_ARGS="--push=$(PUSH)"
	sed -i'' -e 's@image: .*@image: '"$(REGISTRY_AND_USERNAME)/$(NAME):$(TAG)"'@' ./config/default/manager_image_patch.yaml


.PHONY: manifests
manifests: ## Generate manifests (e.g. CRD, RBAC, etc.).
	@$(MAKE) local-$@ DEST=./

.PHONY: release
release: manifests container ## Create the release YAML. The build result will be ouput to the specified local destination.
	@$(MAKE) local-$@ DEST=./$(ARTIFACTS)

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
