module github.com/talos-systems/cluster-api-bootstrap-provider-talos

go 1.12

require (
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/talos-systems/talos v0.3.0-alpha.10.0.20191209205235-3a93e65b5480
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.5
	k8s.io/api v0.0.0-20191109101513-0171b7c15da1
	k8s.io/apimachinery v0.0.0-20191111054156-6eb29fdf75dc
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	sigs.k8s.io/cluster-api v0.2.5
	sigs.k8s.io/controller-runtime v0.3.0
)

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190704095032-f4ca3d3bdf1d
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
)
