## [CAPI Bootstrap Provider Talos 0.3.0-alpha.0](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/releases/tag/v0.3.0-alpha.0) (2021-09-21)

Welcome to the v0.3.0-alpha.0 release of CAPI Bootstrap Provider Talos!  
*This is a pre-release of CAPI Bootstrap Provider Talos*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/issues.

### CAPI v1alpha3

This release of CABPT is compatible with CAPI v1alpha3 (v0.3.x).
Next release of CABPT will bring compatibility with CAPI v1alpha4 (v0.4.x).


### Default `talosVersion`

In previous releases of CABPT default value of `talosVersion` field was `v0.8`.
As Talos v0.8 release is almost a year old, new default value of `talosVersion` is to use whatever Talos version CABPT was
built against (in this relase, it's Talos 0.12).

If you're still running Talos v0.8.x, please make sure `talosVersion` is set to `v0.8`.


### Talos 0.12

CABPT supports config generation for Talos 0.12.
Talos majort version can be specified in the spec of `TalosControlPlane` or `MachineDeployment`:

```yaml
  generateType: controlplane
  talosVersion: v0.11
```

It is recommended to specify minor version of Talos to make sure machine configuration stays comptabile with Talos version
being used even if the CABPT is upgraded to new version.


### Contributors

* Alexey Palazhchenko
* Alexey Palazhchenko
* Andrey Smirnov
* Andrey Smirnov
* Serge Logvinov

### Changes
<details><summary>15 commits</summary>
<p>

* [`755a2dd`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/755a2dd90c3668db89f8eae14f60db4564764475) fix: update Talos machinery to 0.12, fix secrets persistence
* [`f91b032`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f91b032935776c1224f824cc860bfa4df5e220b1) fix: use bootstrap data secret names
* [`6bff239`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/6bff2393840655c2361def455b601511b86ba71f) chore: use Go 1.17
* [`56fb73b`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/56fb73b53f41b91b12ba2b3c331d7a04b7263a17) test: add test for the second machine
* [`e5b7738`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/e5b773833120fdd7ca4d57e0a0a4fe781495bf7e) test: add more tests
* [`bc4105d`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/bc4105d9e8366d4e840705a6cecfbc81bdcca00a) test: wait for CAPI availability
* [`c82b8ab`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/c82b8ab47bca5313cb96df1b70de0914da285331) chore: make versions configurable
* [`5594c96`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/5594c96daa55fb9fc9af585e8f2fc26551ce9bb5) chore: use codecov uploader from build-container
* [`cced038`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/cced038257d3eec5b7c48bc524de5165b5734496) chore: fix license headers
* [`7b5dc51`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/7b5dc51e83a54a1f5fa707c66a296ca9514c8722) chore: do not run tests on ARM
* [`d6258cf`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/d6258cf21778149a254d9669b03ac10bae9e0955) chore: improve tests runner
* [`c6ce363`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/c6ce36375ef145760647c632d64a9a3c93574e4b) chore: sign Drone CI configuration
* [`ad592d1`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/ad592d16fa8397f88a28e6a4151bc64b0a1c097d) chore: add basic integration test
* [`9fb0d07`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/9fb0d07ca4d2e8333b0b61ee0fe0ba3e6660489f) chore: add missing LICENSE file
* [`acf18d2`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/acf18d2bb09aab64687c1fccf1e628ef76e9cff8) chore: update machinery to v0.11.3
</p>
</details>

### Changes from talos-systems/crypto
<details><summary>8 commits</summary>
<p>

* [`deec8d4`](https://github.com/talos-systems/crypto/commit/deec8d47700e10e3ea813bdce01377bd93c83367) chore: implement DeepCopy methods for PEMEncoded* types
* [`d3cb772`](https://github.com/talos-systems/crypto/commit/d3cb77220384b3a3119a6f3ddb1340bbc811f1d1) feat: make possible to change KeyUsage
* [`6bc5bb5`](https://github.com/talos-systems/crypto/commit/6bc5bb50c52767296a1b1cab6580e3fcf1358f34) chore: remove unused argument
* [`cd18ef6`](https://github.com/talos-systems/crypto/commit/cd18ef62eb9f65d8b6730a2eb73e47e629949e1b) feat: add support for several organizations
* [`97c888b`](https://github.com/talos-systems/crypto/commit/97c888b3924dd5ac70b8d30dd66b4370b5ab1edc) chore: add options to CSR
* [`7776057`](https://github.com/talos-systems/crypto/commit/7776057f5086157873f62f6a21ec23fa9fd86e05) chore: fix typos
* [`80df078`](https://github.com/talos-systems/crypto/commit/80df078327030af7e822668405bb4853c512bd7c) chore: remove named result parameters
* [`15bdd28`](https://github.com/talos-systems/crypto/commit/15bdd282b74ac406ab243853c1b50338a1bc29d0) chore: minor updates
</p>
</details>

### Dependency Changes

* **github.com/AlekSi/pointer**                     v1.1.0 **_new_**
* **github.com/evanphx/json-patch**                 v4.9.0 -> v4.11.0
* **github.com/stretchr/testify**                   v1.7.0 **_new_**
* **github.com/talos-systems/crypto**               4f80b976b640 -> v0.3.2
* **github.com/talos-systems/talos/pkg/machinery**  828772cec9a3 -> 7e63e43eb399
* **golang.org/x/sys**                              0f9fa26af87c **_new_**
* **gopkg.in/yaml.v2**                              v2.3.0 -> v2.4.0
* **sigs.k8s.io/cluster-api**                       v0.3.12 -> v0.3.22

Previous release can be found at [v0.2.0](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/releases/tag/v0.2.0)


<a name="v0.2.0-alpha.12"></a>
## [v0.2.0-alpha.12](https://github.com/talos-systems/talos/compare/v0.2.0-alpha.11...v0.2.0-alpha.12) (2021-05-14)

### Chore

* rework build, move to ghcr.io, build for arm64/amd64

### Fix

* back down resource requests
* ensure secrets are deleted when cluster is dropped


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
