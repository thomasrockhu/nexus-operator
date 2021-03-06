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

package networking

import (
	ctx "context"
	"fmt"
	"github.com/m88i/nexus-operator/pkg/apis/apps/v1alpha1"
	"github.com/m88i/nexus-operator/pkg/controller/nexus/resource/deployment"
	"github.com/m88i/nexus-operator/pkg/test"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"reflect"
	"testing"
)

var nodePortNexus = &v1alpha1.Nexus{
	ObjectMeta: metav1.ObjectMeta{Namespace: "test", Name: "nexus"},
	Spec: v1alpha1.NexusSpec{
		Networking: v1alpha1.NexusNetworking{Expose: true, NodePort: 31031, ExposeAs: v1alpha1.NodePortExposeType},
	},
}

func TestManager_IngressAvailable(t *testing.T) {
	client := test.NewFakeClientBuilder().Build()
	nexus := v1alpha1.Nexus{}
	mgr, err := NewManager(nexus, client, client)
	assert.Nil(t, err)

	assert.Equal(t, mgr.ingressAvailable, mgr.IngressAvailable())
}

func TestManager_RouteAvailable(t *testing.T) {
	client := test.NewFakeClientBuilder().Build()
	nexus := v1alpha1.Nexus{}
	mgr, err := NewManager(nexus, client, client)
	assert.Nil(t, err)

	assert.Equal(t, mgr.ingressAvailable, mgr.IngressAvailable())
}

func TestNewManager(t *testing.T) {
	k8sClient := test.NewFakeClientBuilder().Build()
	k8sClientWithIngress := test.NewFakeClientBuilder().WithIngress().Build()
	ocpClient := test.NewFakeClientBuilder().OnOpenshift().Build()
	//default-setting logic is tested elsewhere
	//so here we just check if the resulting manager took in the arguments correctly
	tests := []struct {
		name       string
		want       *Manager
		wantClient *test.FakeClient
	}{
		{
			"On Kubernetes with Ingresses available",
			&Manager{
				nexus:            nodePortNexus,
				client:           test.NewFakeClientBuilder().WithIngress().Build(),
				routeAvailable:   false,
				ingressAvailable: true,
				ocp:              false,
			},
			k8sClientWithIngress,
		},
		{
			"On Kubernetes without Ingresses",
			&Manager{
				nexus:            nodePortNexus,
				client:           test.NewFakeClientBuilder().Build(),
				routeAvailable:   false,
				ingressAvailable: false,
				ocp:              false,
			},
			k8sClient,
		},
		{
			"On Openshift",
			&Manager{
				nexus:            nodePortNexus,
				client:           test.NewFakeClientBuilder().OnOpenshift().Build(),
				routeAvailable:   true,
				ingressAvailable: false,
				ocp:              true,
			},
			ocpClient,
		},
	}

	for _, tt := range tests {
		got, err := NewManager(*nodePortNexus, tt.wantClient, tt.wantClient)
		assert.NoError(t, err)
		assert.NotNil(t, got.client)
		assert.NotNil(t, got.nexus)
		assert.Equal(t, tt.want.routeAvailable, got.routeAvailable)
		assert.Equal(t, tt.want.ingressAvailable, got.ingressAvailable)
		assert.Equal(t, tt.want.ocp, got.ocp)
	}

	// simulate discovery 500 response, expect error
	mockErrorMsg := "mock 500"
	k8sClient.SetMockErrorForOneRequest(errors.NewInternalError(fmt.Errorf(mockErrorMsg)))
	mgr, err := NewManager(*nodePortNexus, k8sClient, k8sClient)
	assert.Nil(t, mgr)
	assert.Contains(t, err.Error(), mockErrorMsg)

	// test with an invalid Nexus, expect error
	invalidNexus := v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com", TLS: v1alpha1.NexusNetworkingTLS{Mandatory: true}}}}
	mgr, err = NewManager(invalidNexus, k8sClientWithIngress, k8sClientWithIngress)
	assert.Nil(t, mgr)
	assert.NotNil(t, err)
}

