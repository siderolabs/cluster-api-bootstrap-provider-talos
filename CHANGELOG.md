## [CAPI Bootstrap Provider Talos 0.6.9](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.9) (2025-05-02)

Welcome to the v0.6.9 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Contributors

* Noel Georgi

### Changes
<details><summary>1 commit</summary>
<p>

* [`b7a2f69`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/b7a2f69f323d319907c0fb0cdb63fa3de62c040a) fix(ci): arm64 container image asset
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.6.8](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.8)

## [CAPI Bootstrap Provider Talos 0.6.8](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.8) (2025-05-01)

Welcome to the v0.6.8 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Cluster API

CABPT is now built and tested with Cluster API 1.10.0.


### Talos Linux

CABPT now supports Talos Linux v1.10.x machine configuration generation.


### Contributors

* Andrey Smirnov
* Chris
* Christian Bendieck
* Dmitriy Matrenichev
* Noel Georgi

### Changes
<details><summary>5 commits</summary>
<p>

* [`16c6183`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/16c6183db712a8f2a4e631f0511101379a385c12) feat: update Talos to 1.10.0, CAPI to 1.10.0
* [`636868b`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/636868bcaf14f8d9a12761c4bcd95a255b124ef0) feat: update Talos to 1.10-beta.0, CAPI to 1.10-rc.1
* [`7fcb5b3`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/7fcb5b3859ca024d7276b32664d23d65493b4a91) feat: use kres to manage github actions
* [`0044f9b`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/0044f9b027a0d2ed437f18fcc80d6b6c398e1583) fix: use net.JoinHostPort for control plane endpoint URL
* [`e135465`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e1354658114dfc6f202369b3214dc0929146c84c) feat: add infrastructurename hostname source
</p>
</details>

### Changes since v0.6.8-alpha.2
<details><summary>5 commits</summary>
<p>

* [`16c6183`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/16c6183db712a8f2a4e631f0511101379a385c12) feat: update Talos to 1.10.0, CAPI to 1.10.0
* [`636868b`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/636868bcaf14f8d9a12761c4bcd95a255b124ef0) feat: update Talos to 1.10-beta.0, CAPI to 1.10-rc.1
* [`7fcb5b3`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/7fcb5b3859ca024d7276b32664d23d65493b4a91) feat: use kres to manage github actions
* [`0044f9b`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/0044f9b027a0d2ed437f18fcc80d6b6c398e1583) fix: use net.JoinHostPort for control plane endpoint URL
* [`e135465`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e1354658114dfc6f202369b3214dc0929146c84c) feat: add infrastructurename hostname source
</p>
</details>

### Changes from siderolabs/crypto
<details><summary>1 commit</summary>
<p>

