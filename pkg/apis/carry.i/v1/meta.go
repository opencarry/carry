package v1

import "time"

type UID string

type TypeMeta struct {
	// 小写模式
	Kind string `json:"kind,omitempty"`
	//
	APIVersion string `json:"api_version,omitempty"`
}

type ListMeta struct {
	// etcd ModRevision
	ResourceVersion string `json:"resource_version,omitempty"`
	// 分页时会赋值
	Continue string `json:"continue,omitempty"`
	// 剩余条数
	RemainingItemCount *int64 `json:"remaining_item_count,omitempty"`
}

type ObjectMeta struct {
	// 资源名称，在命令空间下唯一
	// 由域名DNS_LABEL组成，最长128字节
	// 不允许更新
	Name string `json:"name,omitempty"`
	// 命名空间，默认default
	// 由域名DNS_LABEL组成，不允许更新
	Namespace string `json:"namespace,omitempty"`

	Labels map[string]string `json:"labels,omitempty"`
	// 非结构化描述类型信息
	// carry.i/ 开头是系统保留前缀
	Annotations map[string]string `json:"annotations,omitempty"`
	// 依赖的资源列表
	// 例如pod，如果这个列表里资源全部被删除了，那么当前这个Pod也会被回收
	OwnerReferences []OwnerReference `json:"owner_references,omitempty"`

	// 资源创建时间，格式：RFC3339，其它地方时间格式同样
	// 由服务器端设置，不允许更新
	CreationTime time.Time `json:"creation_time,omitempty"`
	// 删除时间，到此时间，资源将被从系统中删除
	// 当资源被请求优雅删除时，系统会设置该时间，只要该字段被设置，系统将启动该资源的删除流程
	// 举例子：某Pod资源被优雅删除，该值设置为30s后的时间点，carry感知到之后开始执行如下流程：
	// 发送TERM信号（若有则执行termination_command）--> （若时间到deletion_time还未终止）发送KILL信号（等2s）  --> 执行uninstallation_containers --> （若卸载成功）carry请求server删除Pod
	// 整个过程status字段相应字段要有对应设置
	// 不允许更新
	DeletionTime time.Time `json:"deletion_time,omitempty"`
	// 留给优雅终止的时间秒数，超过这个时间资源将被从系统中删除
	// 只有deletion_time字段有值时，此字段才能被设置，时间要比deletion_time小
	DeletionGracePeriodSeconds *int64 `json:"deletion_grace_period_seconds,omitempty"`

	// etcd ModRevision
	// 系统内部标识资源的版本，不透明的字符串，client端不要对此字段作任何假设，只需在需要的地方原样传回server端
	// 可用于判断资源是否更新
	// 可用于控制并发更新资源时导致的冲突
	// 可用于watch资源时用的位置游标
	ResourceVersion string `json:"resource_version,omitempty"`

	// 资源的唯一ID，格式UUID
	// 为了区别同样name的资源，比如有个name=foo的Pod被删除后，又创建一个同名的Pod
	// 不允许更新
	UID UID `json:"uid,omitempty"`
}

type OwnerReference struct {
	APIVersion string `json:"api_version"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	UID        UID    `json:"uid"`
	// 控制器类资源只能有一个
	Controller *bool `json:"controller,omitempty"`
}
