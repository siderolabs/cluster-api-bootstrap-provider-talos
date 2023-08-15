# Cluster API Bootstrap Provider Talos (CABPT)

## Intro

The Cluster API Bootstrap Provider Talos (CABPT) is a project by [Sidero Labs](https://siderolabs.com/) that provides a [Cluster API](https://github.com/kubernetes-sigs/cluster-api) (CAPI) bootstrap provider for use in deploying Talos-based Kubernetes nodes across any environment.
Given some basic info, this provider will generate bootstrap configurations for a given machine and reconcile the necessary custom resources for CAPI to pick up the generated data.

## Corequisites

There are a few corequisites and assumptions that go into using this project:

- [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- [Talos](https://talos.dev/)

## Installing

CABPT provider should be installed alongside with [CACPPT](https://github.com/siderolabs/cluster-api-control-plane-provider-talos) provider.

```shell
clusterctl init --bootstrap talos --control-plane talos --infrastructure <infrastructure provider>
```

If you encounter the following error, this is caused by a rename of our GitHub org from `talos-systems` to `siderolabs`.

```bash
$ clusterctl init -b talos -c talos -i sidero
Fetching providers
Error: failed to get provider components for the "talos" provider: target namespace can't be defaulted. Please specify a target namespace
```

This can be worked around by adding the following to `~/.cluster-api/clusterctl.yaml` and rerunning the init command:

```yaml
providers:
  - name: "talos"
    url: "https://github.com/siderolabs/cluster-api-bootstrap-provider-talos/releases/latest/bootstrap-components.yaml"
    type: "BootstrapProvider"
  - name: "talos"
    url: "https://github.com/siderolabs/cluster-api-control-plane-provider-talos/releases/latest/control-plane-components.yaml"
    type: "ControlPlaneProvider"
  - name: "sidero"
    url: "https://github.com/siderolabs/sidero/releases/latest/infrastructure-components.yaml"
    type: "InfrastructureProvider"
```

## Compatibility

This provider's versions are compatible with the following versions of Cluster API:

|                | v1alpha3 (v0.3) | v1alpha4 (v0.4) | v1beta1 (v1.x) |
| -------------- | --------------- | --------------- | -------------- |
| CABPT (v0.5.x) |                 |                 | ✓              |
| CABPT (v0.6.x) |                 |                 | ✓              |

This provider's versions are able to install and manage the following versions of Kubernetes:

|                | v1.19 | v1.20 | v1.21 | v1.22 | v1.23 | v1.24 | v1.25 | v1.26 | v1.27 | v1.28 |
| -------------- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- | ----- |
| CABPT (v0.5.x) | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     | ✓     |       |       |
| CABPT (v0.6.x) |       |       |       |       |       | ✓     | ✓     | ✓     | ✓     | ✓     |

This provider's versions are compatible with the following versions of Talos:

|                  | v1.0  | v1.1  | v1.2  | v1.3  | v1.4  | v1.5  |
| ---------------- | ----- | ----- | ----- | ----- | ----- | ----- |
| CABPT (v0.5.x)   | ✓     | ✓     | ✓     | ✓     |       |       |
| CABPT (v0.6.x)   |       |       | ✓     | ✓     | ✓     | ✓     |

CABPT generates machine configuration compatible with Talos version specified in the `talosVersion:` field (see below).

## Usage

CABPT is not used directly, but rather via CACPPT (`TalosControlPlane`) for control plane nodes or via `MachineDeployment` (`MachinePool`) for worker nodes.
In either case, CABPT settings are passed via `TalosConfigTemplate` resource:

```yaml
apiVersion: controlplane.cluster.x-k8s.io/v1alpha3
kind: TalosControlPlane
metadata:
  name: mycluster-cp
spec:
  controlPlaneConfig:
    controlplane:
      generateType: controlplane
      talosVersion: v1.1
  ...
```

```yaml
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha3
kind: TalosConfigTemplate
metadata:
  name: mycluster-workers
spec:
  template:
    spec:
      generateType: worker
      talosVersion: v1.1
```

Fields available in the `TalosConfigTemplate` (and `TalosConfig`) resources:

- `generateType`: Talos machine configuration type to generate (`controlplane`, `init` (deprecated), `worker`) or `none` for user-supplied configuration (see below)
- `talosVersion`: version of Talos to generate machine configuration for (e.g. `v1.0`, patch version might be omitted).
   CABPT defaults to the latest supported Talos version, but it can generate configuration compatible with previous versions of Talos.
   It is recommended to always set this field explicitly to avoid issues when CABPT is upgraded to the version which supports new Talos version.
- `configPatches` (optional): set of machine configuration patches to apply to the generated configuration.
- `data` (only for `generateType: none`): user-supplied machine configuration.
- `hostname` (optional): configure hostname in the generate machine configuration:
  - `source` (`MachineName`): set the hostname in the generated machine configuration to the `Machine` name (not supported with `MachinePool` deployments)

### Generated Machine Configuration

When `generateType` is set to the machine type of the Talos nodes (`controlplane` for control plane nodes and `worker` for worker nodes), CABPT generates a set of cluster-wide
secrets which are used to provision machine configuration for each node.
Machine configuration generated is compatible with the Talos version set in the `talosVersion` field.

```yaml
spec:
  generateType: controlplane
  talosVersion: v1.5
```

### User-supplied Machine Configuration

In this mode CABPT passes through machine configuration set in `data` field as bootstrap data to the `Machine`.
Machine configuration can be generated with `talosctl gen config`.

```yaml
spec:
  generateType: none
  data: |
    version: v1alpha1
    machine:
      type: controlplane
    ...
    ...
    ...
```

### Configuration Patches

Machine configuration can be customized by applying configuration patches.
Any field of the [Talos machine configuration](https://www.talos.dev/docs/latest/reference/configuration/)
can be overridden on a per-machine basis using this method.
The format of these patches is based on [JSON 6902](http://jsonpatch.com/) that you may be used to in tools like `kustomize`.

```yaml
spec:
  generateType: controlplane
  talosVersion: v1.5
  configPatches:
    - op: replace
      path: /machine/install
      value:
        disk: /dev/sda
    - op: add
      path: /cluster/network/cni
      value:
        name: custom
        urls:
          - https://docs.projectcalico.org/v3.18/manifests/calico.yaml
```

### Retrieving `talosconfig`

Client-side `talosconfig` is required to access the cluster using Talos API.
CABPT generates `talosconfig` for generated machine configuration and stores it as `<cluster>-talosconfig` secret in cluster's namespace.

`talosconfig` can be retrieved with:

```shell
kubectl get secret --namespace <cluster-namespace> <cluster-name>-talosconfig -o jsonpath='{.data.talosconfig}' | base64 -d > cluster-talosconfig
talosctl config merge cluster-talosconfig
talosctl -n <IP> version
```

CABPT updates endpoints in the `talosconfig` based on control plane `Machine` addresses.

### Operation

CABPT reconciles `TalosConfig` resources.
Once `TalosConfig` and its associated `Machine` are ready, CABPT generates machine configuration and stores it in the `<machine>-bootstrap-data` Secret.
Kubernetes cluster CA is stored in the `<cluster>-ca` Secret.
Cluster-wide shared secrets are stored in the `<cluster>-talos` Secret.
Client-side Talos API configuration is stored in the `<cluster>-talosconfig` Secret.

As part of its operation, CABPT sets a number of Conditions on the `TalosConfig/Status` resource:

- `DataSecretAvailable`: CABPT generated machine configuration in `<machine>-bootstrap-data` Secret, CABPT unblocks infrastructure provider to boot the `Machine`.
- `ClientConfigAvailableCondition`: CABPT generated Talos client configuration in `<cluster>-talosconfig` Secret.

```shell
$ kubectl describe talosconfig talosconfig-cp-0
Status:
  Conditions:
    Last Transition Time:  2021-10-25T20:54:09Z
    Status:                True
    Type:                  Ready
    Last Transition Time:  2021-10-25T20:54:09Z
    Status:                True
    Type:                  ClientConfigAvailable
    Last Transition Time:  2021-10-25T20:54:09Z
    Status:                True
    Type:                  DataSecretAvailable
  Data Secret Name:        cp-0-bootstrap-data
```

If CABPT fails to perform its operations, it stores error message for the respective Condition.

```shell
Status:
  Conditions:
    Last Transition Time:  2021-10-25T20:54:09Z
    Message:               failure applying rfc6902 patches to talos machine config: add operation does not apply: doc is missing path: "/machine/time/servers": missing value
    Reason:                DataSecretGenerationFailed
    Severity:              Error
    Status:                False
```

These statuses are also presented in the `clusterctl describe cluster --show-conditions all` output.

## Building

This project can be built simply by running `make release` from the root directory.
Doing so will create a file called `_out/bootstrap-talos/<version>/bootstrap-components.yaml`.
If you wish, you can tweak settings by editing the release yaml.
This file can then be installed into your management cluster with `kubectl apply -f _out/bootstrap-components.yaml`.

## Support

Join our [Slack](https://slack.dev.talos-systems.io)!
