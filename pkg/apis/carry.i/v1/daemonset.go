package v1

import "time"

type DaemonSet struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec DaemonSetSpec `json:"spec,omitempty"`

	Status DaemonSetStatus `json:"status,omitempty"`
}

type DaemonSetList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []DaemonSet `json:"items"`
}

type DaemonSetSpec struct {
	// required
	Selector *LabelSelector `json:"selector"`
	// required, same as pod
	Template PodTemplateSpec `json:"template"`
	// strategy to use to replace existing pods.
	Strategy DaemonSetStrategy `json:"strategy,omitempty"`
	// optional, defaults to 10
	RevisionHistoryLimit *int64 `json:"revision_history_limit,omitempty"`
}

type DaemonSetStrategy struct {
	// required
	Type DaemonSetStrategyType `json:"type"`
}

type DaemonSetStrategyType string

const (
	// InplaceUpdateDaemonSetStrategyType 原地升级，当template有变动，直接将变动同步到对应的Pod，不创建新的Pod
	InplaceUpdateDaemonSetStrategyType DaemonSetStrategyType = "inplace_update"
)

type DaemonSetStatus struct {
	// The total number of nodes that should be running the daemon pod (including nodes correctly running the daemon pod).
	DesiredNumberScheduled int `json:"desired_number_scheduled"`
	// The number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod
	CurrentNumberScheduled int `json:"current_number_scheduled"`
	// numberReady is the number of nodes that should be running the daemon pod and have one or more of the daemon pod running with a Ready Condition
	NumberReady int `json:"number_ready"`

	Conditions []DaemonSetCondition `json:"conditions,omitempty"`
}

type DaemonSetCondition struct {
	Type  DaemonSetConditionType `json:"type"`
	State ConditionState         `json:"state"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`
	LastUpdateTime     time.Time `json:"last_update_time,omitempty"`

	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type DaemonSetConditionType string

const (
	DaemonSetAvailable DaemonSetConditionType = "available"
)
