package v1

import "time"

type Node struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec NodeSpec `json:"spec,omitempty"`

	Status NodeStatus `json:"status,omitempty"`
}

type NodeList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Node `json:"items"`
}

type NodeSpec struct {
	Unschedulable bool `json:"unschedulable,omitempty"`
	// 允许使用的资源容量
	Capacity ResourceList `json:"capacity,omitempty"`
}

type NodeStatus struct {
	Capacity   ResourceList     `json:"capacity,omitempty"`
	Phase      NodePhase        `json:"phase,omitempty"`
	Conditions []NodeCondition  `json:"conditions,omitempty"`
	Addresses  []NodeAddress    `json:"addresses,omitempty"`
	NodeInfo   NodeSystemInfo   `json:"node_info,omitempty"`
	Images     []ContainerImage `json:"images,omitempty"`
}

type ContainerImage struct {
	Names []string `json:"names"`
	// The size of the image in bytes.
	SizeBytes int64 `json:"size_bytes,omitempty"`
}

type NodePhase string

const (
	NodePending    NodePhase = "pending"
	NodeRunning    NodePhase = "running"
	NodeTerminated NodePhase = "terminated"
)

type NodeAddress struct {
	Type    NodeAddressType `json:"type"`
	Address string          `json:"address"`
}
type NodeAddressType string

const (
	NodeHostName   NodeAddressType = "hostname"
	NodeInternalIP NodeAddressType = "internal_ip"
)

type NodeSystemInfo struct {
	Architecture    string `json:"architecture"`
	OperatingSystem string `json:"operating_system"`
	KernelVersion   string `json:"kernel_version"`
	CarryVersion    string `json:"carry_version"`
	OSImage         string `json:"os_image"`
	Cpu             string `json:"cpu"`
	Memory          int64  `json:"memory"`
	Disk            int64  `json:"disk"`
}

type NodeConditionType string

const (
	NodeReady          NodeConditionType = "ready"
	NodeMemoryPressure NodeConditionType = "memory_pressure"
	NodeDiskPressure   NodeConditionType = "disk_pressure"
	NodePIDPressure    NodeConditionType = "pid_pressure"
)

type NodeCondition struct {
	Type NodeConditionType `json:"type"`

	State ConditionState `json:"state"`

	LastProbeTime time.Time `json:"last_probe_time,omitempty"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`

	Reason string `json:"reason,omitempty"`

	Message string `json:"message,omitempty"`
}
