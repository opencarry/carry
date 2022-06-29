package v1

type Namespace struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`
	Spec       NamespaceSpec   `json:"spec,omitempty"`
	Status     NamespaceStatus `json:"status,omitempty"`
}

type NamespaceList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Namespace `json:"items"`
}

type NamespaceSpec struct {
}

type NamespaceStatus struct {
	Phase NamespacePhase `json:"phase,omitempty"`
}

type NamespacePhase string

const (
	NamespaceActive      NamespacePhase = "active"
	NamespaceTerminating NamespacePhase = "terminating"
)
