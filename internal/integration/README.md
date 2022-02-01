# Running integration tests

Bring up Talos cluster (or any other Kubernetes clusters).

Put `kubeconfig` to the root of the repository (for Talos: `talosctl -n <IP> kubeconfig -f kubeconfig`).

Create release with fixed tag: `make release-manifests TAG=v0.5.0` (update release if the CRDs are updated).

Run tests: `make test TAG=v0.5.0`.

Tests clean up after the run, so they can be run repeatedly against the cluster.
