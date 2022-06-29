package v1

import (
	"time"
)

type ReplicaSet struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec ReplicaSetSpec `json:"spec,omitempty"`

	Status ReplicaSetStatus `json:"status,omitempty"`
}

type ReplicaSetList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []ReplicaSet `json:"items"`
}

type ReplicaSetSpec struct {
	// required
	Selector *LabelSelector `json:"selector"`
	// required, same as pod
	Template PodTemplateSpec `json:"template"`
	// optional, defaults to 1
	Replicas *int64 `json:"replicas,omitempty"`

	Strategy ReplicaSetStrategy `json:"strategy,omitempty"`
}

type ReplicaSetStrategy struct {
	// required
	Type ReplicaSetStrategyType `json:"type"`
}

type ReplicaSetStrategyType string

const (
	InplaceUpdateReplicaSetStrategyType ReplicaSetStrategyType = "inplace_update"
)

type ReplicaSetStatus struct {
	Replicas             int `json:"replicas,omitempty"`
	FullyLabeledReplicas int `json:"fully_labeled_replicas,omitempty"`
	UpdatedReplicas      int `json:"updated_replicas,omitempty"`
	ReadyReplicas        int `json:"ready_replicas,omitempty"`

	Conditions []ReplicaSetCondition `json:"conditions,omitempty"`
}

type ReplicaSetCondition struct {
	Type  ReplicaSetConditionType `json:"type"`
	State ConditionState          `json:"state"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`
	LastUpdateTime     time.Time `json:"last_update_time,omitempty"`

	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type ReplicaSetConditionType string

const (
	ReplicaSetReplicaFailure ReplicaSetConditionType = "replica_failure"
)
