package v1

import "time"

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
	// strategy to use to replace existing pods.
	Strategy ReplicaSetStrategy `json:"strategy,omitempty"`
	// default to 0
	MinReadySeconds int64 `json:"min_ready_seconds,omitempty"`
}

type ReplicaSetStrategy struct {
	// required
	Type ReplicaSetStrategyType `json:"type"`
}

type ReplicaSetStrategyType string

const (
	// InplaceUpdateReplicaSetStrategyType 原地升级，当template有变动，直接将变动同步到对应的Pod，不创建新的Pod
	InplaceUpdateReplicaSetStrategyType ReplicaSetStrategyType = "inplace_update"
)

type ReplicaSetStatus struct {
	Replicas             int64 `json:"replicas,omitempty"`
	FullyLabeledReplicas int64 `json:"fully_labeled_replicas,omitempty"`
	UpdatedReplicas      int64 `json:"updated_replicas,omitempty"`
	ReadyReplicas        int64 `json:"ready_replicas,omitempty"`
	AvailableReplicas    int64 `json:"available_replicas,omitempty"`

	// ObservedGeneration reflects the generation of the most recently observed ReplicaSet.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

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
	// ReplicaSetReplicaFailure means one of its pods fails to be created
	ReplicaSetReplicaFailure ReplicaSetConditionType = "replica_failure"
)
