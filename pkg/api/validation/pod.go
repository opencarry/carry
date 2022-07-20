package validation

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/google/go-cmp/cmp"

	v1 "github.com/opencarry/carry/pkg/apis/carry.i/v1"
	"github.com/opencarry/carry/pkg/resource"
	"github.com/opencarry/carry/pkg/util/sets"
	"github.com/opencarry/carry/pkg/util/validation"
	"github.com/opencarry/carry/pkg/util/validation/field"
)

var ValidateNodeName = NameIsDNSSubdomain

var ValidatePodName = NameIsDNSSubdomain

func ValidatePod(pod *v1.Pod) field.ErrorList {
	fldPath := field.NewPath("metadata")
	allErrs := ValidateObjectMeta(&pod.ObjectMeta, true, NameIsDNSSubdomain, fldPath)
	allErrs = append(allErrs, ValidatePodSpecificAnnotations(pod.ObjectMeta.Annotations, &pod.Spec, fldPath.Child("annotations"))...)
	allErrs = append(allErrs, ValidatePodSpec(&pod.Spec, field.NewPath("spec"))...)

	return allErrs
}

func ValidatePodSpecificAnnotations(annotations map[string]string, spec *v1.PodSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	// BindNodesAnnotationKey
	return allErrs
}

// ValidatePodSpec tests that the specified PodSpec has valid data.
// This includes checking formatting and uniqueness.  It also canonicalizes the
// structure by setting default values and implementing any backwards-compatibility
// tricks.
func ValidatePodSpec(spec *v1.PodSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// volumes
	vols, vErrs := ValidateVolumes(spec.Volumes, fldPath.Child("volumes"))
	allErrs = append(allErrs, vErrs...)

	// containers
	allErrs = append(allErrs, validateContainers(spec.Containers, vols, fldPath.Child("containers"))...)
	// init_containers
	otherContainers := spec.Containers
	allErrs = append(allErrs, validateInitContainers(spec.InitContainers, otherContainers, vols, fldPath.Child("init_containers"))...)
	// installation_containers
	otherContainers = append(otherContainers, spec.InitContainers...)
	allErrs = append(allErrs, validateInstallationContainers(spec.InstallationContainers, otherContainers, vols, fldPath.Child("installation_containers"))...)
	// uninstallation_containers
	otherContainers = append(otherContainers, spec.InstallationContainers...)
	allErrs = append(allErrs, validateUninstallationContainers(spec.UninstallationContainers, otherContainers, vols, fldPath.Child("uninstallation_containers"))...)
	// affinity
	allErrs = append(allErrs, validateAffinity(spec.Affinity, fldPath.Child("affinity"))...)
	// node_name
	if len(spec.NodeName) > 0 {
		for _, msg := range ValidateNodeName(spec.NodeName, false) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("node_name"), spec.NodeName, msg))
		}
	}
	// node_selector
	allErrs = append(allErrs, ValidateLabels(spec.NodeSelector, fldPath.Child("node_selector"))...)
	// security_context
	allErrs = append(allErrs, ValidatePodSecurityContext(spec.SecurityContext, fldPath.Child("security_context"))...)
	// restart_policy
	allErrs = append(allErrs, validateRestartPolicy(&spec.RestartPolicy, fldPath.Child("restart_policy"))...)
	// termination_grace_period_seconds
	if spec.TerminationGracePeriodSeconds != nil {
		value := *spec.TerminationGracePeriodSeconds
		if value < 1 || value > math.MaxInt32 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("termination_grace_period_seconds"), value, validation.InclusiveRangeError(1, math.MaxInt32)))
		}
	}

	// active_deadline_seconds
	if spec.ActiveDeadlineSeconds != nil {
		value := *spec.ActiveDeadlineSeconds
		if value < 1 || value > math.MaxInt32 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("active_deadline_seconds"), value, validation.InclusiveRangeError(1, math.MaxInt32)))
		}
	}

	// suspended

	return allErrs
}

