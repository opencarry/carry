package v1

import "github.com/opencarry/carry/pkg/resource"

const (
	DefaultTerminationGracePeriodSeconds = 30

	DefaultActiveDeadlineSeconds = 36000

	DefaultSuspended = false

	ConfigMapVolumeSourceDefaultMode int64 = 0644

	// DefaultSchedulerName "default-scheduler" is the name of default scheduler.
	DefaultSchedulerName = "default-scheduler"
)

type ResourceName string

const (
	// ResourceCPU CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// ResourceMemory Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
)

type ResourceList map[ResourceName]resource.Quantity

// Cpu Returns the CPU limit if specified.
func (rl *ResourceList) Cpu() *resource.Quantity {
	if val, ok := (*rl)[ResourceCPU]; ok {
		return &val
	}
	return &resource.Quantity{}
}

// Memory Returns the Memory limit if specified.
func (rl *ResourceList) Memory() *resource.Quantity {
	if val, ok := (*rl)[ResourceMemory]; ok {
		return &val
	}
	return &resource.Quantity{}
}

type ResourceRequirements struct {
	Limits   ResourceList `json:"limits,omitempty"`
	Requests ResourceList `json:"requests,omitempty"`
}

type LabelSelector struct {
	MatchLabels map[string]string `json:"match_labels,omitempty"`

	MatchExpressions []LabelSelectorRequirement `json:"match_expressions,omitempty"`
}

type LabelSelectorRequirement struct {
	Key string `json:"key"`

	Operator LabelSelectorOperator `json:"operator"`

	Values []string `json:"values,omitempty"`
}

type LabelSelectorOperator string

const (
	LabelSelectorOpIn           LabelSelectorOperator = "in"
	LabelSelectorOpNotIn        LabelSelectorOperator = "not_in"
	LabelSelectorOpExists       LabelSelectorOperator = "exists"
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "does_not_exist"
)

type ConditionState string

const (
	ConditionTrue    ConditionState = "true"
	ConditionFalse   ConditionState = "false"
	ConditionUnknown ConditionState = "unknown"
)

type Protocol string

const (
	// ProtocolTCP is the TCP protocol.
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP is the UDP protocol.
	ProtocolUDP Protocol = "UDP"
)

type ObjectReference struct {
	Kind            string `json:"kind"`
	Namespace       string `json:"namespace"`
	Name            string `json:"name"`
	UID             UID    `json:"uid"`
	APIVersion      string `json:"api_version"`
	ResourceVersion string `json:"resource_version"`
	// Optional. If referring to a piece of an object instead of an entire object, this string
	// should contain information to identify the sub-object. For example, if the object
	// reference is to a container within a pod, this would take on a value like:
	// "spec.containers{name}" (where "name" refers to the name of the container that triggered
	// the event) or if no container name is specified "spec.containers[2]" (container with
	// index 2 in this pod). This syntax is chosen only to have some well-defined way of
	// referencing a part of an object.
	FieldPath string `json:"field_path,omitempty"`
}

func (o ObjectReference) GetObjectKind() string {
	return o.Kind
}

type Binding struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	PodID string `json:"pod_id"`
	Host  string `json:"host"`
}
