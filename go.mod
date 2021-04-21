module github.com/talos-systems/cluster-api-bootstrap-provider-talos

go 1.16

require (
	github.com/evanphx/json-patch v4.9.0+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	github.com/spf13/pflag v1.0.5
	github.com/talos-systems/crypto v0.2.1-0.20210202170911-39584f1b6e54
	github.com/talos-systems/talos/pkg/machinery v0.0.0-20210216142802-8d7a36cc0cc2
	gopkg.in/yaml.v2 v2.3.0
	k8s.io/api v0.17.9
	k8s.io/apiextensions-apiserver v0.17.9
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v0.17.9
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19
	sigs.k8s.io/cluster-api v0.3.12
	sigs.k8s.io/controller-runtime v0.5.14
)