func validateContainers(containers []v1.Container, volumes map[string]v1.VolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(containers) == 0 {
		return append(allErrs, field.Required(fldPath, ""))
	}

	allNames := sets.String{}
	for i, ctr := range containers {
		idxPath := fldPath.Index(i)
		namePath := idxPath.Child("name")

		// name
		if len(ctr.Name) == 0 {
			allErrs = append(allErrs, field.Required(namePath, ""))
		} else {
			allErrs = append(allErrs, ValidateDNS1123Label(ctr.Name, namePath)...)
		}
		if allNames.Has(ctr.Name) {
			allErrs = append(allErrs, field.Duplicate(namePath, ctr.Name))
		} else {
			allNames.Insert(ctr.Name)
		}
		// image
		if len(ctr.Image) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("image"), ""))
		}
		if len(ctr.Image) != len(strings.TrimSpace(ctr.Image)) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("image"), ctr.Image, "must not have leading or trailing whitespace"))
		}
		// image_pull_policy
		allErrs = append(allErrs, validatePullPolicy(ctr.ImagePullPolicy, idxPath.Child("image_pull_policy"))...)
		// image_deployment_dir
		if len(ctr.ImageDeploymentDir) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("image_deployment_dir"), ""))
		}
		// command todo
		// if len(ctr.Command) == 0 {
		// 	allErrs = append(allErrs, field.Required(idxPath.Child("command"), ""))
		// }
		// termination_command todo

		// security_context
		allErrs = append(allErrs, ValidateSecurityContext(ctr.SecurityContext, fldPath.Child("security_context"))...)
		// env
		allErrs = append(allErrs, ValidateEnv(ctr.Env, idxPath.Child("env"))...)
		// working_dir todo

		// ports
		allErrs = append(allErrs, validateContainerPorts(ctr.Ports, idxPath.Child("ports"))...)
		// volume_mounts
		allErrs = append(allErrs, ValidateVolumeMounts(ctr.VolumeMounts, volumes, idxPath.Child("volume_mounts"))...)
		// resources
		allErrs = append(allErrs, ValidateResourceRequirements(&ctr.Resources, idxPath.Child("resources"))...)
	}
	// Check for colliding ports across all containers.
	allErrs = append(allErrs, checkHostPortConflicts(containers, fldPath)...)

	return allErrs
}

func validateInitContainers(containers, otherContainers []v1.Container, deviceVolumes map[string]v1.VolumeSource, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if len(containers) > 0 {
		allErrs = append(allErrs, validateContainers(containers, deviceVolumes, fldPath)...)
	}

	allNames := sets.String{}
	for _, ctr := range otherContainers {
		allNames.Insert(ctr.Name)
	}
	for i, ctr := range containers {
		idxPath := fldPath.Index(i)
		if allNames.Has(ctr.Name) {
			allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), ctr.Name))
		}
		if len(ctr.Name) > 0 {
			allNames.Insert(ctr.Name)
		}
	}

	return allErrs
}

func validateInstallationContainers(containers, otherContainers []v1.Container, deviceVolumes map[string]v1.VolumeSource, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if len(containers) > 0 {
		allErrs = append(allErrs, validateContainers(containers, deviceVolumes, fldPath)...)
	}

	allNames := sets.String{}
	for _, ctr := range otherContainers {
		allNames.Insert(ctr.Name)
	}
	for i, ctr := range containers {
		idxPath := fldPath.Index(i)
		if allNames.Has(ctr.Name) {
			allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), ctr.Name))
		}
		if len(ctr.Name) > 0 {
			allNames.Insert(ctr.Name)
		}
	}

	return allErrs
}

func validateUninstallationContainers(containers, otherContainers []v1.Container, deviceVolumes map[string]v1.VolumeSource, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if len(containers) > 0 {
		allErrs = append(allErrs, validateContainers(containers, deviceVolumes, fldPath)...)
	}

	allNames := sets.String{}
	for _, ctr := range otherContainers {
		allNames.Insert(ctr.Name)
	}
	for i, ctr := range containers {
		idxPath := fldPath.Index(i)
		if allNames.Has(ctr.Name) {
			allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), ctr.Name))
		}
		if len(ctr.Name) > 0 {
			allNames.Insert(ctr.Name)
		}
	}

	return allErrs
}

