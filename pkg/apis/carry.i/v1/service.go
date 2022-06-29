package v1

import "time"

type Service struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec ServiceSpec `json:"spec,omitempty"`

	Status ServiceStatus `json:"status,omitempty"`
}

type ServiceList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Service `json:"items"`
}

type ServiceSpec struct {
	Selector map[string]string `json:"selector,omitempty"`

	Ports []ServicePort `json:"ports,omitempty"`
}

type ServicePort struct {
	Name string `json:"name,omitempty"`

	Address string `json:"address"`

	Port int `json:"port"`

	Protocol Protocol `json:"protocol,omitempty"`

	TargetPort int `json:"target_port"`

	// 端口注解，用于扩展service，第三方插件可将自定义的注解信息写在此处供插件自身使用
	Annotations map[string]string `json:"annotations,omitempty"`
}

type ServiceStatus struct {
	Conditions []ServiceCondition `json:"conditions,omitempty"`
}

type ServiceConditionType string

const (
	ServiceAvailable ServiceConditionType = "available"
)

type ServiceCondition struct {
	Type ServiceConditionType `json:"type"`

	State ConditionState `json:"state"`

	LastProbeTime time.Time `json:"last_probe_time,omitempty"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`

	Reason string `json:"reason,omitempty"`

	Message string `json:"message,omitempty"`
}