func TestNewManager_setDefaults(t *testing.T) {
	tests := []struct {
		name             string
		ocp              bool
		routeAvailable   bool
		ingressAvailable bool
		input            v1alpha1.Nexus
		want             v1alpha1.Nexus
	}{
		{
			"'spec.networking.expose' left blank or set to false",
			false, // unimportant
			false, // unimportant
			false, // unimportant
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: false}}},
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: false}}},
		},
		{
			"'spec.networking.exposeAs' left blank on OCP",
			true,
			true,
			false, // unimportant
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true}}},
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType}}},
		},
		{
			"'spec.networking.exposeAs' left blank on K8s",
			false,
			false,
			true,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true}}},
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType}}},
		},
		{
			"'spec.networking.exposeAs' left blank on K8s, but Ingress unavailable",
			false,
			false,
			false,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true}}},
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.NodePortExposeType}}},
		},
	}

	for _, tt := range tests {
		manager := &Manager{
			nexus:            &tt.input,
			routeAvailable:   tt.routeAvailable,
			ingressAvailable: tt.ingressAvailable,
			ocp:              tt.ocp,
		}
		manager.setDefaults()
		if !reflect.DeepEqual(*manager.nexus, tt.want) {
			t.Errorf("TestManager_setDefaults() - %s\nWant: %v\tGot: %v", tt.name, tt.want, *manager.nexus)
		}
	}
}

func TestManager_validate(t *testing.T) {
	tests := []struct {
		name             string
		ocp              bool
		routeAvailable   bool
		ingressAvailable bool
		input            v1alpha1.Nexus
		wantError        bool
	}{
		{
			"'spec.networking.expose' left blank or set to false",
			false, // unimportant
			false, // unimportant
			false, // unimportant
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: false}}},
			false,
		},
		{
			"Valid Nexus with Ingress on K8s",
			false,
			false,
			true,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com"}}},
			false,
		},
		{
			"Valid Nexus with Ingress and TLS secret on K8s",
			false,
			false,
			true,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com", TLS: v1alpha1.NexusNetworkingTLS{SecretName: "tt-secret"}}}},
			false,
		},
		{
			"Valid Nexus with Ingress on K8s, but Ingress unavailable (Kubernetes < 1.14)",
			false,
			false,
			false,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com"}}},
			true,
		},
		{
			"Invalid Nexus with Ingress on K8s and no 'spec.networking.host'",
			false,
			false,
			true,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType}}},
			true,
		},
		{
			"Invalid Nexus with Ingress on K8s and 'spec.networking.mandatory' set to 'true'",
			false,
			false,
			true,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com", TLS: v1alpha1.NexusNetworkingTLS{Mandatory: true}}}},
			true,
		},
		{
			"Invalid Nexus with Route on K8s",
			false,
			false,
			true,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType}}},
			true,
		},
		{
			"Valid Nexus with Route on OCP",
			true,
			true,
			false,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType}}},
			false,
		},
		{
			"Valid Nexus with Route on OCP with mandatory TLS",
			true,
			true,
			false,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType, TLS: v1alpha1.NexusNetworkingTLS{Mandatory: true}}}},
			false,
		},
		{
			"Invalid Nexus with Route on OCP, but using secret name",
			true,
			true,
			false,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.RouteExposeType, TLS: v1alpha1.NexusNetworkingTLS{SecretName: "test-secret"}}}},
			true,
		},
		{
			"Invalid Nexus with Ingress on OCP",
			true,
			true,
			false,
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.IngressExposeType, Host: "example.com"}}},
			true,
		},
		{
			"Valid Nexus with Node Port",
			false, // unimportant
			false, // unimportant
			false, // unimportant
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.NodePortExposeType, NodePort: 8080}}},
			false,
		},
		{
			"Invalid Nexus with Node Port and no port informed",
			false, // unimportant
			false, // unimportant
			false, // unimportant
			v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: true, ExposeAs: v1alpha1.NodePortExposeType}}},
			true,
		},
	}

	for _, tt := range tests {
		manager := &Manager{
			nexus:            &tt.input,
			routeAvailable:   tt.routeAvailable,
			ingressAvailable: tt.ingressAvailable,
			ocp:              tt.ocp,
		}
		if err := manager.validate(); (err != nil) != tt.wantError {
			t.Errorf("TestManager_validate() - %s\nWantError: %v\tError: %v", tt.name, tt.wantError, err)
		}
	}
}