// ValidateEnv validates env vars
func ValidateEnv(vars []v1.EnvVar, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, ev := range vars {
		idxPath := fldPath.Index(i)
		if len(ev.Name) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("name"), ""))
		} else {
			for _, msg := range validation.IsEnvVarName(ev.Name) {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), ev.Name, msg))
			}
		}
		allErrs = append(allErrs, validateEnvVarValueFrom(ev, idxPath.Child("value_from"))...)
	}
	return allErrs
}

var supportedPortProtocols = sets.NewString(string(v1.ProtocolTCP), string(v1.ProtocolUDP))

func validateContainerPorts(ports []v1.ContainerPort, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allNames := sets.String{}
	for i, port := range ports {
		idxPath := fldPath.Index(i)
		if len(port.Name) > 0 {
			if msgs := validation.IsValidPortName(port.Name); len(msgs) != 0 {
				for i = range msgs {
					allErrs = append(allErrs, field.Invalid(idxPath.Child("name"), port.Name, msgs[i]))
				}
			} else if allNames.Has(port.Name) {
				allErrs = append(allErrs, field.Duplicate(idxPath.Child("name"), port.Name))
			} else {
				allNames.Insert(port.Name)
			}
		}
		if port.ContainerPort == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("container_port"), ""))
		} else {
			for _, msg := range validation.IsValidPortNum(int(port.ContainerPort)) {
				allErrs = append(allErrs, field.Invalid(idxPath.Child("container_port"), port.ContainerPort, msg))
			}
		}
		if len(port.Protocol) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("protocol"), ""))
		} else if !supportedPortProtocols.Has(string(port.Protocol)) {
			allErrs = append(allErrs, field.NotSupported(idxPath.Child("protocol"), port.Protocol, supportedPortProtocols.List()))
		}
	}
	return allErrs
}

func ValidateVolumeMounts(mounts []v1.VolumeMount, volumes map[string]v1.VolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	mountPoints := sets.NewString()

	for i, mnt := range mounts {
		idxPath := fldPath.Index(i)
		if len(mnt.Name) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("name"), ""))
		}
		if !IsMatchedVolume(mnt.Name, volumes) {
			allErrs = append(allErrs, field.NotFound(idxPath.Child("name"), mnt.Name))
		}
		if len(mnt.MountPath) == 0 {
			allErrs = append(allErrs, field.Required(idxPath.Child("mount_path"), ""))
		}
		if mountPoints.Has(mnt.MountPath) {
			allErrs = append(allErrs, field.Invalid(idxPath.Child("mount_path"), mnt.MountPath, "must be unique"))
		}
		mountPoints.Insert(mnt.MountPath)

		if len(mnt.SubPath) > 0 {
			allErrs = append(allErrs, validateLocalDescendingPath(mnt.SubPath, fldPath.Child("sub_path"))...)
		}

	}
	return allErrs
}

func IsMatchedVolume(name string, volumes map[string]v1.VolumeSource) bool {
	if _, ok := volumes[name]; ok {
		return true
	} else {
		return false
	}
}

func validateEnvVarValueFrom(ev v1.EnvVar, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if ev.ValueFrom == nil {
		return allErrs
	}
	if ev.ValueFrom.FieldRef == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("field_ref"), ""))
	} else {
		allErrs = append(allErrs, validateObjectFieldSelector(ev.ValueFrom.FieldRef, fldPath.Child("field_ref"))...)
	}

	return allErrs
}

func validateObjectFieldSelector(fs *v1.ObjectFieldSelector, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(fs.FieldPath) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("field_path"), ""))
		return allErrs
	}

	return allErrs
}

func GetVolumeMountMap(mounts []v1.VolumeMount) map[string]string {
	volMounts := make(map[string]string)

	for _, mnt := range mounts {
		volMounts[mnt.Name] = mnt.MountPath
	}

	return volMounts
}