* [`0d45dee`](https://github.com/siderolabs/crypto/commit/0d45deefbcdd4bd6b6e549433b859083df55fc16) chore: bump deps
</p>
</details>

### Changes from siderolabs/go-pointer
<details><summary>1 commit</summary>
<p>

* [`347ee9b`](https://github.com/siderolabs/go-pointer/commit/347ee9b78f625d420254f4ab01bb1d6174474bf4) chore: rekres, update dependencies
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**              v5.9.0 -> v5.9.11
* **github.com/google/go-cmp**                   v0.6.0 -> v0.7.0
* **github.com/siderolabs/crypto**               v0.5.0 -> v0.5.1
* **github.com/siderolabs/go-pointer**           v1.0.0 -> v1.0.1
* **github.com/siderolabs/talos/pkg/machinery**  v1.9.0 -> v1.10.0
* **github.com/spf13/pflag**                     v1.0.5 -> v1.0.6
* **golang.org/x/sys**                           v0.28.0 -> v0.32.0
* **k8s.io/api**                                 v0.31.3 -> v0.32.3
* **k8s.io/apiextensions-apiserver**             v0.31.3 -> v0.32.3
* **k8s.io/apimachinery**                        v0.31.3 -> v0.32.3
* **k8s.io/client-go**                           v0.31.3 -> v0.32.3
* **k8s.io/component-base**                      v0.31.3 -> v0.32.3
* **sigs.k8s.io/cluster-api**                    v1.9.0 -> v1.10.1
* **sigs.k8s.io/controller-runtime**             v0.19.3 -> v0.20.4

Previous release can be found at [v0.6.7](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.7)

## [CAPI Bootstrap Provider Talos 0.6.7](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.7) (2024-12-17)

Welcome to the v0.6.7 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Cluster API

CABPT is now built and tested with Cluster API 1.9.0.


### Talos Linux

CABPT now supports Talos Linux v1.9.x machine configuration generation.


### Contributors

* Andrey Smirnov
* Dmitriy Matrenichev

### Changes
<details><summary>2 commits</summary>
<p>

* [`da3e7c8`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/da3e7c89b6c2279ba7562045ac3e38edee163cba) feat: update Talos to 1.9.0 final
* [`7cd504f`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/7cd504f954143900264f47afa13b774d41e5fd2d) feat: update Talos v1.9.0-beta.0
</p>
</details>

### Changes from siderolabs/crypto
<details><summary>1 commit</summary>
<p>

* [`58b2f92`](https://github.com/siderolabs/crypto/commit/58b2f9291c7e763a7210cfa681f88a7fa2230bf3) chore: use HTTP/2 ALPN by default
</p>
</details>

### Dependency Changes

* **github.com/siderolabs/crypto**               v0.4.4 -> v0.5.0
* **github.com/siderolabs/talos/pkg/machinery**  v1.8.0 -> v1.9.0
* **github.com/stretchr/testify**                v1.9.0 -> v1.10.0
* **golang.org/x/sys**                           v0.25.0 -> v0.28.0
* **k8s.io/api**                                 v0.31.1 -> v0.31.3
* **k8s.io/apiextensions-apiserver**             v0.31.1 -> v0.31.3
* **k8s.io/client-go**                           v0.31.1 -> v0.31.3
* **k8s.io/component-base**                      v0.31.1 -> v0.31.3
* **sigs.k8s.io/cluster-api**                    v1.8.3 -> v1.9.0
* **sigs.k8s.io/controller-runtime**             v0.19.0 -> v0.19.3

Previous release can be found at [v0.6.6](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.6)

## [CAPI Bootstrap Provider Talos 0.6.6](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.6) (2024-09-23)

Welcome to the v0.6.6 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Talos Linux

CABPT now supports Talos Linux v1.8.x machine configuration generation.


### Contributors

* Andrey Smirnov

### Changes
<details><summary>3 commits</summary>
<p>

* [`7b6b1fa`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/7b6b1fa007799da4575405ee486b255a94e73a0d) fix: generate configs without comments/examples
* [`d08a0bb`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d08a0bbe388c08efedc8b2547c50e7359f8f87f0) fix: remove CA bundle
* [`c0a6152`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/c0a615238eef0dc78790e94500e5a449973e885e) feat: update Talos to 1.8.0-beta.0
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                    v1.4.1 -> v1.4.2
* **github.com/siderolabs/talos/pkg/machinery**  v1.7.0 -> v1.8.0-beta.0
* **golang.org/x/sys**                           v0.19.0 -> v0.25.0
* **k8s.io/api**                                 v0.29.3 -> v0.31.0
* **k8s.io/apiextensions-apiserver**             v0.29.3 -> v0.31.0
* **k8s.io/apimachinery**                        v0.29.3 -> v0.31.0
* **k8s.io/client-go**                           v0.29.3 -> v0.31.0
* **k8s.io/component-base**                      v0.29.3 -> v0.31.0
* **k8s.io/klog/v2**                             v2.110.1 -> v2.130.1
* **sigs.k8s.io/cluster-api**                    v1.7.0 -> v1.8.2
* **sigs.k8s.io/controller-runtime**             v0.17.3 -> v0.19.0

Previous release can be found at [v0.6.5](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.5)

## [CAPI Bootstrap Provider Talos 0.6.5](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.5) (2024-04-19)

Welcome to the v0.6.5 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Patches

CABPT now supports Talos Linux machine configuration strategic merge patches via 'strategicPatches' field on the `TalosConfig` CRD.


### Contributors

* Andrey Smirnov
* Ksawery Kuczyński

### Changes
<details><summary>3 commits</summary>
<p>

* [`d0de67b`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d0de67bae46c7969efa19bc8a35516fe196c3de5) feat: update Talos to final 1.7.0
* [`77a404e`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/77a404e76c05f218ba372cf847331c665a5f01b9) feat: support for strategic merge patches
* [`171daf4`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/171daf47dbfa01cc8ef4ab4e116db2385d13f5a5) feat: update Talos to v1.7.0-beta.1
</p>
</details>

### Changes from siderolabs/crypto
<details><summary>3 commits</summary>
<p>

* [`c240482`](https://github.com/siderolabs/crypto/commit/c2404820ab1c1346c76b5b0f9b7632ca9d51e547) feat: provide dynamic client CA matching
* [`2f4f911`](https://github.com/siderolabs/crypto/commit/2f4f911da321ade3cedacc3b6abfef5f119f7508) feat: add PEMEncodedCertificate wrapper
* [`1c94bb3`](https://github.com/siderolabs/crypto/commit/1c94bb3967a427ba52c779a1b705f5aea466dc57) chore: bump dependencies
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**              v5.7.0 -> v5.9.0
* **github.com/go-logr/logr**                    v1.3.0 -> v1.4.1
* **github.com/siderolabs/crypto**               v0.4.1 -> v0.4.4
* **github.com/siderolabs/talos/pkg/machinery**  v1.6.0 -> v1.7.0
* **github.com/stretchr/testify**                v1.8.4 -> v1.9.0
* **golang.org/x/sys**                           e4099bfacb8c -> v0.19.0
* **k8s.io/api**                                 v0.28.4 -> v0.29.3
* **k8s.io/apiextensions-apiserver**             v0.28.4 -> v0.29.3
* **k8s.io/apimachinery**                        v0.28.4 -> v0.29.3
* **k8s.io/client-go**                           v0.28.4 -> v0.29.3
* **k8s.io/component-base**                      v0.28.4 -> v0.29.3
* **k8s.io/klog/v2**                             v2.110.1 **_new_**
* **sigs.k8s.io/cluster-api**                    v1.6.0 -> v1.7.0
* **sigs.k8s.io/controller-runtime**             v0.16.3 -> v0.17.3

Previous release can be found at [v0.6.4](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.4)

## [CAPI Bootstrap Provider Talos 0.6.4](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.4) (2024-01-23)

Welcome to the v0.6.4 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Contributors

* Andrey Smirnov

### Changes
<details><summary>1 commit</summary>
<p>

* [`604978d`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/604978dc6f6b78c7e966542ed9bf89168e3d8a16) fix: set a default controller runtime log
</p>
</details>

### Dependency Changes

This release has no dependency changes

Previous release can be found at [v0.6.3](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.3)

## [CAPI Bootstrap Provider Talos 0.6.3](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.3) (2023-12-15)

Welcome to the v0.6.3 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Talos 1.6

CABPT now supports machine configuration generation for Talos 1.6.


### Contributors

* Andrey Smirnov

### Changes
<details><summary>1 commit</summary>
<p>

* [`540603a`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/540603ae0b3527c022db0514bcfaa862272a0dbe) feat: update to Talos 1.6.0
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**              v5.6.0 -> v5.7.0
* **github.com/go-logr/logr**                    v1.2.4 -> v1.3.0
* **github.com/google/go-cmp**                   v0.5.9 -> v0.6.0
* **github.com/siderolabs/talos/pkg/machinery**  v1.5.2 -> v1.6.0
* **golang.org/x/sys**                           v0.10.0 -> e4099bfacb8c
* **k8s.io/api**                                 v0.27.2 -> v0.28.4
* **k8s.io/apiextensions-apiserver**             v0.27.2 -> v0.28.4
* **k8s.io/apimachinery**                        v0.27.2 -> v0.28.4
* **k8s.io/client-go**                           v0.27.2 -> v0.28.4
* **k8s.io/component-base**                      v0.28.4 **_new_**
* **sigs.k8s.io/cluster-api**                    v1.5.0 -> v1.6.0
* **sigs.k8s.io/controller-runtime**             v0.15.0 -> v0.16.3

Previous release can be found at [v0.6.2](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.2)

## [CAPI Bootstrap Provider Talos 0.6.2](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.2) (2023-09-07)

Welcome to the v0.6.2 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Talos 1.5

CABPT now supports machine configuration generation for Talos 1.5.


### Contributors

* Andrey Smirnov

### Changes
<details><summary>1 commit</summary>
<p>

* [`100d3d5`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/100d3d521e3db90b6449752e2d1e51dff43d25bd) fix: update Talos machinery to 1.5.2
</p>
</details>

### Dependency Changes

* **github.com/siderolabs/talos/pkg/machinery**  v1.5.0 -> v1.5.2

Previous release can be found at [v0.6.1](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.1)

## [CAPI Bootstrap Provider Talos 0.6.1](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.1) (2023-08-17)

Welcome to the v0.6.1 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Talos 1.5

CABPT now supports machine configuration generation for Talos 1.5.


### Contributors

* Andrey Smirnov
* Andrey Smirnov
* Utku Ozdemir

### Changes
<details><summary>1 commit</summary>
<p>

* [`fc4ef4e`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/fc4ef4e6165ffdef65249cb519214b4867c4f88d) feat: update to Talos 1.5, latest CAPI
</p>
</details>

### Changes from siderolabs/crypto
<details><summary>2 commits</summary>
<p>

* [`8f77da3`](https://github.com/siderolabs/crypto/commit/8f77da30a5193d207a6660b562a273a06d73aae0) feat: add a method to load PEM key from file
* [`c03ff58`](https://github.com/siderolabs/crypto/commit/c03ff58af5051acb9b56e08377200324a3ea1d5e) feat: add a way to represent redacted x509 private keys
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                    v1.2.3 -> v1.2.4
* **github.com/siderolabs/crypto**               v0.4.0 -> v0.4.1
* **github.com/siderolabs/talos/pkg/machinery**  v1.4.0 -> v1.5.0
* **github.com/stretchr/testify**                v1.8.2 -> v1.8.4
* **golang.org/x/sys**                           v0.7.0 -> v0.10.0
* **k8s.io/api**                                 v0.26.1 -> v0.27.2
* **k8s.io/apiextensions-apiserver**             v0.26.1 -> v0.27.2
* **k8s.io/apimachinery**                        v0.26.1 -> v0.27.2
* **k8s.io/client-go**                           v0.26.1 -> v0.27.2
* **sigs.k8s.io/cluster-api**                    v1.4.1 -> v1.5.0
* **sigs.k8s.io/controller-runtime**             v0.14.5 -> v0.15.0

Previous release can be found at [v0.6.0](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.0)

## [CAPI Bootstrap Provider Talos 0.6.0](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.0) (2023-05-03)

Welcome to the v0.6.0 release of CAPI Bootstrap Provider Talos!



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Talos 1.4

CABPT now supports machine configuration generation for Talos 1.4.


### Contributors

* Andrey Smirnov
* Andrey Smirnov
* Alexey Palazhchenko
* Spencer Smith
* Noel Georgi
* Andrew Rynhard
* Artem Chernyshev
* Artem Chernyshev
* Benjamin Gentil
* Dmitriy Matrenichev
* Serge Logvinov

### Changes
<details><summary>22 commits</summary>
<p>

* [`fee35a4`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/fee35a4161cfe5b2ae0dee2fb4f8230db9cceaa8) release(v0.6.0-alpha.1): prepare release
* [`28f4212`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/28f4212c1025940a64ddf6a12331d3eeedd7398f) chore: add 0.6 series to CAPI metadata
* [`0c61a33`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/0c61a336f754cd933fc4b26675a08bbe7d03c002) release(v0.6.0-alpha.0): prepare release
* [`d25c6a4`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d25c6a46bf946ac0b7bb1365ecdb593031cec789) feat: update Talos to 1.4.0
* [`d3adcdb`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d3adcdba9356331c5243e8dcbd3e9afdf6ba08ac) chore: bump dependencies
* [`6c9d018`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/6c9d018f1a1c908c9cd514f717aadcd62a829c4f) feat: add Tilt support
* [`376eb01`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/376eb01b82ff65f26fae1ef8df7d2301ea900585) feat: update CABPT to Talos 1.3.0
* [`4f2f856`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/4f2f85639f84521b99e830bc60dffc3ba574343d) feat: update to Talos 1.2.0
* [`a7fef2c`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/a7fef2c021a48c4436508d7ff1b7dce27e61be39) feat: update Talos to 1.2.0-beta.2
* [`2f3b21f`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/2f3b21f0926f08f4e480dac26e4db97699ee005b) feat: bump Talos to 1.1.0
* [`8b180df`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/8b180dfd87358f9d4d41f277aad4d592a7a8a1a5) feat: make `talosconfig` and `talosconfigtemplate` immutable
* [`e66b203`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e66b2037f81ef5a07540929184c305fcaf7c2bab) docs: update README for Talos 1.0
* [`ff9d1e8`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/ff9d1e86ee0731cfb3fea2f994a89b4633923d78) feat: update to Talos 1.0
* [`4eb3093`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/4eb30934d3e1cd29fd79d768f8c0ec3ae5151f33) chore: update after org rename
* [`e3a1f5a`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e3a1f5afc8d8af3deeb483c9f2b64e9a60c31a87) docs: add note for clusterctl rename bug
* [`7a4bc89`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/7a4bc890a8d0ea68c4a08451e3e767157d6e008f) chore: update GPG org
* [`3bc5406`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/3bc5406fcbbfac8024771b7e7b11228dd798fed1) chore: bump cert-manager to v1
* [`f2b1060`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/f2b1060ce79dfc7b271b353053911efbb8b0356c) chore: bump CAPI to 1.0.4
* [`b27f976`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/b27f9765e87e9c3aa6730286579a4ca4ba05d384) feat: add readiness/liveness checks
* [`c7a7265`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/c7a7265d3b86c69fc8d35bf39a62c4ea719c9a25) feat: support setting hostname to the machine name
* [`36fb7cc`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/36fb7cc14e4523e50a32d4e1a0a22d8085f361f9) fix: ensure shebang on packet machine configs
* [`8e39bd7`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/8e39bd731f06f14e642d7ffad202aec44948a492) feat: update Talos to 0.14.0
</p>
</details>

### Changes since v0.6.0-alpha.1
<details><summary>0 commit</summary>
<p>

</p>
</details>

### Changes from siderolabs/crypto
<details><summary>27 commits</summary>
<p>

* [`c3225ee`](https://github.com/siderolabs/crypto/commit/c3225eee603a8d1218c67e1bfe33ddde7953ed74) feat: allow CSR template subject field to be overridden
* [`8570669`](https://github.com/siderolabs/crypto/commit/85706698dac8cddd0e9f41006bed059347d2ea26) chore: rename to siderolabs/crypto
* [`e9df1b8`](https://github.com/siderolabs/crypto/commit/e9df1b8ca74c6efdc7f72191e5d2613830162fd5) feat: add support for generating keys from RSA-SHA256 CAs
* [`510b0d2`](https://github.com/siderolabs/crypto/commit/510b0d2753a89170d0c0f60e052a66484997a5b2) chore: add json tags
* [`6fa2d93`](https://github.com/siderolabs/crypto/commit/6fa2d93d0382299d5471e0de8e831c923398aaa8) fix: deepcopy nil fields as `nil`
* [`9a63cba`](https://github.com/siderolabs/crypto/commit/9a63cba8dabd278f3080fa8c160613efc48c43f8) fix: add back support for generating ECDSA keys with P-256 and SHA512
* [`893bc66`](https://github.com/siderolabs/crypto/commit/893bc66e4716a4cb7d1d5e66b5660ffc01f22823) fix: use SHA256 for ECDSA-P256
* [`deec8d4`](https://github.com/siderolabs/crypto/commit/deec8d47700e10e3ea813bdce01377bd93c83367) chore: implement DeepCopy methods for PEMEncoded* types
* [`d3cb772`](https://github.com/siderolabs/crypto/commit/d3cb77220384b3a3119a6f3ddb1340bbc811f1d1) feat: make possible to change KeyUsage
* [`6bc5bb5`](https://github.com/siderolabs/crypto/commit/6bc5bb50c52767296a1b1cab6580e3fcf1358f34) chore: remove unused argument
* [`cd18ef6`](https://github.com/siderolabs/crypto/commit/cd18ef62eb9f65d8b6730a2eb73e47e629949e1b) feat: add support for several organizations
* [`97c888b`](https://github.com/siderolabs/crypto/commit/97c888b3924dd5ac70b8d30dd66b4370b5ab1edc) chore: add options to CSR
* [`7776057`](https://github.com/siderolabs/crypto/commit/7776057f5086157873f62f6a21ec23fa9fd86e05) chore: fix typos
* [`80df078`](https://github.com/siderolabs/crypto/commit/80df078327030af7e822668405bb4853c512bd7c) chore: remove named result parameters
* [`15bdd28`](https://github.com/siderolabs/crypto/commit/15bdd282b74ac406ab243853c1b50338a1bc29d0) chore: minor updates
* [`4f80b97`](https://github.com/siderolabs/crypto/commit/4f80b976b640d773fb025d981bf85bcc8190815b) fix: verify CSR signature before issuing a certificate
* [`39584f1`](https://github.com/siderolabs/crypto/commit/39584f1b6e54e9966db1f16369092b2215707134) feat: support for key/certificate types RSA, Ed25519, ECDSA
* [`cf75519`](https://github.com/siderolabs/crypto/commit/cf75519cab82bd1b128ae9b45107c6bb422bd96a) fix: function NewKeyPair should create certificate with proper subject
* [`751c95a`](https://github.com/siderolabs/crypto/commit/751c95aa9434832a74deb6884cff7c5fd785db0b) feat: add 'PEMEncodedKey' which allows to transport keys in YAML
* [`562c3b6`](https://github.com/siderolabs/crypto/commit/562c3b66f89866746c0ba47927c55f41afed0f7f) feat: add support for public RSA key in RSAKey
* [`bda0e9c`](https://github.com/siderolabs/crypto/commit/bda0e9c24e80c658333822e2002e0bc671ac53a3) feat: enable more conversions between encoded and raw versions
* [`e0dd56a`](https://github.com/siderolabs/crypto/commit/e0dd56ac47456f85c0b247999afa93fb87ebc78b) feat: add NotBefore option for x509 cert creation
* [`12a4897`](https://github.com/siderolabs/crypto/commit/12a489768a6bb2c13e16e54617139c980f99a658) feat: add support for SPKI fingerprint generation and matching
* [`d0c3eef`](https://github.com/siderolabs/crypto/commit/d0c3eef149ec9b713e7eca8c35a6214bd0a64bc4) fix: implement NewKeyPair
* [`196679e`](https://github.com/siderolabs/crypto/commit/196679e9ec77cb709db54879ddeddd4eaafaea01) feat: move `pkg/grpc/tls` from `github.com/talos-systems/talos` as `./tls`
* [`1ff6242`](https://github.com/siderolabs/crypto/commit/1ff6242c91bb298ceeb4acd65685cba952fe4178) chore: initial version as imported from talos-systems/talos
* [`835063e`](https://github.com/siderolabs/crypto/commit/835063e055b28a525038b826a6d80cbe76402414) chore: initial commit
</p>
</details>

### Changes from siderolabs/go-pointer
<details><summary>2 commits</summary>
<p>

* [`71ccdf0`](https://github.com/siderolabs/go-pointer/commit/71ccdf0d65330596f4def36da37625e4f362f2a9) chore: implement main functionality
* [`c1c3b23`](https://github.com/siderolabs/go-pointer/commit/c1c3b235d30cb0de97ed0645809f2b21af3b021e) Initial commit
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**              v4.11.0 -> v5.6.0
* **github.com/go-logr/logr**                    v0.4.0 -> v1.2.3
* **github.com/google/go-cmp**                   v0.5.9 **_new_**
* **github.com/siderolabs/crypto**               v0.4.0 **_new_**
* **github.com/siderolabs/go-pointer**           v1.0.0 **_new_**
* **github.com/siderolabs/talos/pkg/machinery**  v1.4.0 **_new_**
* **github.com/stretchr/testify**                v1.7.0 -> v1.8.2
* **golang.org/x/sys**                           39ccf1dd6fa6 -> v0.7.0
* **k8s.io/api**                                 v0.22.2 -> v0.26.1
* **k8s.io/apiextensions-apiserver**             v0.22.2 -> v0.26.1
* **k8s.io/apimachinery**                        v0.22.2 -> v0.26.1
* **k8s.io/client-go**                           v0.22.2 -> v0.26.1
* **sigs.k8s.io/cluster-api**                    v1.0.0 -> v1.4.1
* **sigs.k8s.io/controller-runtime**             v0.10.2 -> v0.14.5

Previous release can be found at [v0.5.0](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.5.0)

## [CAPI Bootstrap Provider Talos 0.6.0-alpha.1](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.0-alpha.1) (2023-04-19)

Welcome to the v0.6.0-alpha.1 release of CAPI Bootstrap Provider Talos!  
*This is a pre-release of CAPI Bootstrap Provider Talos*



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Talos 1.4

CABPT now supports machine configuration generation for Talos 1.4.


### Contributors

* Andrey Smirnov
* Andrey Smirnov
* Alexey Palazhchenko
* Spencer Smith
* Noel Georgi
* Andrew Rynhard
* Artem Chernyshev
* Artem Chernyshev
* Benjamin Gentil
* Dmitriy Matrenichev
* Serge Logvinov

### Changes
<details><summary>21 commits</summary>
<p>

* [`28f4212`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/28f4212c1025940a64ddf6a12331d3eeedd7398f) chore: add 0.6 series to CAPI metadata
* [`0c61a33`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/0c61a336f754cd933fc4b26675a08bbe7d03c002) release(v0.6.0-alpha.0): prepare release
* [`d25c6a4`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d25c6a46bf946ac0b7bb1365ecdb593031cec789) feat: update Talos to 1.4.0
* [`d3adcdb`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d3adcdba9356331c5243e8dcbd3e9afdf6ba08ac) chore: bump dependencies
* [`6c9d018`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/6c9d018f1a1c908c9cd514f717aadcd62a829c4f) feat: add Tilt support
* [`376eb01`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/376eb01b82ff65f26fae1ef8df7d2301ea900585) feat: update CABPT to Talos 1.3.0
* [`4f2f856`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/4f2f85639f84521b99e830bc60dffc3ba574343d) feat: update to Talos 1.2.0
* [`a7fef2c`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/a7fef2c021a48c4436508d7ff1b7dce27e61be39) feat: update Talos to 1.2.0-beta.2
* [`2f3b21f`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/2f3b21f0926f08f4e480dac26e4db97699ee005b) feat: bump Talos to 1.1.0
* [`8b180df`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/8b180dfd87358f9d4d41f277aad4d592a7a8a1a5) feat: make `talosconfig` and `talosconfigtemplate` immutable
* [`e66b203`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e66b2037f81ef5a07540929184c305fcaf7c2bab) docs: update README for Talos 1.0
* [`ff9d1e8`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/ff9d1e86ee0731cfb3fea2f994a89b4633923d78) feat: update to Talos 1.0
* [`4eb3093`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/4eb30934d3e1cd29fd79d768f8c0ec3ae5151f33) chore: update after org rename
* [`e3a1f5a`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e3a1f5afc8d8af3deeb483c9f2b64e9a60c31a87) docs: add note for clusterctl rename bug
* [`7a4bc89`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/7a4bc890a8d0ea68c4a08451e3e767157d6e008f) chore: update GPG org
* [`3bc5406`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/3bc5406fcbbfac8024771b7e7b11228dd798fed1) chore: bump cert-manager to v1
* [`f2b1060`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/f2b1060ce79dfc7b271b353053911efbb8b0356c) chore: bump CAPI to 1.0.4
* [`b27f976`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/b27f9765e87e9c3aa6730286579a4ca4ba05d384) feat: add readiness/liveness checks
* [`c7a7265`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/c7a7265d3b86c69fc8d35bf39a62c4ea719c9a25) feat: support setting hostname to the machine name
* [`36fb7cc`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/36fb7cc14e4523e50a32d4e1a0a22d8085f361f9) fix: ensure shebang on packet machine configs
* [`8e39bd7`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/8e39bd731f06f14e642d7ffad202aec44948a492) feat: update Talos to 0.14.0
</p>
</details>

### Changes since v0.6.0-alpha.0
<details><summary>1 commit</summary>
<p>

* [`28f4212`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/28f4212c1025940a64ddf6a12331d3eeedd7398f) chore: add 0.6 series to CAPI metadata
</p>
</details>

### Changes from siderolabs/crypto
<details><summary>27 commits</summary>
<p>

* [`c3225ee`](https://github.com/siderolabs/crypto/commit/c3225eee603a8d1218c67e1bfe33ddde7953ed74) feat: allow CSR template subject field to be overridden
* [`8570669`](https://github.com/siderolabs/crypto/commit/85706698dac8cddd0e9f41006bed059347d2ea26) chore: rename to siderolabs/crypto
* [`e9df1b8`](https://github.com/siderolabs/crypto/commit/e9df1b8ca74c6efdc7f72191e5d2613830162fd5) feat: add support for generating keys from RSA-SHA256 CAs
* [`510b0d2`](https://github.com/siderolabs/crypto/commit/510b0d2753a89170d0c0f60e052a66484997a5b2) chore: add json tags
* [`6fa2d93`](https://github.com/siderolabs/crypto/commit/6fa2d93d0382299d5471e0de8e831c923398aaa8) fix: deepcopy nil fields as `nil`
* [`9a63cba`](https://github.com/siderolabs/crypto/commit/9a63cba8dabd278f3080fa8c160613efc48c43f8) fix: add back support for generating ECDSA keys with P-256 and SHA512
* [`893bc66`](https://github.com/siderolabs/crypto/commit/893bc66e4716a4cb7d1d5e66b5660ffc01f22823) fix: use SHA256 for ECDSA-P256
* [`deec8d4`](https://github.com/siderolabs/crypto/commit/deec8d47700e10e3ea813bdce01377bd93c83367) chore: implement DeepCopy methods for PEMEncoded* types
* [`d3cb772`](https://github.com/siderolabs/crypto/commit/d3cb77220384b3a3119a6f3ddb1340bbc811f1d1) feat: make possible to change KeyUsage
* [`6bc5bb5`](https://github.com/siderolabs/crypto/commit/6bc5bb50c52767296a1b1cab6580e3fcf1358f34) chore: remove unused argument
* [`cd18ef6`](https://github.com/siderolabs/crypto/commit/cd18ef62eb9f65d8b6730a2eb73e47e629949e1b) feat: add support for several organizations
* [`97c888b`](https://github.com/siderolabs/crypto/commit/97c888b3924dd5ac70b8d30dd66b4370b5ab1edc) chore: add options to CSR
* [`7776057`](https://github.com/siderolabs/crypto/commit/7776057f5086157873f62f6a21ec23fa9fd86e05) chore: fix typos
* [`80df078`](https://github.com/siderolabs/crypto/commit/80df078327030af7e822668405bb4853c512bd7c) chore: remove named result parameters
* [`15bdd28`](https://github.com/siderolabs/crypto/commit/15bdd282b74ac406ab243853c1b50338a1bc29d0) chore: minor updates
* [`4f80b97`](https://github.com/siderolabs/crypto/commit/4f80b976b640d773fb025d981bf85bcc8190815b) fix: verify CSR signature before issuing a certificate
* [`39584f1`](https://github.com/siderolabs/crypto/commit/39584f1b6e54e9966db1f16369092b2215707134) feat: support for key/certificate types RSA, Ed25519, ECDSA
* [`cf75519`](https://github.com/siderolabs/crypto/commit/cf75519cab82bd1b128ae9b45107c6bb422bd96a) fix: function NewKeyPair should create certificate with proper subject
* [`751c95a`](https://github.com/siderolabs/crypto/commit/751c95aa9434832a74deb6884cff7c5fd785db0b) feat: add 'PEMEncodedKey' which allows to transport keys in YAML
* [`562c3b6`](https://github.com/siderolabs/crypto/commit/562c3b66f89866746c0ba47927c55f41afed0f7f) feat: add support for public RSA key in RSAKey
* [`bda0e9c`](https://github.com/siderolabs/crypto/commit/bda0e9c24e80c658333822e2002e0bc671ac53a3) feat: enable more conversions between encoded and raw versions
* [`e0dd56a`](https://github.com/siderolabs/crypto/commit/e0dd56ac47456f85c0b247999afa93fb87ebc78b) feat: add NotBefore option for x509 cert creation
* [`12a4897`](https://github.com/siderolabs/crypto/commit/12a489768a6bb2c13e16e54617139c980f99a658) feat: add support for SPKI fingerprint generation and matching
* [`d0c3eef`](https://github.com/siderolabs/crypto/commit/d0c3eef149ec9b713e7eca8c35a6214bd0a64bc4) fix: implement NewKeyPair
* [`196679e`](https://github.com/siderolabs/crypto/commit/196679e9ec77cb709db54879ddeddd4eaafaea01) feat: move `pkg/grpc/tls` from `github.com/talos-systems/talos` as `./tls`
* [`1ff6242`](https://github.com/siderolabs/crypto/commit/1ff6242c91bb298ceeb4acd65685cba952fe4178) chore: initial version as imported from talos-systems/talos
* [`835063e`](https://github.com/siderolabs/crypto/commit/835063e055b28a525038b826a6d80cbe76402414) chore: initial commit
</p>
</details>

### Changes from siderolabs/go-pointer
<details><summary>2 commits</summary>
<p>

* [`71ccdf0`](https://github.com/siderolabs/go-pointer/commit/71ccdf0d65330596f4def36da37625e4f362f2a9) chore: implement main functionality
* [`c1c3b23`](https://github.com/siderolabs/go-pointer/commit/c1c3b235d30cb0de97ed0645809f2b21af3b021e) Initial commit
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**              v4.11.0 -> v5.6.0
* **github.com/go-logr/logr**                    v0.4.0 -> v1.2.3
* **github.com/google/go-cmp**                   v0.5.9 **_new_**
* **github.com/siderolabs/crypto**               v0.4.0 **_new_**
* **github.com/siderolabs/go-pointer**           v1.0.0 **_new_**
* **github.com/siderolabs/talos/pkg/machinery**  v1.4.0 **_new_**
* **github.com/stretchr/testify**                v1.7.0 -> v1.8.2
* **golang.org/x/sys**                           39ccf1dd6fa6 -> v0.7.0
* **k8s.io/api**                                 v0.22.2 -> v0.26.1
* **k8s.io/apiextensions-apiserver**             v0.22.2 -> v0.26.1
* **k8s.io/apimachinery**                        v0.22.2 -> v0.26.1
* **k8s.io/client-go**                           v0.22.2 -> v0.26.1
* **sigs.k8s.io/cluster-api**                    v1.0.0 -> v1.4.1
* **sigs.k8s.io/controller-runtime**             v0.10.2 -> v0.14.5

Previous release can be found at [v0.5.0](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.5.0)

## [CAPI Bootstrap Provider Talos 0.6.0-alpha.0](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.6.0-alpha.0) (2023-04-19)

Welcome to the v0.6.0-alpha.0 release of CAPI Bootstrap Provider Talos!  
*This is a pre-release of CAPI Bootstrap Provider Talos*



Please try out the release binaries and report any issues at
https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/issues.

### Talos 1.4

CABPT now supports machine configuration generation for Talos 1.4.


### Contributors

* Andrey Smirnov
* Andrey Smirnov
* Alexey Palazhchenko
* Spencer Smith
* Noel Georgi
* Andrew Rynhard
* Artem Chernyshev
* Artem Chernyshev
* Benjamin Gentil
* Dmitriy Matrenichev
* Serge Logvinov

### Changes
<details><summary>19 commits</summary>
<p>

* [`d25c6a4`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d25c6a46bf946ac0b7bb1365ecdb593031cec789) feat: update Talos to 1.4.0
* [`d3adcdb`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/d3adcdba9356331c5243e8dcbd3e9afdf6ba08ac) chore: bump dependencies
* [`6c9d018`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/6c9d018f1a1c908c9cd514f717aadcd62a829c4f) feat: add Tilt support
* [`376eb01`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/376eb01b82ff65f26fae1ef8df7d2301ea900585) feat: update CABPT to Talos 1.3.0
* [`4f2f856`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/4f2f85639f84521b99e830bc60dffc3ba574343d) feat: update to Talos 1.2.0
* [`a7fef2c`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/a7fef2c021a48c4436508d7ff1b7dce27e61be39) feat: update Talos to 1.2.0-beta.2
* [`2f3b21f`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/2f3b21f0926f08f4e480dac26e4db97699ee005b) feat: bump Talos to 1.1.0
* [`8b180df`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/8b180dfd87358f9d4d41f277aad4d592a7a8a1a5) feat: make `talosconfig` and `talosconfigtemplate` immutable
* [`e66b203`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e66b2037f81ef5a07540929184c305fcaf7c2bab) docs: update README for Talos 1.0
* [`ff9d1e8`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/ff9d1e86ee0731cfb3fea2f994a89b4633923d78) feat: update to Talos 1.0
* [`4eb3093`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/4eb30934d3e1cd29fd79d768f8c0ec3ae5151f33) chore: update after org rename
* [`e3a1f5a`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/e3a1f5afc8d8af3deeb483c9f2b64e9a60c31a87) docs: add note for clusterctl rename bug
* [`7a4bc89`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/7a4bc890a8d0ea68c4a08451e3e767157d6e008f) chore: update GPG org
* [`3bc5406`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/3bc5406fcbbfac8024771b7e7b11228dd798fed1) chore: bump cert-manager to v1
* [`f2b1060`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/f2b1060ce79dfc7b271b353053911efbb8b0356c) chore: bump CAPI to 1.0.4
* [`b27f976`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/b27f9765e87e9c3aa6730286579a4ca4ba05d384) feat: add readiness/liveness checks
* [`c7a7265`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/c7a7265d3b86c69fc8d35bf39a62c4ea719c9a25) feat: support setting hostname to the machine name
* [`36fb7cc`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/36fb7cc14e4523e50a32d4e1a0a22d8085f361f9) fix: ensure shebang on packet machine configs
* [`8e39bd7`](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/commit/8e39bd731f06f14e642d7ffad202aec44948a492) feat: update Talos to 0.14.0
</p>
</details>

### Changes from siderolabs/crypto
<details><summary>27 commits</summary>
<p>

* [`c3225ee`](https://github.com/siderolabs/crypto/commit/c3225eee603a8d1218c67e1bfe33ddde7953ed74) feat: allow CSR template subject field to be overridden
* [`8570669`](https://github.com/siderolabs/crypto/commit/85706698dac8cddd0e9f41006bed059347d2ea26) chore: rename to siderolabs/crypto
* [`e9df1b8`](https://github.com/siderolabs/crypto/commit/e9df1b8ca74c6efdc7f72191e5d2613830162fd5) feat: add support for generating keys from RSA-SHA256 CAs
* [`510b0d2`](https://github.com/siderolabs/crypto/commit/510b0d2753a89170d0c0f60e052a66484997a5b2) chore: add json tags
* [`6fa2d93`](https://github.com/siderolabs/crypto/commit/6fa2d93d0382299d5471e0de8e831c923398aaa8) fix: deepcopy nil fields as `nil`
* [`9a63cba`](https://github.com/siderolabs/crypto/commit/9a63cba8dabd278f3080fa8c160613efc48c43f8) fix: add back support for generating ECDSA keys with P-256 and SHA512
* [`893bc66`](https://github.com/siderolabs/crypto/commit/893bc66e4716a4cb7d1d5e66b5660ffc01f22823) fix: use SHA256 for ECDSA-P256
* [`deec8d4`](https://github.com/siderolabs/crypto/commit/deec8d47700e10e3ea813bdce01377bd93c83367) chore: implement DeepCopy methods for PEMEncoded* types
* [`d3cb772`](https://github.com/siderolabs/crypto/commit/d3cb77220384b3a3119a6f3ddb1340bbc811f1d1) feat: make possible to change KeyUsage
* [`6bc5bb5`](https://github.com/siderolabs/crypto/commit/6bc5bb50c52767296a1b1cab6580e3fcf1358f34) chore: remove unused argument
* [`cd18ef6`](https://github.com/siderolabs/crypto/commit/cd18ef62eb9f65d8b6730a2eb73e47e629949e1b) feat: add support for several organizations
* [`97c888b`](https://github.com/siderolabs/crypto/commit/97c888b3924dd5ac70b8d30dd66b4370b5ab1edc) chore: add options to CSR
* [`7776057`](https://github.com/siderolabs/crypto/commit/7776057f5086157873f62f6a21ec23fa9fd86e05) chore: fix typos
* [`80df078`](https://github.com/siderolabs/crypto/commit/80df078327030af7e822668405bb4853c512bd7c) chore: remove named result parameters
* [`15bdd28`](https://github.com/siderolabs/crypto/commit/15bdd282b74ac406ab243853c1b50338a1bc29d0) chore: minor updates
* [`4f80b97`](https://github.com/siderolabs/crypto/commit/4f80b976b640d773fb025d981bf85bcc8190815b) fix: verify CSR signature before issuing a certificate
* [`39584f1`](https://github.com/siderolabs/crypto/commit/39584f1b6e54e9966db1f16369092b2215707134) feat: support for key/certificate types RSA, Ed25519, ECDSA
* [`cf75519`](https://github.com/siderolabs/crypto/commit/cf75519cab82bd1b128ae9b45107c6bb422bd96a) fix: function NewKeyPair should create certificate with proper subject
* [`751c95a`](https://github.com/siderolabs/crypto/commit/751c95aa9434832a74deb6884cff7c5fd785db0b) feat: add 'PEMEncodedKey' which allows to transport keys in YAML
* [`562c3b6`](https://github.com/siderolabs/crypto/commit/562c3b66f89866746c0ba47927c55f41afed0f7f) feat: add support for public RSA key in RSAKey
* [`bda0e9c`](https://github.com/siderolabs/crypto/commit/bda0e9c24e80c658333822e2002e0bc671ac53a3) feat: enable more conversions between encoded and raw versions
* [`e0dd56a`](https://github.com/siderolabs/crypto/commit/e0dd56ac47456f85c0b247999afa93fb87ebc78b) feat: add NotBefore option for x509 cert creation
* [`12a4897`](https://github.com/siderolabs/crypto/commit/12a489768a6bb2c13e16e54617139c980f99a658) feat: add support for SPKI fingerprint generation and matching
* [`d0c3eef`](https://github.com/siderolabs/crypto/commit/d0c3eef149ec9b713e7eca8c35a6214bd0a64bc4) fix: implement NewKeyPair
* [`196679e`](https://github.com/siderolabs/crypto/commit/196679e9ec77cb709db54879ddeddd4eaafaea01) feat: move `pkg/grpc/tls` from `github.com/talos-systems/talos` as `./tls`
* [`1ff6242`](https://github.com/siderolabs/crypto/commit/1ff6242c91bb298ceeb4acd65685cba952fe4178) chore: initial version as imported from talos-systems/talos
* [`835063e`](https://github.com/siderolabs/crypto/commit/835063e055b28a525038b826a6d80cbe76402414) chore: initial commit
</p>
</details>

### Changes from siderolabs/go-pointer
<details><summary>2 commits</summary>
<p>

* [`71ccdf0`](https://github.com/siderolabs/go-pointer/commit/71ccdf0d65330596f4def36da37625e4f362f2a9) chore: implement main functionality
* [`c1c3b23`](https://github.com/siderolabs/go-pointer/commit/c1c3b235d30cb0de97ed0645809f2b21af3b021e) Initial commit
</p>
</details>

### Dependency Changes

* **github.com/evanphx/json-patch**              v4.11.0 -> v5.6.0
* **github.com/go-logr/logr**                    v0.4.0 -> v1.2.3
* **github.com/google/go-cmp**                   v0.5.9 **_new_**
* **github.com/siderolabs/crypto**               v0.4.0 **_new_**
* **github.com/siderolabs/go-pointer**           v1.0.0 **_new_**
* **github.com/siderolabs/talos/pkg/machinery**  v1.4.0 **_new_**
* **github.com/stretchr/testify**                v1.7.0 -> v1.8.2
* **golang.org/x/sys**                           39ccf1dd6fa6 -> v0.7.0
* **k8s.io/api**                                 v0.22.2 -> v0.26.1
* **k8s.io/apiextensions-apiserver**             v0.22.2 -> v0.26.1
* **k8s.io/apimachinery**                        v0.22.2 -> v0.26.1
* **k8s.io/client-go**                           v0.22.2 -> v0.26.1
* **sigs.k8s.io/cluster-api**                    v1.0.0 -> v1.4.1
* **sigs.k8s.io/controller-runtime**             v0.10.2 -> v0.14.5

Previous release can be found at [v0.5.0](https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/tag/v0.5.0)


## [CAPI Bootstrap Provider Talos 0.5.0-alpha.0](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/releases/tag/v0.5.0-alpha.0) (2021-10-27)

Welcome to the v0.5.0-alpha.0 release of CAPI Bootstrap Provider Talos!  
*This is a pre-release of CAPI Bootstrap Provider Talos*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/issues.

### CAPI v1beta1

CABPT now supports CAPI version 1.0.x (v1beta1).


### `talosconfig` Generation

CABPT now generates client-side Talos API credentials (`talosconfig`) in the `<cluster>-talosconfig` Secret.
Generated `talosconfig` will be updated with the endpoints of the control plane `Machine`s.


### Contributors

* Andrey Smirnov

### Changes
<details><summary>4 commits</summary>
<p>

* [`d124c07`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/d124c072c9db8d402b353a73646d2d197bae76a4) docs: update README with usage and compatibility matrix
* [`20792f3`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/20792f345b7ff3c8ffa9d65c9ca8dcab1932f49e) feat: generate talosconfig as a secret with proper endpoints
* [`abd206f`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/abd206fd8a98f5478f8ffd0f8686e32be3b7defe) feat: update to CAPI v1.0.x contract (v1beta1)
* [`b7faf9e`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/b7faf9e730b7c9f50ffa94be194ddcf908708a2c) feat: update Talos machinery to 0.13.0
</p>
</details>

### Changes from talos-systems/crypto
<details><summary>2 commits</summary>
<p>

* [`9a63cba`](https://github.com/talos-systems/crypto/commit/9a63cba8dabd278f3080fa8c160613efc48c43f8) fix: add back support for generating ECDSA keys with P-256 and SHA512
* [`893bc66`](https://github.com/talos-systems/crypto/commit/893bc66e4716a4cb7d1d5e66b5660ffc01f22823) fix: use SHA256 for ECDSA-P256
</p>
</details>

### Dependency Changes

* **github.com/talos-systems/crypto**  v0.3.2 -> v0.3.4
* **golang.org/x/sys**                 bfb29a6856f2 -> 39ccf1dd6fa6
* **inet.af/netaddr**                  85fa6c94624e **_new_**
* **k8s.io/api**                       v0.21.4 -> v0.22.2
* **k8s.io/apiextensions-apiserver**   v0.21.4 -> v0.22.2
* **k8s.io/apimachinery**              v0.21.4 -> v0.22.2
* **k8s.io/client-go**                 v0.21.4 -> v0.22.2
* **sigs.k8s.io/cluster-api**          v0.4.3 -> v1.0.0
* **sigs.k8s.io/controller-runtime**   v0.9.7 -> v0.10.2

Previous release can be found at [v0.4.0](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/releases/tag/v0.4.0)

## [CAPI Bootstrap Provider Talos 0.4.0-alpha.0](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/releases/tag/v0.4.0-alpha.0) (2021-10-01)

Welcome to the v0.4.0-alpha.0 release of CAPI Bootstrap Provider Talos!  
*This is a pre-release of CAPI Bootstrap Provider Talos*



Please try out the release binaries and report any issues at
https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/issues.

### CAPI v1alpha4

CABPT now supports CAPI v1alpha4.


### Contributors

* Andrey Smirnov
* Spencer Smith

### Changes
<details><summary>3 commits</summary>
<p>

* [`8c7fec8`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/8c7fec8e373bd12609f6274d79ca07d187212d91) fix: don't write incomplete `<cluster>-ca` secret for configtype none
* [`f46c83d`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f46c83d328ee44db2ccb5eef67b366cc73c13319) feat: bump Talos machinery to 0.12.3
* [`7b760cf`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/7b760cf69ecab93200821dded931171657a5dedc) feat: support CAPI v1alpha4
</p>
</details>

### Dependency Changes

* **github.com/go-logr/logr**                       v0.1.0 -> v0.4.0
* **github.com/talos-systems/talos/pkg/machinery**  7e63e43eb399 -> v0.12.3
* **golang.org/x/sys**                              0f9fa26af87c -> bfb29a6856f2
* **k8s.io/api**                                    v0.17.9 -> v0.21.4
* **k8s.io/apiextensions-apiserver**                v0.17.9 -> v0.21.4
* **k8s.io/apimachinery**                           v0.17.9 -> v0.21.4
* **k8s.io/client-go**                              v0.17.9 -> v0.21.4
* **sigs.k8s.io/cluster-api**                       v0.3.22 -> v0.4.3
* **sigs.k8s.io/controller-runtime**                v0.5.14 -> v0.9.7

Previous release can be found at [v0.3.0](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/releases/tag/v0.3.0)

## [CAPI Bootstrap Provider Talos 0.3.0-alpha.1](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/releases/tag/v0.3.0-alpha.1) (2021-09-21)

Welcome to the v0.3.0-alpha.1 release of CAPI Bootstrap Provider Talos!  
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
* Spencer Smith

### Changes
<details><summary>18 commits</summary>
<p>

* [`977121a`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/977121ad14dc0637f7c4282e69a4ee26e28372d4) fix: construct properly data secret name
* [`f8c75c8`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f8c75c89c4653de30165fb1538e906256a4eec66) fix: update metadata.yaml for v0.3 of CABPT
* [`db60f9e`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/db60f9eb0697c4949be9c00cf8dc7787d383bad2) release(v0.3.0-alpha.0): prepare release
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

### Changes since v0.3.0-alpha.0
<details><summary>2 commits</summary>
<p>

* [`977121a`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/977121ad14dc0637f7c4282e69a4ee26e28372d4) fix: construct properly data secret name
* [`f8c75c8`](https://github.com/talos-systems/cluster-api-bootstrap-provider-talos/commit/f8c75c89c4653de30165fb1538e906256a4eec66) fix: update metadata.yaml for v0.3 of CABPT
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
