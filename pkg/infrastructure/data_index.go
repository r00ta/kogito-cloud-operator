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

package infrastructure

import (
	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/kubernetes"
	appsv1 "k8s.io/api/apps/v1"
)

const (
	// DefaultDataIndexImageName is just the image name for the Data Index Service
	DefaultDataIndexImageName = "kogito-data-index-infinispan"
	// DefaultDataIndexName is the default name for the Data Index instance service
	DefaultDataIndexName = "data-index"
	// Data index HTTP URL env
	dataIndexHTTPRouteEnv = "KOGITO_DATAINDEX_HTTP_URL"
	// Data index WS URL env
	dataIndexWSRouteEnv = "KOGITO_DATAINDEX_WS_URL"
)

// InjectDataIndexURLIntoKogitoRuntimeServices will query for every KogitoRuntime in the given namespace to inject the Data Index route to each one
// Won't trigger an update if the KogitoRuntime already has the route set to avoid unnecessary reconciliation triggers
func InjectDataIndexURLIntoKogitoRuntimeServices(client *client.Client, namespace string) error {
	log.Debugf("Injecting Data-Index Route in kogito apps")
	return injectSupportingServiceURLIntoKogitoRuntime(client, namespace, dataIndexHTTPRouteEnv, dataIndexWSRouteEnv, v1alpha1.DataIndex)
}

// InjectDataIndexURLIntoDeployment will inject data-index route URL in to kogito runtime deployment env var
func InjectDataIndexURLIntoDeployment(client *client.Client, namespace string, deployment *appsv1.Deployment) error {
	log.Debugf("Injecting Data-Index URL in kogito Runtime deployment")
	return injectSupportingServiceURLIntoDeployment(client, namespace, dataIndexHTTPRouteEnv, dataIndexWSRouteEnv, deployment, v1alpha1.DataIndex)
}

// InjectDataIndexURLIntoSupportingService will query for Supporting service deployment in the given namespace to inject the Data Index route to each one
// Won't trigger an update if the SupportingService already has the route set to avoid unnecessary reconciliation triggers
func InjectDataIndexURLIntoSupportingService(client *client.Client, namespace string, serviceTypes ...v1alpha1.ServiceType) error {
	for _, serviceType := range serviceTypes {
		log.Debugf("Injecting Data-Index Route in %s", serviceType)
		deployment, err := getSupportingServiceDeployment(namespace, client, serviceType)
		if err != nil {
			return err
		}
		if deployment == nil {
			log.Debugf("No deployment found for %s, skipping to inject %s URL into %s", serviceType, v1alpha1.DataIndex, serviceType)
			return nil
		}

		log.Debugf("Querying %s route to inject into %s", v1alpha1.DataIndex, serviceType)
		serviceEndpoints, err := getServiceEndpoints(client, namespace, dataIndexHTTPRouteEnv, dataIndexWSRouteEnv, v1alpha1.DataIndex)
		if err != nil {
			return err
		}
		if serviceEndpoints != nil {
			log.Debugf("The %s route is '%s'", v1alpha1.DataIndex, serviceEndpoints.HTTPRouteURI)

			updateHTTP, updateWS := updateServiceEndpointIntoDeploymentEnv(deployment, serviceEndpoints)
			// update only once
			if updateWS || updateHTTP {
				if err := kubernetes.ResourceC(client).Update(deployment); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