func ValidateVolumes(volumes []v1.Volume, fldPath *field.Path) (map[string]v1.VolumeSource, field.ErrorList) {
	allErrs := field.ErrorList{}

	allNames := sets.String{}
	vols := make(map[string]v1.VolumeSource)
	for i, vol := range volumes {
		idxPath := fldPath.Index(i)
		namePath := idxPath.Child("name")
		el := validateVolumeSource(&vol.VolumeSource, idxPath, vol.Name)
		if len(vol.Name) == 0 {
			el = append(el, field.Required(namePath, ""))
		} else {
			el = append(el, ValidateDNS1123Label(vol.Name, namePath)...)
		}
		if allNames.Has(vol.Name) {
			el = append(el, field.Duplicate(namePath, vol.Name))
		}
		if len(el) == 0 {
			allNames.Insert(vol.Name)
			vols[vol.Name] = vol.VolumeSource
		} else {
			allErrs = append(allErrs, el...)
		}

	}
	return vols, allErrs
}

func validateVolumeSource(source *v1.VolumeSource, fldPath *field.Path, volName string) field.ErrorList {
	allErrs := field.ErrorList{}
	if source.ConfigMap != nil {
		allErrs = append(allErrs, validateConfigMapVolumeSource(source.ConfigMap, fldPath.Child("config_map"))...)
	}
	return allErrs
}

func validateConfigMapVolumeSource(configMapSource *v1.ConfigMapVolumeSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(configMapSource.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	}

	configMapMode := configMapSource.DefaultMode
	if configMapMode != nil && (*configMapMode > 0777 || *configMapMode < 0) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("default_mode"), *configMapMode, fileModeErrorMsg))
	}

	itemsPath := fldPath.Child("items")
	for i, kp := range configMapSource.Items {
		itemPath := itemsPath.Index(i)
		allErrs = append(allErrs, validateKeyToPath(&kp, itemPath)...)
	}
	return allErrs
}

func validateKeyToPath(kp *v1.KeyToPath, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(kp.Key) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("key"), ""))
	}
	if len(kp.Path) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("path"), ""))
	}
	allErrs = append(allErrs, validateLocalNonReservedPath(kp.Path, fldPath.Child("path"))...)
	if kp.Mode != nil && (*kp.Mode > 0777 || *kp.Mode < 0) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("mode"), *kp.Mode, fileModeErrorMsg))
	}

	return allErrs
}

func validateRestartPolicy(restartPolicy *v1.RestartPolicy, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}
	switch *restartPolicy {
	case v1.RestartPolicyAlways, v1.RestartPolicyOnFailure, v1.RestartPolicyNever:
		break
	case "":
		allErrors = append(allErrors, field.Required(fldPath, ""))
	default:
		validValues := []string{string(v1.RestartPolicyAlways), string(v1.RestartPolicyOnFailure), string(v1.RestartPolicyNever)}
		allErrors = append(allErrors, field.NotSupported(fldPath, *restartPolicy, validValues))
	}

	return allErrors
}

// validateAffinity checks if given affinities are valid
func validateAffinity(affinity *v1.Affinity, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if affinity != nil {
		if affinity.PodAffinity != nil {
			allErrs = append(allErrs, validatePodAffinity(affinity.PodAffinity, fldPath.Child("pod_affinity"))...)
		}
		if affinity.PodAntiAffinity != nil {
			allErrs = append(allErrs, validatePodAntiAffinity(affinity.PodAntiAffinity, fldPath.Child("pod_anti_affinity"))...)
		}
	}

	return allErrs
}

func validatePodAffinity(podAffinity *v1.PodAffinity, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if podAffinity.Required != nil {
		allErrs = append(allErrs, ValidateLabels(podAffinity.Required.MatchLabels, fldPath.Child("required").Child("match_labels"))...)
	}

	return allErrs
}

