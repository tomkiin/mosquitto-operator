package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MosquittoSpec struct {
	Count int32  `json:"count"`
	Image string `json:"image"`
	Conf  string `json:"conf"`
}

type MosquittoStatus struct {
	ClusterReloading bool `json:"clusterReloading"`
}

type Mosquitto struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MosquittoSpec   `json:"spec,omitempty"`
	Status MosquittoStatus `json:"status,omitempty"`
}

type MosquittoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mosquitto `json:"items"`
}
