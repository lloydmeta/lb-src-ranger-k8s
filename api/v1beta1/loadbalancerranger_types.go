/*

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LbSrcRangerSpec defines the desired state of LbSrcRanger
type LbSrcRangerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// TargetLabels is a selector for finding LoadBalancer Services
	// that need tending-to
	TargetLabels map[string]string `json:"target_labels,omitempty"`

	// UpdateEvery is the duration to wait between reconciles
	UpdateEvery metav1.Duration `json:"update_every"`

	// SrcIpUrls holds urls that return newline-separated lists of IP
	// address (ranges) to use as source ip ranges in the selected LoadBalancer
	// serviceds
	SrcIPUrls []string `json:"src_ip_urls,omitempty"`
}

// LbSrcRangerStatus defines the observed state of LbSrcRanger
type LbSrcRangerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// UpdatedCount is the number of servies last updated
	LastUpdatedCount int `json:"last_updated_count"`

	// LastRunAt is the time that there was a load balancer run
	LastRunAt metav1.Time `json:"last_run_at"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// LbSrcRanger is the Schema for the loadbalancerrangers API
type LbSrcRanger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LbSrcRangerSpec   `json:"spec,omitempty"`
	Status LbSrcRangerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LbSrcRangerList contains a list of LbSrcRanger
type LbSrcRangerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LbSrcRanger `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LbSrcRanger{}, &LbSrcRangerList{})
}
