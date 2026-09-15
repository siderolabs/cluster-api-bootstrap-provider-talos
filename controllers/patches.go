// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package controllers

import (
	"github.com/siderolabs/talos/pkg/machinery/config"
	configconfig "github.com/siderolabs/talos/pkg/machinery/config/config"
	"github.com/siderolabs/talos/pkg/machinery/config/container"
	"github.com/siderolabs/talos/pkg/machinery/config/encoder"
	"github.com/siderolabs/talos/pkg/machinery/config/types/k8s"
	"github.com/siderolabs/talos/pkg/machinery/config/types/meta"
	"gopkg.in/yaml.v2"
)

func patchPodServiceSubnets(versionContract *config.VersionContract, podBlocks, serviceBlocks []string) ([]byte, error) {
	if versionContract.MultidocKubernetesConfigSupported() {
		podSubnets := make([]meta.Prefix, len(podBlocks))

		for i := range podSubnets {
			if err := podSubnets[i].UnmarshalText([]byte(podBlocks[i])); err != nil {
				return nil, err
			}
		}

		serviceSubnets := make([]meta.Prefix, len(serviceBlocks))

		for i := range serviceSubnets {
			if err := serviceSubnets[i].UnmarshalText([]byte(serviceBlocks[i])); err != nil {
				return nil, err
			}
		}

		networkConfig := k8s.NewKubeNetworkConfigV1Alpha1()
		networkConfig.NetworkPodSubnets = podSubnets
		networkConfig.NetworkServiceSubnets = serviceSubnets

		return patchFromDocument(networkConfig)
	}

	return patchFromV1Alpha1(map[string]any{
		"cluster": map[string]any{
			"network": map[string]any{
				"podSubnets":     podBlocks,
				"serviceSubnets": serviceBlocks,
			},
		},
	})

}

func patchFromDocument(doc configconfig.Document) ([]byte, error) {
	ctr, err := container.New(doc)
	if err != nil {
		return nil, err
	}

	return ctr.EncodeBytes(encoder.WithComments(encoder.CommentsDisabled))
}

func patchFromV1Alpha1(doc any) ([]byte, error) {
	return yaml.Marshal(doc)
}