func validatePodAntiAffinity(podAntiAffinity *v1.PodAntiAffinity, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if podAntiAffinity.Required != nil {
		allErrs = append(allErrs, ValidateLabels(podAntiAffinity.Required.MatchLabels, fldPath.Child("required").Child("match_labels"))...)
	}
	return allErrs
}

var supportedPullPolicies = sets.NewString(string(v1.PullAlways), string(v1.PullIfNotPresent), string(v1.PullNever))

func validatePullPolicy(policy v1.PullPolicy, fldPath *field.Path) field.ErrorList {
	allErrors := field.ErrorList{}

	switch policy {
	case v1.PullAlways, v1.PullIfNotPresent, v1.PullNever:
		break
	case "":
		allErrors = append(allErrors, field.Required(fldPath, ""))
	default:
		allErrors = append(allErrors, field.NotSupported(fldPath, policy, supportedPullPolicies.List()))
	}

	return allErrors
}

// Validate compute resource typename.
func validateResourceName(value string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	switch v1.ResourceName(value) {
	case v1.ResourceCPU, v1.ResourceMemory:
		break
	default:
		allErrs = append(allErrs, field.NotSupported(fldPath, value, []string{string(v1.ResourceMemory), string(v1.ResourceCPU)}))
	}

	return allErrs
}

func validateBasicResource(quantity resource.Quantity, fldPath *field.Path) field.ErrorList {
	if quantity.Value() < 0 {
		return field.ErrorList{field.Invalid(fldPath, quantity.Value(), "must be a valid resource quantity")}
	}
	return field.ErrorList{}
}

func ValidateResourceRequirements(requirements *v1.ResourceRequirements, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	limPath := fldPath.Child("limits")
	reqPath := fldPath.Child("requests")

	for resourceName, quantity := range requirements.Limits {
		fldPath := limPath.Key(string(resourceName))
		// Validate resource name.
		allErrs = append(allErrs, validateResourceName(string(resourceName), fldPath)...)
		allErrs = append(allErrs, validateBasicResource(quantity, fldPath)...)

	}
	for resourceName, quantity := range requirements.Requests {
		fldPath := reqPath.Key(string(resourceName))
		// Validate resource name.
		allErrs = append(allErrs, validateResourceName(string(resourceName), fldPath)...)
		allErrs = append(allErrs, validateBasicResource(quantity, fldPath)...)
	}
	return allErrs
}

// AccumulateUniqueHostPorts extracts each HostPort of each Container,
// accumulating the results and returning an error if any ports conflict.
func AccumulateUniqueHostPorts(containers []v1.Container, accumulator *sets.String, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for ci, ctr := range containers {
		idxPath := fldPath.Index(ci)
		portsPath := idxPath.Child("ports")
		for pi := range ctr.Ports {
			idxPath := portsPath.Index(pi)
			port := ctr.Ports[pi].ContainerPort
			if port == 0 {
				continue
			}
			str := fmt.Sprintf("%s/%d", ctr.Ports[pi].Protocol, port)
			if accumulator.Has(str) {
				allErrs = append(allErrs, field.Duplicate(idxPath.Child("hostPort"), str))
			} else {
				accumulator.Insert(str)
			}
		}
	}
	return allErrs
}

// checkHostPortConflicts checks for colliding Port.HostPort values across
// a slice of containers.
func checkHostPortConflicts(containers []v1.Container, fldPath *field.Path) field.ErrorList {
	allPorts := sets.String{}
	return AccumulateUniqueHostPorts(containers, &allPorts, fldPath)
}

// ValidatePodSecurityContext test that the specified PodSecurityContext has valid data.
func ValidatePodSecurityContext(securityContext *v1.PodSecurityContext, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if securityContext != nil {
		if securityContext.RunAsNonRoot != nil {
			// todo
		}
	}

	return allErrs
}

// ValidateSecurityContext ensure the security context contains valid settings
func ValidateSecurityContext(sc *v1.SecurityContext, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	// this should only be true for testing since SecurityContext is defaulted by the core
	if sc == nil {
		return allErrs
	}

	if len(sc.RunAsUser) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("run_as_user"), ""))
	}

	return allErrs
}

