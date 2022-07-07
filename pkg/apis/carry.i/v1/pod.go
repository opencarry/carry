package v1

import (
	"time"
)

type Pod struct {
	TypeMeta   `json:",omitempty"`
	ObjectMeta `json:"metadata,omitempty"`

	Spec PodSpec `json:"spec,omitempty"`

	Status PodStatus `json:"status,omitempty"`
}

type PodList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Pod `json:"items"`
}

type PodTemplateSpec struct {
	ObjectMeta `json:"metadata,omitempty"`

	Spec PodSpec `json:"spec,omitempty"`
}

type PodSpec struct {
	// 安装容器，负责安装/升级程序，满足目前很多开源中间件使用rpm包安装的需求
	// 可以定义多个容器，按顺序执行，程序结束返回0视为成功，如果运行失败（exit非0），根据spec.restart_policy字段决定是否重试
	// 一个Pod被调度到一台新Node上，或installation_containers里的内容更新时会执行
	// 安装容器需要自己负责判断是安装或升级的场景
	// 单纯Pod重启过程不会执行installation_containers
	// !!!不建议使用installation_containers，能不用就不用，因为安装过程脱离了carry的控制，不确定性的错误增加
	InstallationContainers []Container `json:"installation_containers,omitempty"`

	// 卸载容器，对应installation_containers，负责卸载程序
	// 可以定义多个容器，按顺序执行，程序结束返回0视为成功
	// 删除Pod时，主容器全部停止后，触发卸载容器执行
	UninstallationContainers []Container `json:"uninstallation_containers,omitempty"`

	// 初始化容器，安装容器执行完毕后执行
	// 多个按顺序执行，只要有一个执行失败，则视为Pod失败，重启策略取决spec.restart_policy
	// 启动/重启Pod时执行
	InitContainers []Container `json:"init_containers,omitempty"`

	// 主容器列表，同时启动，多个容器不保证启动顺序有规律
	// 【必填】至少要有一个容器
	Containers []Container `json:"containers"`

	Volumes []Volume `json:"volumes,omitempty"`

	// pod被调度到此Node，如果此值为空，scheduler负责填充
	NodeName string `json:"node_name,omitempty"`

	NodeSelector map[string]string `json:"node_selector,omitempty"`

	// 一组亲和性调度规则
	Affinity Affinity `json:"affinity,omitempty"`

	// Defaults to always
	RestartPolicy RestartPolicy `json:"restart_policy,omitempty"`

	SecurityContext *PodSecurityContext `json:"security_context,omitempty"`

	// 优雅删除时间周期，若为0，则立即删除
	// 默认为30
	TerminationGracePeriodSeconds *int64 `json:"termination_grace_period_seconds,omitempty"`

	// Pod存活时间，时间到如果还未结束则carry会主动终止
	// 可用于Job类型的Pod
	ActiveDeadlineSeconds *int64 `json:"active_deadline_seconds,omitempty"`

	// 如果想停止Pod，将此值置为true，默认为false
	// carry 感知到此值为true后，将Pod中所有容器停止
	Suspended *bool `json:"suspended,omitempty"`

	// 可选，默认 default-scheduler
	SchedulerName string `json:"scheduler_name,omitempty"`
}

type RestartPolicy string

const (
	RestartPolicyAlways    RestartPolicy = "always"
	RestartPolicyOnFailure RestartPolicy = "on_failure"
	RestartPolicyNever     RestartPolicy = "never"
)

type PodSecurityContext struct {
	RunAsUser    string `json:"run_as_user,omitempty"`
	RunAsGroup   string `json:"run_as_group,omitempty"`
	RunAsNonRoot *bool  `json:"run_as_non_root,omitempty"`
}

