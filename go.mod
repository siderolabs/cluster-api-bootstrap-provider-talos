module github.com/talos-systems/cluster-api-bootstrap-provider-talos

go 1.16

require (
	github.com/evanphx/json-patch v4.11.0+incompatible
	github.com/go-logr/logr v0.1.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/talos-systems/crypto v0.3.1
	github.com/talos-systems/talos/pkg/machinery v0.11.3
	golang.org/x/sys v0.0.0-20210816074244-15123e1e1f71
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.17.9
	k8s.io/apiextensions-apiserver v0.17.9
	k8s.io/apimachinery v0.17.9
	k8s.io/client-go v0.17.9
	k8s.io/utils v0.0.0-20200619165400-6e3d28b6ed19
	sigs.k8s.io/cluster-api v0.3.22
	sigs.k8s.io/controller-runtime v0.5.14
)