func ValidatePodUpdate(newPod, oldPod *v1.Pod) field.ErrorList {
	fldPath := field.NewPath("metadata")
	allErrs := ValidateObjectMetaUpdate(&newPod.ObjectMeta, &oldPod.ObjectMeta, fldPath)
	allErrs = append(allErrs, ValidatePodSpecificAnnotationUpdates(newPod, oldPod, fldPath.Child("annotations"))...)
	specPath := field.NewPath("spec")
	// validate updatable fields:
	// 1.  spec.containers[*].image
	// 2.  spec.initContainers[*].image
	// 3.  spec.InstallationContainers[*].image
	// 4.  spec.UninstallationContainers[*].image
	// 5.  spec.activeDeadlineSeconds
	// 6.  spec.terminationGracePeriodSeconds
	containerErrs, stop := ValidateContainerUpdates(newPod.Spec.Containers, oldPod.Spec.Containers, specPath.Child("containers"))
	allErrs = append(allErrs, containerErrs...)
	if stop {
		return allErrs
	}
	containerErrs, stop = ValidateContainerUpdates(newPod.Spec.InitContainers, oldPod.Spec.InitContainers, specPath.Child("init_containers"))
	allErrs = append(allErrs, containerErrs...)
	if stop {
		return allErrs
	}
	containerErrs, stop = ValidateContainerUpdates(newPod.Spec.InstallationContainers, oldPod.Spec.InstallationContainers, specPath.Child("installation_containers"))
	allErrs = append(allErrs, containerErrs...)
	if stop {
		return allErrs
	}
	containerErrs, stop = ValidateContainerUpdates(newPod.Spec.UninstallationContainers, oldPod.Spec.UninstallationContainers, specPath.Child("uninstallation_containers"))
	allErrs = append(allErrs, containerErrs...)
	if stop {
		return allErrs
	}

	// validate updated spec.activeDeadlineSeconds.  two types of updates are allowed:
	// 1.  from nil to a positive value
	// 2.  from a positive value to a lesser, non-negative value
	if newPod.Spec.ActiveDeadlineSeconds != nil {
		newActiveDeadlineSeconds := *newPod.Spec.ActiveDeadlineSeconds
		if newActiveDeadlineSeconds < 0 || newActiveDeadlineSeconds > math.MaxInt32 {
			allErrs = append(allErrs, field.Invalid(specPath.Child("active_deadline_seconds"), newActiveDeadlineSeconds, validation.InclusiveRangeError(0, math.MaxInt32)))
			return allErrs
		}
		if oldPod.Spec.ActiveDeadlineSeconds != nil {
			oldActiveDeadlineSeconds := *oldPod.Spec.ActiveDeadlineSeconds
			if oldActiveDeadlineSeconds < newActiveDeadlineSeconds {
				allErrs = append(allErrs, field.Invalid(specPath.Child("active_deadline_seconds"), newActiveDeadlineSeconds, "must be less than or equal to previous value"))
				return allErrs
			}
		}
	} else if oldPod.Spec.ActiveDeadlineSeconds != nil {
		allErrs = append(allErrs, field.Invalid(specPath.Child("active_deadline_seconds"), newPod.Spec.ActiveDeadlineSeconds, "must not update from a positive integer to nil value"))
	}

	// the last thing to check is pod spec equality.  If the pod specs are equal, then we can simply return the errors we have
	// so far and save the cost of a deep copy.
	if reflect.DeepEqual(newPod.Spec, oldPod.Spec) {
		return allErrs
	}

	// handle updatable fields by munging those fields prior to deep equal comparison.
	mungedPodSpec := *newPod.Spec.DeepCopy()
	// munged spec.containers[*].image
	var newContainers []v1.Container
	for ix, container := range mungedPodSpec.Containers {
		container.Image = oldPod.Spec.Containers[ix].Image
		newContainers = append(newContainers, container)
	}
	mungedPodSpec.Containers = newContainers
	// munged spec.initContainers[*].image
	var newInitContainers []v1.Container
	for ix, container := range mungedPodSpec.InitContainers {
		container.Image = oldPod.Spec.InitContainers[ix].Image
		newInitContainers = append(newInitContainers, container)
	}
	mungedPodSpec.InitContainers = newInitContainers
	// munged spec.initContainers[*].image
	var newInstallationContainers []v1.Container
	for ix, container := range mungedPodSpec.InstallationContainers {
		container.Image = oldPod.Spec.InstallationContainers[ix].Image
		newInstallationContainers = append(newInstallationContainers, container)
	}
	mungedPodSpec.InstallationContainers = newInstallationContainers
	// munged spec.initContainers[*].image
	var newUninstallationContainers []v1.Container
	for ix, container := range mungedPodSpec.UninstallationContainers {
		container.Image = oldPod.Spec.UninstallationContainers[ix].Image
		newUninstallationContainers = append(newUninstallationContainers, container)
	}
	mungedPodSpec.UninstallationContainers = newUninstallationContainers

	// munged spec.activeDeadlineSeconds
	mungedPodSpec.ActiveDeadlineSeconds = nil
	if oldPod.Spec.ActiveDeadlineSeconds != nil {
		activeDeadlineSeconds := *oldPod.Spec.ActiveDeadlineSeconds
		mungedPodSpec.ActiveDeadlineSeconds = &activeDeadlineSeconds
	}

	// Relax validation of immutable fields to allow it to be set to 1 if it was previously negative.
	if oldPod.Spec.TerminationGracePeriodSeconds != nil && *oldPod.Spec.TerminationGracePeriodSeconds < 0 &&
		mungedPodSpec.TerminationGracePeriodSeconds != nil && *mungedPodSpec.TerminationGracePeriodSeconds == 1 {
		mungedPodSpec.TerminationGracePeriodSeconds = oldPod.Spec.TerminationGracePeriodSeconds
	}

	if !reflect.DeepEqual(mungedPodSpec, oldPod.Spec) {
		specDiff := cmp.Diff(oldPod.Spec, mungedPodSpec)
		allErrs = append(allErrs, field.Forbidden(specPath, fmt.Sprintf("pod updates may not change fields other than `spec.containers[*].image`, `spec.init_containers[*].image`, `spec.active_deadline_seconds` or `spec.termination_grace_period_seconds` (allow it to be set to 1 if it was previously negative)\n%v", specDiff)))
	}

	return allErrs
}

