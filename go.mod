module github.com/talos-systems/cluster-api-bootstrap-provider-talos

go 1.13

replace github.com/kubernetes-sigs/bootkube => github.com/talos-systems/bootkube v0.14.1-0.20200131192519-720c01d02032

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/talos-systems/talos v0.4.0-alpha.5
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	k8s.io/utils v0.0.0-20200229041039-0a110f9eb7ab
	sigs.k8s.io/cluster-api v0.3.3
	sigs.k8s.io/controller-runtime v0.5.2
)
