package v1

type ConfigMap struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	// 每个键必须由字母数字和字符'-', '_' or '.'组成
	Data map[string]string `json:"data,omitempty"`
	// 每个键必须由字母数字和字符'-', '_' or '.'组成
	BinaryData map[string][]byte `json:"binary_data,omitempty"`
}

type ConfigMapList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []ConfigMap `json:"items"`
}
