/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/ptr"
	v1 "knative.dev/serving/pkg/apis/serving/v1"
)

func (r *Revision) SetDefaults(ctx context.Context) {
	r.Spec.SetDefaults(apis.WithinSpec(ctx))
}

// SetDefaults implements apis.Defaultable
func (rts *RevisionTemplateSpec) SetDefaults(ctx context.Context) {
	rts.Spec.SetDefaults(apis.WithinSpec(ctx))
}

func (rs *RevisionSpec) SetDefaults(ctx context.Context) {
	if v1.IsUpgradeViaDefaulting(ctx) {
		v1 := v1.RevisionSpec{}
		if rs.ConvertTo(ctx, &v1) == nil {
			alpha := RevisionSpec{}
			if alpha.ConvertFrom(ctx, v1) == nil {
				*rs = alpha
			}
		}
	}

	// When ConcurrencyModel is specified but ContainerConcurrency
	// is not (`nil`), use the ConcurrencyModel value.
	if rs.DeprecatedConcurrencyModel == DeprecatedRevisionRequestConcurrencyModelSingle && rs.ContainerConcurrency == nil {
		rs.ContainerConcurrency = ptr.Int64(1)
	}

	// When the PodSpec has no containers, move the single Container
	// into the PodSpec for the scope of defaulting and then move
	// it back as we return.
	if len(rs.Containers) == 0 {
		if rs.DeprecatedContainer == nil {
			rs.DeprecatedContainer = &corev1.Container{}
		}
		rs.Containers = []corev1.Container{*rs.DeprecatedContainer}
		defer func() {
			rs.DeprecatedContainer = &rs.Containers[0]
			rs.Containers = nil
		}()
	}
	rs.RevisionSpec.SetDefaults(ctx)
}
