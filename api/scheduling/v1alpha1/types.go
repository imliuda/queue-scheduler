/*
Copyright 2020 The Kubernetes Authors.

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:shortName={queue,queues},scope=Cluster

type Queue struct {
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata"`

	// QueueSpec defines the queue config
	Queues []Queues `json:"queues"`
}

// Queues defines the the queues
type Queues struct {
	// Name is the queue name in current level
	Name string `json:"name"`
	// Min is the lower limit of desired guaranteed resources
	// +optional
	Min v1.ResourceList `json:"min,omitempty"`

	// Guarantee is the upper limit of desired guaranteed resources
	// +optional
	Max v1.ResourceList `json:"max,omitempty"`

	// Weight is the weight through same level queue
	// +optional
	Weight int `json:"weight,omitempty"`

	// Child queues
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Queues []Queues `json:"queues"`

	// Properties define queue custom configs
	// +optional
	Properties map[string]string `json:"properties"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// QueueList is a list of Queue items
type QueueList struct {
	metav1.TypeMeta `json:",inline"`

	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of Queue objects.
	Items []Queue `json:"items"`
}
