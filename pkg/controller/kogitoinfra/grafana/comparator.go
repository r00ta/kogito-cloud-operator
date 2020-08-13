// Copyright 2019 Red Hat, Inc. and/or its affiliates
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

package grafana

import (
	"reflect"

	grafanav1 "github.com/integr8ly/grafana-operator/pkg/apis/integreatly/v1alpha1"

	"github.com/RHsyseng/operator-utils/pkg/resource"
	"github.com/kiegroup/kogito-cloud-operator/pkg/framework"
)

// GetComparators gets the comparator for Grafana resources
func GetComparators() []framework.Comparator {
	return []framework.Comparator{createGrafanaComparator()}
}

func createGrafanaComparator() framework.Comparator {
	return framework.Comparator{
		ResourceType: reflect.TypeOf(grafanav1.Grafana{}),
		CompFunc: func(deployed resource.KubernetesResource, requested resource.KubernetesResource) bool {
			grafanaDep := deployed.(*grafanav1.Grafana)
			grafanaReq := requested.(*grafanav1.Grafana).DeepCopy()
			// we just care for the instance name, other attributes can be changed at will by the user
			return reflect.DeepEqual(grafanaDep.Name, grafanaReq.Name)
		},
	}
}
