# syntax = docker/dockerfile-upstream:1.1.4-experimental

FROM golang:1.13 AS build
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED 0
WORKDIR /tmp
RUN go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN go mod verify
COPY ./ ./
RUN go list -mod=readonly all >/dev/null
RUN ! go mod tidy -v 2>&1 | grep .

FROM build AS manifests-build
ARG NAME
RUN controller-gen rbac:roleName=manager-role crd paths="./..." output:rbac:artifacts:config=config/rbac output:crd:artifacts:config=config/crd/bases
FROM scratch AS manifests
COPY --from=manifests-build /src/config/crd /config/crd
COPY --from=manifests-build /src/config/rbac /config/rbac

FROM build AS generate-build
RUN controller-gen object:headerFile=./hack/boilerplate.go.txt paths="./..."
FROM scratch AS generate
COPY --from=generate-build /src/api /api

FROM k8s.gcr.io/hyperkube:v1.17.0 AS release-build
RUN apt update -y \
  && apt install -y curl \
  && curl -LO https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.4.0/kustomize_v3.4.0_linux_amd64.tar.gz \
  && tar -xf kustomize_v3.4.0_linux_amd64.tar.gz -C /usr/local/bin \
  && rm kustomize_v3.4.0_linux_amd64.tar.gz
COPY ./config ./config
ARG REGISTRY_AND_USERNAME
ARG NAME
ARG TAG
RUN cd config/manager \
  && kustomize edit set image controller=${REGISTRY_AND_USERNAME}/${NAME}:${TAG} \
  && cd - \
  && kubectl kustomize config/default >/bootstrap-components.yaml \
  && cp config/metadata/metadata.yaml /metadata.yaml
FROM scratch AS release
COPY --from=release-build /bootstrap-components.yaml /bootstrap-components.yaml
COPY --from=release-build /metadata.yaml /metadata.yaml

FROM build AS binary
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux go build -ldflags "-s -w" -o /manager
RUN chmod +x /manager

FROM scratch AS container
COPY --from=docker.io/autonomy/ca-certificates:v0.1.0 / /
COPY --from=docker.io/autonomy/fhs:v0.1.0 / /
COPY --from=binary /manager /manager
ENTRYPOINT [ "/manager" ]