type Container struct {
	// DNS_LABEL.
	Name string `json:"name"`
	// mountain image name
	Image string `json:"image"`
	// Defaults to always if :latest tag is specified, or if_not_present otherwise
	ImagePullPolicy PullPolicy `json:"image_pull_policy,omitempty"`
	// 部署目录，绝对路径
	ImageDeploymentDir string `json:"image_deployment_dir"`

	Env []EnvVar `json:"env,omitempty"`
	// 容器进程工作目录
	WorkingDir string `json:"working_dir,omitempty"`
	// 启动程序命令
	// 建议不放在shell里启动，比如程序bin文件是passport，那么command直接写成passport，而不是封装一个脚本start.sh来启动
	Command []string `json:"command"`
	// 停止程序命令
	TerminationCommand []string `json:"termination_command,omitempty"`

	SecurityContext *SecurityContext `json:"security_context,omitempty"`

	// 容器暴露的监听端口列表，为了让外界知道怎么连接进来
	Ports []ContainerPort `json:"ports,omitempty"`

	VolumeMounts []VolumeMount `json:"volume_mounts,omitempty"`

	Resources ResourceRequirements `json:"resources,omitempty"`
}

// PullPolicy describes a policy for if/when to pull a container image
type PullPolicy string

const (
	PullAlways       PullPolicy = "always"
	PullNever        PullPolicy = "never"
	PullIfNotPresent PullPolicy = "if_not_present"
)

type EnvVar struct {
	Name      string        `json:"name"`
	Value     string        `json:"value,omitempty"`
	ValueFrom *EnvVarSource `json:"value_from,omitempty"`
}

type EnvVarSource struct {
	FieldRef *ObjectFieldSelector `json:"field_ref,omitempty"`
}

type ObjectFieldSelector struct {
	FieldPath string `json:"field_path,omitempty"`
}

type SecurityContext struct {
	RunAsUser string `json:"run_as_user,omitempty"`
}

type ContainerPort struct {
	Name          string   `json:"name,omitempty"`
	ContainerPort int      `json:"container_port"`
	Protocol      Protocol `json:"protocol,omitempty"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mount_path"`
	SubPath   string `json:"sub_path,omitempty"`
}

type Affinity struct {
	// Pod节点亲和性
	PodAffinity *PodAffinity `json:"pod_affinity,omitempty"`

	// Pod反亲和性
	PodAntiAffinity *PodAntiAffinity `json:"pod_anti_affinity,omitempty"`
}

type PodAffinity struct {
	// 本Pod必须调度到匹配这些规则的Pod所运行Node上
	Required *Required `json:"required,omitempty"`
}

type PodAntiAffinity struct {
	// 本Pod必须不能调度到匹配这些规则的Pod所运行Node上
	Required *Required `json:"required,omitempty"`
}

type Required struct {
	MatchLabels map[string]string `json:"match_labels,omitempty"`
}

type Volume struct {
	Name         string `json:"name"`
	VolumeSource `json:",inline"`
}

type VolumeSource struct {
	ConfigMap *ConfigMapVolumeSource `json:"configMap,omitempty"`
}

type ConfigMapVolumeSource struct {
	LocalObjectReference `json:",inline"`
	// If unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	Items []KeyToPath `json:"items,omitempty"`
	// Optional: mode bits to use on created files by default. Must be a
	// value between 0 and 0777. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	DefaultMode *int64 `json:"default_mode,omitempty"`
	// Specify whether the ConfigMap or it's keys must be defined
	Optional *bool `json:"optional,omitempty"`
}

type LocalObjectReference struct {
	// Name of the referent.
	Name string `json:"name,omitempty"`
}

// KeyToPath Maps a string key to a path within a volume.
type KeyToPath struct {
	// The key to project.
	Key string `json:"key"`

	// The relative path of the file to map the key to.
	// May not be an absolute path.
	// May not contain the path element '..'.
	// May not start with the string '..'.
	Path string `json:"path"`
	// Optional: mode bits to use on this file, must be a value between 0
	// and 0777. If not specified, the volume defaultMode will be used.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	Mode *int32 `json:"mode,omitempty"`
}

// PodPhase is a label for the condition of a pod at the current time.
type PodPhase string

// These are the valid statuses of pods.
const (
	// PodPending Pod被创建后，在任何一个容器启动之前，其中包括调度时间、制品下载时间、安装时间
	PodPending PodPhase = "pending"
	// PodRunning 至少有一个容器在运行，或容器在重启过程中
	PodRunning PodPhase = "running"
	// PodSucceeded 所有容器都已运行成功并停止，并且不会再重启（Job类型的Pod）
	PodSucceeded PodPhase = "succeeded"
	// PodFailed 所有容器都已停止，并且至少有一个容器因为错误停止，比如进程退出码非0，或被carry终止
	PodFailed PodPhase = "failed"
	// PodSuspended 主动暂停（spec.suspended=true），所有容器都被carry停止
	PodSuspended PodPhase = "suspended"
	// PodUnknown 未知状态，比如无法与node网络通讯
	PodUnknown PodPhase = "unknown"
)

