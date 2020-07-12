//     Copyright 2020 Nexus Operator and/or its authors
//
//     This file is part of Nexus Operator.
//
//     Nexus Operator is free software: you can redistribute it and/or modify
//     it under the terms of the GNU General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.
//
//     Nexus Operator is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU General Public License for more details.
//
//     You should have received a copy of the GNU General Public License
//     along with Nexus Operator.  If not, see <https://www.gnu.org/licenses/>.

package server

import (
	"net/url"
	"testing"

	"github.com/m88i/nexus-operator/pkg/apis/apps/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/controller/nexus/resource/meta"
	"github.com/m88i/nexus-operator/pkg/test"
	"github.com/stretchr/testify/assert"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func Test_server_getNexusEndpoint(t *testing.T) {
	nexus := &v1alpha1.Nexus{
		Spec:       v1alpha1.NexusSpec{},
		ObjectMeta: v1.ObjectMeta{Name: "nexus3", Namespace: t.Name()},
	}
	svc := &corev1.Service{
		ObjectMeta: meta.DefaultObjectMeta(nexus),
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Protocol: corev1.ProtocolTCP,
					Port:     8081,
					TargetPort: intstr.IntOrString{
						IntVal: 8081,
					},
				},
			},
			Selector:        meta.GenerateLabels(nexus),
			SessionAffinity: corev1.ServiceAffinityNone,
		},
	}
	cli := test.NewFakeClientBuilder(nexus, svc).Build()
	s := server{
		nexus:     nexus,
		k8sclient: cli,
	}
	URL, err := s.getNexusEndpoint()
	assert.NoError(t, err)
	assert.NotEmpty(t, URL)
	assert.Contains(t, URL, nexus.Name)
	assert.Contains(t, URL, "8081")
	_, err = url.Parse(URL)
	assert.NoError(t, err)
}

func Test_server_isServerReady(t *testing.T) {
	nexus := &v1alpha1.Nexus{
		Spec:       v1alpha1.NexusSpec{},
		ObjectMeta: v1.ObjectMeta{Name: "nexus3", Namespace: t.Name()},
		Status: v1alpha1.NexusStatus{
			DeploymentStatus: appv1.DeploymentStatus{
				AvailableReplicas: 1,
			},
		},
	}
	s := server{nexus: nexus}
	assert.True(t, s.isServerReady())
}

// TODO: add fakes/mock to AICURA to make it easy to mock the server when using the API

func Test_server_getOrCreateOperatorCredentials(t *testing.T) {
	nexus := &v1alpha1.Nexus{
		Spec:       v1alpha1.NexusSpec{},
		ObjectMeta: v1.ObjectMeta{Name: "nexus3", Namespace: t.Name()},
	}
	cli := test.NewFakeClientBuilder(nexus).Build()
	s := server{
		nexus:     nexus,
		k8sclient: cli,
	}
	secret, err := s.getOrCreateOperatorCredentials()
	assert.NoError(t, err)
	assert.NotNil(t, secret)
	assert.NotEmpty(t, secret.StringData)
	assert.Equal(t, operatorUsername, secret.StringData[secretKeyUsername])
	assert.NotEmpty(t, secret.StringData[secretKeyPassword])
}
