// Copyright 2020 Red Hat, Inc. and/or its affiliates
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

package kogitosupportingservice

import (
	"github.com/kiegroup/kogito-cloud-operator/pkg/apis/app/v1alpha1"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/kubernetes"
	"github.com/kiegroup/kogito-cloud-operator/pkg/client/meta"
	"github.com/kiegroup/kogito-cloud-operator/pkg/infrastructure"
	"github.com/kiegroup/kogito-cloud-operator/pkg/test"
	imagev1 "github.com/openshift/api/image/v1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestReconcileKogitoSupportingServiceMgmtConsole_Reconcile(t *testing.T) {
	replicas := int32(1)
	instance := &v1alpha1.KogitoSupportingService{
		ObjectMeta: v1.ObjectMeta{Name: infrastructure.DefaultMgmtConsoleName, Namespace: t.Name()},
		Spec: v1alpha1.KogitoSupportingServiceSpec{
			ServiceType:       v1alpha1.MgmtConsole,
			KogitoServiceSpec: v1alpha1.KogitoServiceSpec{Replicas: &replicas},
		},
	}
	cli := test.NewFakeClientBuilder().AddK8sObjects(instance).OnOpenShift().Build()

	r := &MgmtConsoleSupportingServiceResource{}
	// first reconciliation
	requeueAfter, err := r.Reconcile(cli, instance, meta.GetRegisteredSchema())
	assert.NoError(t, err)
	assert.True(t, requeueAfter == 0)
	// second time
	requeueAfter, err = r.Reconcile(cli, instance, meta.GetRegisteredSchema())
	assert.NoError(t, err)
	assert.True(t, requeueAfter == 0)

	_, err = kubernetes.ResourceC(cli).Fetch(instance)
	assert.NoError(t, err)
	assert.NotNil(t, instance.Status)
	assert.Len(t, instance.Status.Conditions, 1)
}

// see: https://issues.redhat.com/browse/KOGITO-2535
func TestReconcileKogitoSupportingServiceMgmtConsole_CustomImage(t *testing.T) {
	replicas := int32(1)
	instance := &v1alpha1.KogitoSupportingService{
		ObjectMeta: v1.ObjectMeta{Name: infrastructure.DefaultMgmtConsoleName, Namespace: t.Name()},
		Spec: v1alpha1.KogitoSupportingServiceSpec{
			ServiceType: v1alpha1.MgmtConsole,
			KogitoServiceSpec: v1alpha1.KogitoServiceSpec{
				Replicas: &replicas,
				Image:    "quay.io/mynamespace/super-mgmt-console:0.1.3",
			},
		},
	}
	cli := test.NewFakeClientBuilder().AddK8sObjects(instance).OnOpenShift().Build()

	r := &MgmtConsoleSupportingServiceResource{}
	requeueAfter, err := r.Reconcile(cli, instance, meta.GetRegisteredSchema())
	assert.NoError(t, err)
	assert.True(t, requeueAfter == 0)
	// image stream
	is := imagev1.ImageStream{
		ObjectMeta: v1.ObjectMeta{Name: infrastructure.DefaultMgmtConsoleImageName, Namespace: instance.Namespace},
	}
	exists, err := kubernetes.ResourceC(cli).Fetch(&is)
	assert.True(t, exists)
	assert.NoError(t, err)
	assert.Len(t, is.Spec.Tags, 1)
	assert.Equal(t, "0.1.3", is.Spec.Tags[0].Name)
	assert.Equal(t, "quay.io/mynamespace/super-mgmt-console:0.1.3", is.Spec.Tags[0].From.Name)

}
