package v1

import "time"

type Job struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec JobSpec `json:"spec,omitempty"`

	Status JobStatus `json:"status,omitempty"`
}

type JobList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Job `json:"items"`
}

type JobSpec struct {
	// Specifies the number of retries before marking this job failed. Defaults to 6
	BackoffLimit *int `json:"backoff_limit,omitempty"`

	Completions *int `json:"completions,omitempty"`

	Parallelism *int `json:"parallelism,omitempty"`
	// Specifies the duration in seconds relative to the startTime that the job may be continuously active
	// before the system tries to terminate it; value must be positive integer.
	ActiveDeadlineSeconds *int64 `json:"active_deadline_seconds,omitempty"`

	Selector *LabelSelector `json:"selector,omitempty"`

	Template PodTemplateSpec `json:"template"`
}

type JobStatus struct {
	CompletionTime time.Time `json:"completion_time,omitempty"`

	StartTime time.Time `json:"start_time,omitempty"`
	// The number of pods which reached phase Succeeded.
	Succeeded int `json:"succeeded,omitempty"`
	// The number of pods which reached phase Failed.
	Failed int `json:"failed,omitempty"`
	// The number of pending and running pods.
	Active int `json:"active,omitempty"`

	Conditions []JobCondition `json:"conditions,omitempty"`
}

type JobConditionType string

// These are built-in conditions of a job.
const (
	// JobSuspended means the job has been suspended.
	JobSuspended JobConditionType = "suspended"
	// JobComplete means the job has completed its execution.
	JobComplete JobConditionType = "complete"
	// JobFailed means the job has failed its execution.
	JobFailed JobConditionType = "failed"
)

type JobCondition struct {
	Type JobConditionType `json:"type"`

	LastProbeTime time.Time `json:"last_probe_time,omitempty"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`

	State ConditionState `json:"state"`

	Reason string `json:"reason,omitempty"`

	Message string `json:"message,omitempty"`
}