func ValidatePodSpecificAnnotationUpdates(newPod, oldPod *v1.Pod, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	newAnnotations := newPod.Annotations
	oldAnnotations := oldPod.Annotations
	for k, oldVal := range oldAnnotations {
		if newVal, exists := newAnnotations[k]; exists && newVal == oldVal {
			continue // No change.
		}
		// special keys
	}
	// Check for additions
	for k := range newAnnotations {
		if _, ok := oldAnnotations[k]; ok {
			continue // No change.
		}
		// special keys
	}
	allErrs = append(allErrs, ValidatePodSpecificAnnotations(newAnnotations, &newPod.Spec, fldPath)...)
	return allErrs
}

func ValidateContainerUpdates(newContainers, oldContainers []v1.Container, fldPath *field.Path) (allErrs field.ErrorList, stop bool) {
	allErrs = field.ErrorList{}
	if len(newContainers) != len(oldContainers) {
		allErrs = append(allErrs, field.Forbidden(fldPath, "pod updates may not add or remove containers"))
		return allErrs, true
	}

	// validate updated container images
	for i, ctr := range newContainers {
		if len(ctr.Image) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Index(i).Child("image"), ""))
		}
		// this is only called from ValidatePodUpdate so its safe to check leading/trailing whitespace.
		if len(strings.TrimSpace(ctr.Image)) != len(ctr.Image) {
			allErrs = append(allErrs, field.Invalid(fldPath.Index(i).Child("image"), ctr.Image, "must not have leading or trailing whitespace"))
		}
	}
	return allErrs, false
}
