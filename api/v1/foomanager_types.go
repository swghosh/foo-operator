/*
Copyright 2025.

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

package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// FooManagerSpec defines the desired state of FooManager.
type FooManagerSpec struct {
	// UseIn is an example field that decides what we use for :x
	UseIn string `json:"useIn,omitempty"`

	// Template for pods in the deployment
	// +optional
	TemplateSpec *corev1.PodSpec `json:"templateSpec,omitempty"`
}

// FooManagerStatus defines the observed state of FooManager.
type FooManagerStatus struct {
	// WasUsedFor reports what this field was last used for :x
	WasUsedFor string `json:"wasUsedFor"`

	// LastTransitionTime time when we last observed a change
	LastTransitionTime *metav1.Time `json:"lastTransitionTime"`

	// WeirdReport of deployment
	WeirdReport string `json:"weirdReport,omitempty"`

	// DeploymentStatus
	// +optional
	DeploymentStatus *appsv1.DeploymentStatus `json:"deploymentStatus,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// FooManager is the Schema for the foomanagers API.
type FooManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooManagerSpec   `json:"spec,omitempty"`
	Status FooManagerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// FooManagerList contains a list of FooManager.
type FooManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FooManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FooManager{}, &FooManagerList{})
}
