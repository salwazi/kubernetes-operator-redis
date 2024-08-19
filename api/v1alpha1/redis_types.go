/*
Copyright 2024.

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
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisSpec defines the desired state of Redis
type RedisSpec struct {

	// Image is the Redis Docker image
	Image string `json:"image"`

	// Version is the version of Redis to deploy
	Version string `json:"version"`

	// Storage defines the storage requirements for Redis
	Storage RedisStorage `json:"storage"`

	// Replicas is the number of replicas for each Redis node
	Replicas int32 `json:"replicas"`

	// SecretName is the name of the Kubernetes Secret object that stores the Redis password
	SecretName string `json:"secretName,omitempty"`

	// Resources defines the CPU and memory resource requirements
	Resources RedisResources `json:"resources"`
}

// RedisStorage defines the storage requirements for Redis
type RedisStorage struct {
	// Size is the size of the storage to allocate to each Redis instance
	Size string `json:"size"`

	// StorageClassName is the name of the StorageClass used for provisioning volumes
	StorageClassName string `json:"storageClassName,omitempty"`
}

// RedisResources defines the CPU and memory resource requirements
type RedisResources struct {
	// Requests specifies the minimum amount of compute resources required.
	Requests Requests `json:"requests"`
	// Limits specifies the maximum amount of compute resources required.
	Limits Limits `json:"limits,omitempty"`
}

type Requests struct {
	// CPU request and limit
	CPU string `json:"cpu"`
	// Memory request and limit
	Memory string `json:"memory"`
}
type Limits struct {
	// CPU request and limit
	CPU string `json:"cpu"`
	// Memory request and limit
	Memory string `json:"memory"`
}

// RedisStatus defines the observed state of Redis
type RedisStatus struct {

	// ReadyReplicas is the number of replicas that are ready and serving requests.
	ReadyReplicas int32 `json:"readyReplicas"`
	// TotalReplicas is the total number of desired replicas.
	TotalReplicas int32 `json:"totalReplicas"`

	// Conditions represent the latest available observations of an object's state.
	Conditions []appsv1.DeploymentCondition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Redis is the Schema for the redis API
type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisSpec   `json:"spec,omitempty"`
	Status RedisStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RedisList contains a list of Redis
type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Redis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Redis{}, &RedisList{})
}
