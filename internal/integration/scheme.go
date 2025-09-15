// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package integration

import (
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	capdv1 "sigs.k8s.io/cluster-api/test/infrastructure/docker/api/v1beta2"
)

func init() {
	utilruntime.Must(capdv1.AddToScheme(scheme.Scheme))
}
