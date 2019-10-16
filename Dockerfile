ARG KUBEBUILDER_VERSION=2.0.0

FROM golang:1.13-alpine as base
RUN apk add --no-cache make curl git

FROM base AS modules
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/talos-systems/cluster-api-bootstrap-provider-talos
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN go mod verify
COPY ./api ./api
COPY ./config ./config
COPY ./controllers ./controllers
COPY ./hack ./hack
COPY main.go main.go
COPY Makefile Makefile
COPY PROJECT PROJECT
RUN go mod vendor
RUN go list -mod=readonly all >/dev/null
RUN ! go mod tidy -v 2>&1 | grep .

FROM modules AS test
RUN mkdir -p /usr/local/kubebuilder/bin
ARG KUBEBUILDER_VERSION
RUN curl -L https://github.com/kubernetes-sigs/kubebuilder/releases/download/v${KUBEBUILDER_VERSION}/kubebuilder_${KUBEBUILDER_VERSION}_linux_amd64.tar.gz | tar -xvz --strip-components=2 -C /usr/local/kubebuilder/bin
RUN make generate fmt vet manifests && go test ./... -coverprofile cover.out

# Build the manager binary
FROM test AS build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

# Copy the controller-manager into a thin image
FROM scratch
WORKDIR /
COPY --from=build /go/src/github.com/talos-systems/cluster-api-bootstrap-provider-talos/manager .
ENTRYPOINT ["/manager"]
