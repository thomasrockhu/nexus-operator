// +build !ignore_autogenerated

//     Copyright 2019 Nexus Operator and/or its authors
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
// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"./pkg/apis/apps/v1alpha1.Nexus":            schema_pkg_apis_apps_v1alpha1_Nexus(ref),
		"./pkg/apis/apps/v1alpha1.NexusPersistence": schema_pkg_apis_apps_v1alpha1_NexusPersistence(ref),
		"./pkg/apis/apps/v1alpha1.NexusProbe":       schema_pkg_apis_apps_v1alpha1_NexusProbe(ref),
		"./pkg/apis/apps/v1alpha1.NexusSpec":        schema_pkg_apis_apps_v1alpha1_NexusSpec(ref),
		"./pkg/apis/apps/v1alpha1.NexusStatus":      schema_pkg_apis_apps_v1alpha1_NexusStatus(ref),
	}
}

func schema_pkg_apis_apps_v1alpha1_Nexus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "Nexus custom resource to deploy the Nexus Server",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/apps/v1alpha1.NexusSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("./pkg/apis/apps/v1alpha1.NexusStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"./pkg/apis/apps/v1alpha1.NexusSpec", "./pkg/apis/apps/v1alpha1.NexusStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_apps_v1alpha1_NexusPersistence(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NexusPersistence is the structure for the data persistent",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"persistent": {
						SchemaProps: spec.SchemaProps{
							Description: "Flag to indicate if this instance will be persistent or not",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"volumeSize": {
						SchemaProps: spec.SchemaProps{
							Description: "If persistent, the size of the Volume. Defaults: 10Gi",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"storageClass": {
						SchemaProps: spec.SchemaProps{
							Description: "StorageClass used by the managed PVC.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
				Required: []string{"persistent"},
			},
		},
	}
}

func schema_pkg_apis_apps_v1alpha1_NexusProbe(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NexusProbe describes a health check to be performed against a container to determine whether it is alive or ready to receive traffic.",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"initialDelaySeconds": {
						SchemaProps: spec.SchemaProps{
							Description: "Number of seconds after the container has started before probes are initiated. Defaults to 240 seconds. Minimum value is 0.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"timeoutSeconds": {
						SchemaProps: spec.SchemaProps{
							Description: "Number of seconds after which the probe times out. Defaults to 15 seconds. Minimum value is 1.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"periodSeconds": {
						SchemaProps: spec.SchemaProps{
							Description: "How often (in seconds) to perform the probe. Defaults to 10 seconds. Minimum value is 1.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"successThreshold": {
						SchemaProps: spec.SchemaProps{
							Description: "Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"failureThreshold": {
						SchemaProps: spec.SchemaProps{
							Description: "Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
				},
			},
		},
	}
}

func schema_pkg_apis_apps_v1alpha1_NexusSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NexusSpec defines the desired state of Nexus",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"replicas": {
						SchemaProps: spec.SchemaProps{
							Description: "Number of pod replicas desired. Defaults to 0.",
							Type:        []string{"integer"},
							Format:      "int32",
						},
					},
					"image": {
						SchemaProps: spec.SchemaProps{
							Description: "Full image tag name for this specific deployment. Will be ignored if `spec.useRedHatImage` is set to `true`. Default: docker.io/sonatype/nexus3:latest",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"imagePullPolicy": {
						SchemaProps: spec.SchemaProps{
							Description: "The image pull policy for the Nexus image. If left blank behavior will be determined by the image tag (`Always` if \"latest\" and `IfNotPresent` otherwise). Possible values: `Always`, `IfNotPresent` or `Never`.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"resources": {
						SchemaProps: spec.SchemaProps{
							Description: "Defined Resources for the Nexus instance",
							Ref:         ref("k8s.io/api/core/v1.ResourceRequirements"),
						},
					},
					"persistence": {
						SchemaProps: spec.SchemaProps{
							Description: "Persistence definition",
							Ref:         ref("./pkg/apis/apps/v1alpha1.NexusPersistence"),
						},
					},
					"useRedHatImage": {
						SchemaProps: spec.SchemaProps{
							Description: "If you have access to Red Hat Container Catalog, set this to `true` to use the certified image provided by Sonatype Defaults to `false`",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"generateRandomAdminPassword": {
						SchemaProps: spec.SchemaProps{
							Description: "GenerateRandomAdminPassword enables the random password generation. Defaults to `false`: the default password for a newly created instance is 'admin123', which should be changed in the first login. If set to `true`, you must use the automatically generated 'admin' password, stored in the container's file system at `/nexus-data/admin.password`.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"networking": {
						SchemaProps: spec.SchemaProps{
							Description: "Networking definition",
							Ref:         ref("./pkg/apis/apps/v1alpha1.NexusNetworking"),
						},
					},
					"serviceAccountName": {
						SchemaProps: spec.SchemaProps{
							Description: "ServiceAccountName is the name of the ServiceAccount used to run the Pods. If left blank, a default ServiceAccount is created with the same name as the Nexus CR (`metadata.name`).",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"livenessProbe": {
						SchemaProps: spec.SchemaProps{
							Description: "LivenessProbe describes how the Nexus container liveness probe should work",
							Ref:         ref("./pkg/apis/apps/v1alpha1.NexusProbe"),
						},
					},
					"readinessProbe": {
						SchemaProps: spec.SchemaProps{
							Description: "ReadinessProbe describes how the Nexus container readiness probe should work",
							Ref:         ref("./pkg/apis/apps/v1alpha1.NexusProbe"),
						},
					},
				},
				Required: []string{"replicas", "persistence", "useRedHatImage"},
			},
		},
		Dependencies: []string{
			"./pkg/apis/apps/v1alpha1.NexusNetworking", "./pkg/apis/apps/v1alpha1.NexusPersistence", "./pkg/apis/apps/v1alpha1.NexusProbe", "k8s.io/api/core/v1.ResourceRequirements"},
	}
}

func schema_pkg_apis_apps_v1alpha1_NexusStatus(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "NexusStatus defines the observed state of Nexus",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"deploymentStatus": {
						SchemaProps: spec.SchemaProps{
							Description: "Condition status for the Nexus deployment",
							Ref:         ref("k8s.io/api/apps/v1.DeploymentStatus"),
						},
					},
					"nexusStatus": {
						SchemaProps: spec.SchemaProps{
							Description: "Will be \"OK\" when all objects are created successfully",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"nexusRoute": {
						SchemaProps: spec.SchemaProps{
							Description: "Route for external service access",
							Type:        []string{"string"},
							Format:      "",
						},
					},
				},
			},
		},
		Dependencies: []string{
			"k8s.io/api/apps/v1.DeploymentStatus"},
	}
}
