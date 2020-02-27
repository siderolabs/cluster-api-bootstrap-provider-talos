# cluster-api-bootstrap-provider-talos

## Intro

The Cluster API Bootstrap Provider Talos (CABPT) is a project by [Talos Systems](https://www.talos-systems.com/) that provides a [Cluster API](https://github.com/kubernetes-sigs/cluster-api)(CAPI) bootstrap provider for use in deploying Talos-based Kubernetes nodes across any environment.
Given some basic info, this provider will generate bootstrap configurations for a given machine and reconcile the necessary custom resources for CAPI to pick up the generated data.

## Corequisites

There are a few corequisites and assumptions that go into using this project:

- [Cluster API](https://github.com/kubernetes-sigs/cluster-api)
- [Cluster API Provider Metal](https://github.com/talos-systems/cluster-api-provider-metal) (optional)

## Building and Installing

This project can be built simply by running `make release` from the root directory.
Doing so will create a file called `_out/release.yaml`.
If you wish, you can tweak settings by editing the release yaml.
This file can then be installed into your management cluster with `kubectl apply -f _out/release.yaml`.

Note that CABPT should be deployed as part of a set of controllers for Cluster API.
You will need at least the upstream CAPI components and an infrastructure provider for v1alpha2 CAPI capabilities.

## Usage

CAPM supports a single API type, a TalosConfig.
You can create YAML definitions of a TalosConfig and `kubectl apply` them as part of a larger CAPI cluster deployment.
Below is a bare-minimum example.

A basic config:

```yaml
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha2
kind: TalosConfig
metadata:
  name: talos-0
  labels:
    cluster.x-k8s.io/cluster-name: talos
spec:
  generateType: init
```

Note the generateType mentioned above.
This is a required value in the spec for a TalosConfig.
For a no-frills bootstrap config, you can simply specify `init`, `controlplane`, or `worker` depending on what type of Talos node this is.
When creating a TalosConfig this way, you can then retrieve the talosconfig file that allows for osctl interaction with your nodes by doing something like `kubectl get talosconfig -o yaml talos-0 -o jsonpath='{.status.talosConfig}'` after creation.

If you wish to do something more complex, we allow for the ability to supply an entire Talos config file to the resource.
This can be done by setting the generateType to `none` and specifying a `data` field.
This config file can be generated with `osctl config generate` and the edited to supply the various options you may desire.
This full config is blindly copied from the `data` section of the spec and presented under `.status.bootstrapData` so that the upstream CAPI controllers can see it and make use.

An example of a more complex config:

```yaml
apiVersion: bootstrap.cluster.x-k8s.io/v1alpha2
kind: TalosConfig
metadata:
  name: talos-0
  labels:
    cluster.x-k8s.io/cluster-name: talos
spec:
  generateType: none
  data: |
    version: v1alpha1
    machine:
      type: init
      token: xxxxxx
    ...
    ...
    ...
```

Note that specifying the full config above removes the ability for our bootstrap provider to generate a talosconfig for use.
As such, you should keep track of the talosconfig that's generated when running `osctl config generate`.
