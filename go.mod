module github.com/talos-systems/cluster-api-bootstrap-provider-talos

go 1.16

require (
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.16.2
	github.com/onsi/gomega v1.13.0
	github.com/spf13/pflag v1.0.5
	github.com/talos-systems/crypto v0.3.1
	github.com/talos-systems/talos/pkg/machinery v0.11.3
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.17.9
	k8s.io/apiextensions-apiserver v0.17.9
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v0.17.9
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19
	sigs.k8s.io/cluster-api v0.3.12
	sigs.k8s.io/controller-runtime v0.5.14
)
