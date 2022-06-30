package v1

import "time"

type StatefulSet struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec StatefulSetSpec `json:"spec,omitempty"`

	Status StatefulSetStatus `json:"status,omitempty"`
}

type StatefulSetList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []StatefulSet `json:"items"`
}

type StatefulSetSpec struct {
	// required
	Selector *LabelSelector `json:"selector"`
	// required, same as pod
	Template PodTemplateSpec `json:"template"`
	// optional, defaults to 1
	Replicas *int64 `json:"replicas,omitempty"`
	// strategy to use to replace existing pods.
	Strategy StatefulSetStrategy `json:"strategy,omitempty"`
	// optional, defaults to 10
	RevisionHistoryLimit *int64 `json:"revision_history_limit,omitempty"`
}

type StatefulSetStrategy struct {
	// required
	Type StatefulSetStrategyType `json:"type"`
}

type StatefulSetStrategyType string

const (
	// InplaceUpdateStatefulSetStrategyType 原地升级，当template有变动，直接将变动同步到对应的Pod，不创建新的Pod
	InplaceUpdateStatefulSetStrategyType StatefulSetStrategyType = "inplace_update"
)

type StatefulSetStatus struct {
	// replicas is the number of Pods created by the StatefulSet controller.
	Replicas int `json:"replicas,omitempty"`
	// readyReplicas is the number of pods created for this StatefulSet with a Ready Condition.
	ReadyReplicas int `json:"ready_replicas,omitempty"`
	// currentReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by currentRevision.
	CurrentReplicas int `json:"current_replicas,omitempty"`
	// updatedReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version indicated by updateRevision.
	UpdatedReplicas int `json:"updated_replicas,omitempty"`
	// currentRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas).
	CurrentRevision string `json:"current_revision,omitempty"`
	// updateRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas)
	UpdateRevision string `json:"update_revision,omitempty"`

	Conditions []StatefulSetCondition `json:"conditions,omitempty"`
}

type StatefulSetCondition struct {
	Type  StatefulSetConditionType `json:"type"`
	State ConditionState           `json:"state"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`
	LastUpdateTime     time.Time `json:"last_update_time,omitempty"`

	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type StatefulSetConditionType string

const (
	StatefulSetAvailable StatefulSetConditionType = "available"
)
