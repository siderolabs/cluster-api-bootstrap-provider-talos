
<a name="v0.2.0"></a>
## [v0.2.0](https://github.com/talos-systems/talos/compare/v0.2.0-beta.0...v0.2.0) (2021-06-21)


<a name="v0.2.0-beta.0"></a>
## [v0.2.0-beta.0](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.13...v0.2.0-beta.0) (2021-06-08)

### Chore

* update machinery to latest stable

### Fix

* remove unused kube-rbac-proxy, protect metrics-addr

### Release

* **v0.2.0-beta.0:** prepare release


<a name="v0.2.0-alpha.12"></a>
## [v0.2.0-alpha.12](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.11...v0.2.0-alpha.12) (2021-05-14)

### Chore

* rework build, move to ghcr.io, build for arm64/amd64

### Fix

* back down resource requests
* ensure secrets are deleted when cluster is dropped

### Release

* **v0.2.0-alpha.12:** prepare release


<a name="v0.2.0-alpha.11"></a>
## [v0.2.0-alpha.11](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.10...v0.2.0-alpha.11) (2021-02-19)

### Feat

* support EXP_MACHINE_POOL flag

### Release

* **v0.2.0-alpha.11:** prepare release


<a name="v0.2.0-alpha.10"></a>
## [v0.2.0-alpha.10](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.9...v0.2.0-alpha.10) (2021-02-17)

### Feat

* support talosVersion in talosconfig CRD

### Release

* **v0.2.0-alpha.10:** prepare release


<a name="v0.2.0-alpha.9"></a>
## [v0.2.0-alpha.9](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.8...v0.2.0-alpha.9) (2021-02-04)

### Feat

* support machinepools in bootstrapper

### Fix

* ensure proper ordering of packet machine config handling

### Release

* **v0.2.0-alpha.9:** prepare release


<a name="v0.2.0-alpha.8"></a>
## [v0.2.0-alpha.8](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.7...v0.2.0-alpha.8) (2020-12-03)

### Feat

* support config patches at the bootstrap provider level

### Release

* **v0.2.0-alpha.8:** prepare release


<a name="v0.2.0-alpha.7"></a>
## [v0.2.0-alpha.7](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.6...v0.2.0-alpha.7) (2020-11-13)

### Fix

* update talos machinery pkg

### Release

* **v0.2.0-alpha.7:** prepare release


<a name="v0.2.0-alpha.6"></a>
## [v0.2.0-alpha.6](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.5...v0.2.0-alpha.6) (2020-10-20)

### Chore

* update talos machinery v0.7.0-alpha.7

### Release

* **v0.2.0-alpha.6:** prepare release


<a name="v0.2.0-alpha.5"></a>
## [v0.2.0-alpha.5](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.4...v0.2.0-alpha.5) (2020-10-08)

### Fix

* make sure secrets are cluster owned

### Release

* **v0.2.0-alpha.5:** prepare release


<a name="v0.2.0-alpha.4"></a>
## [v0.2.0-alpha.4](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.3...v0.2.0-alpha.4) (2020-10-06)

### Fix

* ensure we have a dns domain

### Release

* **v0.2.0-alpha.4:** prepare release


<a name="v0.2.0-alpha.3"></a>
## [v0.2.0-alpha.3](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.2...v0.2.0-alpha.3) (2020-09-11)

### Fix

* ensure version is not nil

### Release

* **v0.2.0-alpha.3:** prepare release


<a name="v0.2.0-alpha.2"></a>
## [v0.2.0-alpha.2](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.1...v0.2.0-alpha.2) (2020-08-19)

### Fix

* change k8s version if it has leading "v"

### Release

* **v0.2.0-alpha.2:** prepare release


<a name="v0.2.0-alpha.1"></a>
## [v0.2.0-alpha.1](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.0...v0.2.0-alpha.1) (2020-08-17)

### Chore

* update to new talos modules
* update drone pipeline type
* update talos pkg import

### Fix

* ensure proper ownership of certs
* ensure machine configs work in packet

### Release

* **v0.2.0-alpha.1:** prepare release
