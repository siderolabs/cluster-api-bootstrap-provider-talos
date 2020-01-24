module github.com/talos-systems/cluster-api-bootstrap-provider-talos

go 1.12

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.1
	github.com/talos-systems/talos v0.3.1
	gopkg.in/yaml.v2 v2.2.7
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	sigs.k8s.io/cluster-api v0.2.9
	sigs.k8s.io/controller-runtime v0.4.0
)