type PodStatus struct {
	Phase PodPhase `json:"phase,omitempty"`
	// carry首次感知到此Pod的时间，在拉取镜像之前
	StartTime time.Time `json:"start_time,omitempty"`
	// 玷污重启，默认为0
	// 如果想重启，则将此值置为1，carry感知到此值后负责重启Pod，并将此值置为0
	TaintRestarts *int64 `json:"taint_restarts,omitempty"`
	// Pod部署的机器IP
	HostIp string `json:"host_ip,omitempty"`

	Reason string `json:"reason,omitempty"`

	Message string `json:"message,omitempty"`
	// conditions包含详细的Pod状态
	Conditions []PodCondition `json:"conditions,omitempty"`

	ContainerStatuses []ContainerStatus `json:"container_statuses,omitempty"`

	InitContainerStatuses []ContainerStatus `json:"init_container_statuses,omitempty"`

	InstallationContainerStatuses []ContainerStatus `json:"installation_container_statuses,omitempty"`

	UninstallationContainerStatuses []ContainerStatus `json:"uninstallation_container_statuses,omitempty"`
}

// PodConditionType is a valid value for PodCondition.Type
type PodConditionType string

// These are valid conditions of pod.
const (
	// PodScheduled Pod已经被调度到Node
	PodScheduled PodConditionType = "pod_scheduled"
	// PodContainersReady 所有容器都已经ready，已经把所有的制品等信息下载完成
	PodContainersReady PodConditionType = "containers_ready"
	// PodInstalled 所有安装容器成功执行
	PodInstalled PodConditionType = "installed"
	// PodInitialized 所有初始化容器都已经成功运行
	PodInitialized PodConditionType = "initialized"
	// PodReady Pod已经准备好接收请求
	PodReady PodConditionType = "ready"
	// PodReasonUnschedulable Pod未能成功调度到node，比如没有符合Pod资源需求的Node
	PodReasonUnschedulable PodConditionType = "unschedulable"
)

type PodCondition struct {
	// 类型type包含如下几种
	// pod_scheduled: Pod已经被调度到Node
	// containers_ready: 所有容器都已经ready，已经把所有的制品等信息下载完成
	// installed: 所有安装容器成功执行
	// initialized: 所有初始化容器都已经成功运行
	// ready: Pod已经准备好接收请求
	Type PodConditionType `json:"type"`

	State ConditionState `json:"state"`

	LastProbeTime time.Time `json:"last_probe_time,omitempty"`

	LastTransitionTime time.Time `json:"last_transition_time,omitempty"`

	Reason string `json:"reason,omitempty"`

	Message string `json:"message,omitempty"`
}

type ContainerStatus struct {
	// DNS_LABEL, 在同一个pod中必须是唯一的
	Name  string `json:"name"`
	Pid   *int64 `json:"pid,omitempty"`
	Image string `json:"image"`
	// sha256:xxx
	ImageId string `json:"image_id,omitempty"`
	//
	State ContainerState `json:"state,omitempty"`
	// 容器是否能接收流量请求
	Ready bool `json:"ready"`

	RestartCount int64 `json:"restart_count"`
}

type ContainerState struct {
	Waiting    *ContainerStateWaiting    `json:"waiting,omitempty"`
	Running    *ContainerStateRunning    `json:"running,omitempty"`
	Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
}

type ContainerStateWaiting struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type ContainerStateRunning struct {
	// 容器启动（重启）时间
	StartTime time.Time `json:"start_time,omitempty"`
}

type ContainerStateTerminated struct {
	// Exit status from the last termination of the container
	ExitCode int `json:"exit_code"`
	// Signal from the last termination of the container
	Signal int `json:"signal,omitempty"`
	// (brief) reason from the last termination of the container
	Reason string `json:"reason,omitempty"`
	// Message regarding the last termination of the container
	Message string `json:"message,omitempty"`
	// Time at which previous execution of the container started
	StartTime time.Time `json:"start_time,omitempty"`
	// Time at which the container last terminated
	FinishTime time.Time `json:"finish_time,omitempty"`
}
