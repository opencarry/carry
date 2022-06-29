package v1

import (
	"time"
)

type Event struct {
	TypeMeta       `json:",inline"`
	ObjectMeta     `json:"metadata,omitempty"`
	InvolvedObject ObjectReference `json:"involved_object"`
	// Count 事件合并, 相同事件出现的次数
	Count int64 `json:"count,omitempty"`
	// Type Warning...
	Type      string      `json:"type,omitempty"`
	Reason    string      `json:"reason,omitempty"`
	Message   string      `json:"message,omitempty"`
	Source    EventSource `json:"source,omitempty"`
	FirstTime time.Time   `json:"first_time,omitempty"`
	LastTime  time.Time   `json:"last_time,omitempty"`
}

type EventSource struct {
	// Component from which the event is generated.
	Component string
	// Node name on which the event is generated.
	Host string
}

type EventList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Event `json:"items"`
}