func TestManager_GetRequiredResources(t *testing.T) {
	// correctness of the generated resources is tested elsewhere
	// here we just want to check if they have been created and returned
	// first, let's test a Nexus which does not expose
	mgr := &Manager{
		nexus:  &v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{Expose: false}}},
		client: test.NewFakeClientBuilder().Build(),
	}
	resources, err := mgr.GetRequiredResources()
	assert.Nil(t, resources)
	assert.Nil(t, err)

	// now, let's use a route
	mgr = &Manager{
		nexus:          routeNexus,
		client:         test.NewFakeClientBuilder().OnOpenshift().Build(),
		routeAvailable: true,
		ocp:            true,
	}
	resources, err = mgr.GetRequiredResources()
	assert.Nil(t, err)
	assert.Len(t, resources, 1)
	assert.True(t, test.ContainsType(resources, reflect.TypeOf(&routev1.Route{})))

	// still a route, but in a cluster without routes
	mgr = &Manager{
		nexus:  routeNexus,
		client: test.NewFakeClientBuilder().Build(),
	}
	resources, err = mgr.GetRequiredResources()
	assert.Nil(t, resources)
	assert.EqualError(t, err, fmt.Sprintf(resUnavailableFormat, "routes"))

	// now an ingress
	mgr = &Manager{
		nexus:            ingressNexus,
		client:           test.NewFakeClientBuilder().WithIngress().Build(),
		ingressAvailable: true,
	}
	resources, err = mgr.GetRequiredResources()
	assert.Nil(t, err)
	assert.Len(t, resources, 1)
	assert.True(t, test.ContainsType(resources, reflect.TypeOf(&networkingv1beta1.Ingress{})))

	// still an ingress, but in a cluster without ingresses
	mgr = &Manager{
		nexus:  ingressNexus,
		client: test.NewFakeClientBuilder().Build(),
	}
	resources, err = mgr.GetRequiredResources()
	assert.Nil(t, resources)
	assert.EqualError(t, err, fmt.Sprintf(resUnavailableFormat, "ingresses"))
}

func TestManager_createRoute(t *testing.T) {
	mgr := &Manager{nexus: &v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{TLS: v1alpha1.NexusNetworkingTLS{}}}}}

	// first without TLS
	route := mgr.createRoute()
	assert.Nil(t, route.Spec.TLS)

	// now with TLS
	mgr.nexus.Spec.Networking.TLS.Mandatory = true
	route = mgr.createRoute()
	assert.NotNil(t, route.Spec.TLS)
}

func TestManager_createIngress(t *testing.T) {
	mgr := &Manager{nexus: &v1alpha1.Nexus{Spec: v1alpha1.NexusSpec{Networking: v1alpha1.NexusNetworking{TLS: v1alpha1.NexusNetworkingTLS{}}}}}

	// first without TLS
	ingress := mgr.createIngress()
	assert.Empty(t, ingress.Spec.TLS)

	// now with TLS
	mgr.nexus.Spec.Networking.TLS.SecretName = "test-secret"
	ingress = mgr.createIngress()
	assert.Len(t, ingress.Spec.TLS, 1)
}

func TestManager_GetDeployedResources(t *testing.T) {
	// first with no deployed resources
	fakeClient := test.NewFakeClientBuilder().WithIngress().OnOpenshift().Build()
	mgr := &Manager{
		nexus:            nodePortNexus,
		client:           fakeClient,
		ingressAvailable: true,
		routeAvailable:   true,
		ocp:              true,
	}
	resources, err := mgr.GetDeployedResources()
	assert.Nil(t, resources)
	assert.Len(t, resources, 0)
	assert.NoError(t, err)

	// now with deployed resources
	route := &routev1.Route{ObjectMeta: metav1.ObjectMeta{Name: mgr.nexus.Name, Namespace: mgr.nexus.Namespace}}
	assert.NoError(t, mgr.client.Create(ctx.TODO(), route))

	ingress := &networkingv1beta1.Ingress{ObjectMeta: metav1.ObjectMeta{Name: mgr.nexus.Name, Namespace: mgr.nexus.Namespace}}
	assert.NoError(t, mgr.client.Create(ctx.TODO(), ingress))

	resources, err = mgr.GetDeployedResources()
	assert.Nil(t, err)
	assert.Len(t, resources, 2)
	assert.True(t, test.ContainsType(resources, reflect.TypeOf(&routev1.Route{})))
	assert.True(t, test.ContainsType(resources, reflect.TypeOf(&networkingv1beta1.Ingress{})))

	// make the client return a mocked 500 response to test errors other than NotFound
	mockErrorMsg := "mock 500"
	fakeClient.SetMockErrorForOneRequest(errors.NewInternalError(fmt.Errorf(mockErrorMsg)))
	resources, err = mgr.GetDeployedResources()
	assert.Nil(t, resources)
	assert.Contains(t, err.Error(), mockErrorMsg)
}

