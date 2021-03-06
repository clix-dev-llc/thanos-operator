// Copyright 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"github.com/banzaicloud/operator-tools/pkg/reconciler"
	netv1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (e *storeInstance) ingressGRPC() (runtime.Object, reconciler.DesiredState, error) {
	if e.StoreEndpoint != nil && e.StoreEndpoint.Spec.Ingress != nil {
		endpointIngress := e.StoreEndpoint.Spec.Ingress
		pathType := netv1.PathTypeImplementationSpecific
		ingress := &netv1.Ingress{
			ObjectMeta: e.StoreEndpoint.Spec.MetaOverrides.Merge(e.getMeta()),
			Spec: netv1.IngressSpec{
				Rules: []netv1.IngressRule{
					{
						Host: endpointIngress.Host,
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path:     endpointIngress.Path,
										PathType: &pathType,
										Backend: netv1.IngressBackend{
											ServiceName: e.GetName(),
											ServicePort: intstr.FromString("grpc"),
										},
									},
								},
							},
						},
					},
				},
			},
		}
		if endpointIngress.Certificate != "" {
			ingress.Spec.TLS = []netv1.IngressTLS{
				{
					Hosts:      []string{endpointIngress.Host},
					SecretName: endpointIngress.Certificate,
				},
			}
		}
		return ingress, reconciler.StatePresent, nil
	}
	delete := &netv1.Ingress{
		ObjectMeta: e.getMeta(),
	}
	return delete, reconciler.StateAbsent, nil
}
