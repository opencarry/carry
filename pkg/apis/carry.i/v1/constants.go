package v1

const (
	DeploymentKeyPrefix = "deployment.carry.i/"
	SchedulerKeyPrefix  = "scheduler.carry.i/"

	BindNodesAnnotationKey string = SchedulerKeyPrefix + "bind_nodes"

	StatefulSetPodNameLabel = "statefulset.carry.io/pod-name"

	ControllerRevisionHashLabelKey = "controller-revision-hash"
	StatefulSetRevisionLabel       = ControllerRevisionHashLabelKey
)