func TestManager_getDeployedRoute(t *testing.T) {
	mgr := &Manager{
		nexus:  routeNexus,
		client: test.NewFakeClientBuilder().OnOpenshift().Build(),
	}

	// first, test without creating the route
	route, err := mgr.getDeployedRoute()
	assert.Nil(t, route)
	assert.True(t, errors.IsNotFound(err))

	// now test after creating the route
	route = &routev1.Route{ObjectMeta: metav1.ObjectMeta{Name: mgr.nexus.Name, Namespace: mgr.nexus.Namespace}}
	assert.NoError(t, mgr.client.Create(ctx.TODO(), route))
	route, err = mgr.getDeployedRoute()
	assert.NotNil(t, route)
	assert.NoError(t, err)
}

func TestManager_getDeployedIngress(t *testing.T) {
	mgr := &Manager{
		nexus:  ingressNexus,
		client: test.NewFakeClientBuilder().WithIngress().Build(),
	}

	// first, test without creating the ingress
	ingress, err := mgr.getDeployedIngress()
	assert.Nil(t, ingress)
	assert.True(t, errors.IsNotFound(err))

	// now test after creating the ingress
	ingress = &networkingv1beta1.Ingress{ObjectMeta: metav1.ObjectMeta{Name: mgr.nexus.Name, Namespace: mgr.nexus.Namespace}}
	assert.NoError(t, mgr.client.Create(ctx.TODO(), ingress))
	ingress, err = mgr.getDeployedIngress()
	assert.NotNil(t, ingress)
	assert.NoError(t, err)
}

func TestManager_GetCustomComparator(t *testing.T) {
	// the nexus and the client should have no effect on the
	// comparator functions offered by the manager
	mgr := &Manager{}

	// there is no custom comparator function for routes
	routeComp := mgr.GetCustomComparator(reflect.TypeOf(&routev1.Route{}))
	assert.Nil(t, routeComp)
	// there is a custom comparator function for ingresses
	ingressComp := mgr.GetCustomComparator(reflect.TypeOf(&networkingv1beta1.Ingress{}))
	assert.NotNil(t, ingressComp)
}

func TestManager_GetCustomComparators(t *testing.T) {
	// the nexus and the client should have no effect on the
	// comparator functions offered by the manager
	mgr := &Manager{}

	// there is one custom comparator (ingresses)
	comparators := mgr.GetCustomComparators()
	assert.Len(t, comparators, 1)
}

func TestIngressEqual(t *testing.T) {
	// base ingress which will be used in all comparisons
	baseIngress := &networkingv1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{Name: "test", Namespace: "test", UID: "base UID"},
		Spec: networkingv1beta1.IngressSpec{
			Rules: []networkingv1beta1.IngressRule{
				{
					Host: "test.example.com",
					IngressRuleValue: networkingv1beta1.IngressRuleValue{
						HTTP: &networkingv1beta1.HTTPIngressRuleValue{
							Paths: []networkingv1beta1.HTTPIngressPath{
								{
									Path: ingressBasePath,
									Backend: networkingv1beta1.IngressBackend{
										ServiceName: "test",
										ServicePort: intstr.FromInt(deployment.NexusServicePort),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	testCases := []struct {
		name            string
		modifiedIngress *networkingv1beta1.Ingress
		wantEqual       bool
	}{
		{
			"Two ingresses that are equal where it matters and different where it doesn't",
			func() *networkingv1beta1.Ingress {
				ingress := baseIngress.DeepCopy()
				// simulates a field that would be different in runtime
				ingress.UID = "a different UID"
				return ingress
			}(),
			true,
		},
		{
			"All equal except name",
			func() *networkingv1beta1.Ingress {
				ingress := baseIngress.DeepCopy()
				ingress.Name = "a different name"
				return ingress
			}(),
			false,
		},
		{
			"All equal except namespace",
			func() *networkingv1beta1.Ingress {
				ingress := baseIngress.DeepCopy()
				ingress.Namespace = "a different namespace"
				return ingress
			}(),
			false,
		},
		{
			"All equal except the service name",
			func() *networkingv1beta1.Ingress {
				ingress := baseIngress.DeepCopy()
				ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.ServiceName = "a different service"
				return ingress
			}(),
			false,
		},
	}

	for _, testCase := range testCases {
		if ingressEqual(baseIngress, testCase.modifiedIngress) != testCase.wantEqual {
			assert.Failf(t, "%s\nbase: %+v\nmodified: %+v\nwantedEqual: %v", testCase.name, baseIngress, testCase.modifiedIngress, testCase.wantEqual)
		}
	}
}
