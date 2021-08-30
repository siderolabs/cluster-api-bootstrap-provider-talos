# syntax = docker/dockerfile-upstream:1.2.0-labs

# Meta args applied to stage base names.

ARG TOOLS
ARG PKGS

# Resolve package images using ${PKGS} to be used later in COPY --from=.

FROM ghcr.io/talos-systems/ca-certificates:${PKGS} AS pkg-ca-certificates
FROM ghcr.io/talos-systems/fhs:${PKGS} AS pkg-fhs

# The base target provides the base for running various tasks against the source
# code

FROM --platform=${BUILDPLATFORM} ${TOOLS} AS build
SHELL ["/toolchain/bin/bash", "-c"]
ENV PATH /toolchain/bin:/toolchain/go/bin:/go/bin
RUN ["/toolchain/bin/mkdir", "/bin", "/tmp"]
RUN ["/toolchain/bin/ln", "-svf", "/toolchain/bin/bash", "/bin/sh"]
RUN ["/toolchain/bin/ln", "-svf", "/toolchain/etc/ssl", "/etc/ssl"]
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED 0
ENV GOCACHE /.cache/go-build
ENV GOMODCACHE /.cache/mod
ARG CONTROLLER_GEN_VERSION
ARG CONVERSION_GEN_VERSION
RUN --mount=type=cache,target=/.cache go install sigs.k8s.io/controller-tools/cmd/controller-gen@${CONTROLLER_GEN_VERSION}
RUN --mount=type=cache,target=/.cache go install k8s.io/code-generator/cmd/conversion-gen@${CONVERSION_GEN_VERSION}
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
RUN --mount=type=cache,target=/.cache go mod download
RUN --mount=type=cache,target=/.cache go mod verify
COPY ./ ./
RUN --mount=type=cache,target=/.cache go list -mod=readonly all >/dev/null
RUN --mount=type=cache,target=/.cache ! go mod tidy -v 2>&1 | grep .

FROM build AS manifests-build
ARG NAME
RUN --mount=type=cache,target=/.cache controller-gen crd:crdVersions=v1 paths="./api/..." output:crd:dir=config/crd/bases output:webhook:dir=config/webhook webhook
RUN --mount=type=cache,target=/.cache controller-gen rbac:roleName=manager-role paths="./controllers/..." output:rbac:dir=config/rbac

FROM scratch AS manifests
COPY --from=manifests-build /src/config /config

FROM build AS generate-build
RUN --mount=type=cache,target=/.cache controller-gen object:headerFile=./hack/boilerplate.go.txt paths="./..."
RUN --mount=type=cache,target=/.cache conversion-gen --input-dirs=./api/v1alpha2 --output-base ./ --output-file-base=zz_generated.conversion --go-header-file=./hack/boilerplate.go.txt

FROM scratch AS generate
COPY --from=generate-build /src/api /api

FROM build AS integration-test-build
ENV CGO_ENABLED 1
ARG TALOS_VERSION
ARG GO_LDFLAGS="-linkmode=external -extldflags '-static' -X github.com/talos-systems/cluster-api-bootstrap-provider-talos/internal/integration.TalosVersion=${TALOS_VERSION}"
RUN --mount=type=cache,target=/.cache go test -race -ldflags "${GO_LDFLAGS}" -coverpkg=./... -v -c ./internal/integration

FROM scratch AS integration-test
COPY --from=integration-test-build /src/integration.test /integration.test

FROM --platform=${BUILDPLATFORM} alpine:3.13 AS release-build
ADD https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv4.1.0/kustomize_v4.1.0_linux_amd64.tar.gz .
RUN  tar -xf kustomize_v4.1.0_linux_amd64.tar.gz -C /usr/local/bin && rm kustomize_v4.1.0_linux_amd64.tar.gz
COPY ./config ./config
ARG REGISTRY_AND_USERNAME
ARG NAME
ARG TAG
RUN cd config/manager \
  && kustomize edit set image controller=${REGISTRY_AND_USERNAME}/${NAME}:${TAG} \
  && cd - \
  && kustomize build config/default > /bootstrap-components.yaml \
  && cp config/metadata/metadata.yaml /metadata.yaml

FROM scratch AS release
ARG TAG
COPY --from=release-build /bootstrap-components.yaml /bootstrap-talos/${TAG}/bootstrap-components.yaml
COPY --from=release-build /metadata.yaml /bootstrap-talos/${TAG}/metadata.yaml

FROM build AS binary
ARG TARGETARCH
RUN --mount=type=cache,target=/.cache GOOS=linux GOARCH=${TARGETARCH} go build -ldflags "-s -w" -o /manager
RUN chmod +x /manager

FROM scratch AS container
COPY --from=pkg-ca-certificates / /
COPY --from=pkg-fhs / /
COPY --from=binary /manager /manager
LABEL org.opencontainers.image.source https://github.com/talos-systems/cluster-api-bootstrap-provider-talos
ENTRYPOINT [ "/manager" ]
