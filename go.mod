module github.com/talos-systems/cluster-api-bootstrap-provider-talos

go 1.13

replace github.com/kubernetes-sigs/bootkube => github.com/talos-systems/bootkube v0.14.1-0.20200131192519-720c01d02032

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/talos-systems/crypto v0.2.0
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20201006184949-3961f835f502
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v0.18.2
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89
	sigs.k8s.io/cluster-api v0.3.6
	sigs.k8s.io/controller-runtime v0.6.0
)
