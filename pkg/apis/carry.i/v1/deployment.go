package v1

import (
	"time"
)

type Deployment struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec DeploymentSpec `json:"spec,omitempty"`

	Status DeploymentStatus `json:"status,omitempty"`
}

type DeploymentList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Deployment `json:"items"`
}

type DeploymentSpec struct {
	// required
	Selector *LabelSelector `json:"selector"`
	// required, same as pod
	Template PodTemplateSpec `json:"template"`
	// optional, defaults to 1
	Replicas *int64 `json:"replicas,omitempty"`
	// The deployment strategy to use to replace existing pods with new ones.
	Strategy DeploymentStrategy `json:"strategy,omitempty"`
	// optional, defaults to 10
	RevisionHistoryLimit *int64 `json:"revision_history_limit,omitempty"`
	// optional, defaults to 600
	ProgressDeadlineSeconds *int64 `json:"progress_deadline_seconds,omitempty"`
}

type DeploymentStrategy struct {
	// required
	Type DeploymentStrategyType `json:"type"`
}

type DeploymentStrategyType string

const (
	InplaceUpdateDeploymentStrategyType DeploymentStrategyType = "inplace_update"
)

type DeploymentStatus struct {
	Replicas        int `json:"replicas,omitempty"`
	UpdatedReplicas int `json:"updated_replicas,omitempty"`
	ReadyReplicas   int `json:"ready_replicas,omitempty"`

	Conditions []DeploymentCondition `json:"conditions,omitempty"`
}

type DeploymentCondition struct {
	Type  DeploymentConditionType `json:"type"`
	State ConditionState          `json:"state"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`
	LastUpdateTime     time.Time `json:"last_update_time,omitempty"`

	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type DeploymentConditionType string

const (
	DeploymentAvailable DeploymentConditionType = "available"
	// DeploymentProgressing scale
	DeploymentProgressing DeploymentConditionType = "progressing"
	// DeploymentReplicaFailure pod 创建或者删除失败时
	DeploymentReplicaFailure DeploymentConditionType = "replica_failure"
)
