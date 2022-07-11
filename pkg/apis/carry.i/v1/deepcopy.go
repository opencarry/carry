package v1

// DeepCopy copying the receiver, creating a new PodSpec.
func (in *PodSpec) DeepCopy() *PodSpec {
	if in == nil {
		return nil
	}
	out := new(PodSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto copying the receiver, writing into out. in must be non-nil.
func (in *PodSpec) DeepCopyInto(out *PodSpec) {
	*out = *in
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TerminationGracePeriodSeconds != nil {
		in, out := &in.TerminationGracePeriodSeconds, &out.TerminationGracePeriodSeconds
		*out = new(int64)
		**out = **in
	}
	if in.ActiveDeadlineSeconds != nil {
		in, out := &in.ActiveDeadlineSeconds, &out.ActiveDeadlineSeconds
		*out = new(int64)
		**out = **in
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}

	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(Affinity)
		(*in).DeepCopyInto(*out)
	}
	return
}

func (in *Affinity) DeepCopy() *Affinity {
	if in == nil {
		return nil
	}
	out := new(Affinity)
	in.DeepCopyInto(out)
	return out
}

func (in *Affinity) DeepCopyInto(out *Affinity) {
	*out = *in
	if in.PodAffinity != nil {
		in, out := &in.PodAffinity, &out.PodAffinity
		*out = new(PodAffinity)
		(*in).DeepCopyInto(*out)
	}
	if in.PodAntiAffinity != nil {
		in, out := &in.PodAntiAffinity, &out.PodAntiAffinity
		*out = new(PodAntiAffinity)
		(*in).DeepCopyInto(*out)
	}
	return
}

func (in *PodAffinity) DeepCopy() *PodAffinity {
	if in == nil {
		return nil
	}
	out := new(PodAffinity)
	in.DeepCopyInto(out)
	return out
}

func (in *PodAffinity) DeepCopyInto(out *PodAffinity) {
	*out = *in
	if in.Required != nil {
		in, out := &in.Required, &out.Required
		*out = new(Required)
		(*in).DeepCopyInto(*out)
	}
}

func (in *PodAntiAffinity) DeepCopy() *PodAntiAffinity {
	if in == nil {
		return nil
	}
	out := new(PodAntiAffinity)
	in.DeepCopyInto(out)
	return out
}

func (in *PodAntiAffinity) DeepCopyInto(out *PodAntiAffinity) {
	*out = *in
	if in.Required != nil {
		in, out := &in.Required, &out.Required
		*out = new(Required)
		(*in).DeepCopyInto(*out)
	}
}

func (in *Required) DeepCopy() *Required {
	if in == nil {
		return nil
	}
	out := new(Required)
	in.DeepCopyInto(out)
	return out
}

func (in *Required) DeepCopyInto(out *Required) {
	*out = *in
	if in.MatchLabels != nil {
		in, out := &in.MatchLabels, &out.MatchLabels
		*out = make(map[string]string)
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

func (in *Container) DeepCopyInto(out *Container) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]ContainerPort, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(SecurityContext)
		(*in).DeepCopyInto(*out)
	}
	return
}

func (in *Container) DeepCopy() *Container {
	if in == nil {
		return nil
	}
	out := new(Container)
	in.DeepCopyInto(out)
	return out
}

func (in *Volume) DeepCopy() *Volume {
	if in == nil {
		return nil
	}
	out := new(Volume)
	in.DeepCopyInto(out)
	return out
}

func (in *Volume) DeepCopyInto(out *Volume) {
	*out = *in
	in.VolumeSource.DeepCopyInto(&out.VolumeSource)
	return
}

func (in *VolumeSource) DeepCopy() *VolumeSource {
	if in == nil {
		return nil
	}
	out := new(VolumeSource)
	in.DeepCopyInto(out)
	return out
}

func (in *VolumeSource) DeepCopyInto(out *VolumeSource) {
	*out = *in
	if in.ConfigMap != nil {
		in, out := &in.ConfigMap, &out.ConfigMap
		*out = new(ConfigMapVolumeSource)
		(*in).DeepCopyInto(*out)
	}

	return
}

func (in *ConfigMapVolumeSource) DeepCopy() *ConfigMapVolumeSource {
	if in == nil {
		return nil
	}
	out := new(ConfigMapVolumeSource)
	in.DeepCopyInto(out)
	return out
}

func (in *ConfigMapVolumeSource) DeepCopyInto(out *ConfigMapVolumeSource) {
	*out = *in
	out.LocalObjectReference = in.LocalObjectReference
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KeyToPath, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.DefaultMode != nil {
		in, out := &in.DefaultMode, &out.DefaultMode
		*out = new(int64)
		**out = **in
	}
	if in.Optional != nil {
		in, out := &in.Optional, &out.Optional
		*out = new(bool)
		**out = **in
	}
	return
}

func (in *KeyToPath) DeepCopyInto(out *KeyToPath) {
	*out = *in
	if in.Mode != nil {
		in, out := &in.Mode, &out.Mode
		*out = new(int32)
		**out = **in
	}
	return
}

func (in *KeyToPath) DeepCopy() *KeyToPath {
	if in == nil {
		return nil
	}
	out := new(KeyToPath)
	in.DeepCopyInto(out)
	return out
}

func (in *SecurityContext) DeepCopyInto(out *SecurityContext) {
	*out = *in
	return
}

func (in *SecurityContext) DeepCopy() *SecurityContext {
	if in == nil {
		return nil
	}
	out := new(SecurityContext)
	in.DeepCopyInto(out)
	return out
}

func (in *EnvVar) DeepCopy() *EnvVar {
	if in == nil {
		return nil
	}
	out := new(EnvVar)
	in.DeepCopyInto(out)
	return out
}

func (in *EnvVar) DeepCopyInto(out *EnvVar) {
	*out = *in
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(EnvVarSource)
		(*in).DeepCopyInto(*out)
	}
	return
}

func (in *EnvVarSource) DeepCopyInto(out *EnvVarSource) {
	*out = *in
	if in.FieldRef != nil {
		in, out := &in.FieldRef, &out.FieldRef
		*out = new(ObjectFieldSelector)
		**out = **in
	}
	return
}

func (in *EnvVarSource) DeepCopy() *EnvVarSource {
	if in == nil {
		return nil
	}
	out := new(EnvVarSource)
	in.DeepCopyInto(out)
	return out
}

func (in *PodSecurityContext) DeepCopy() *PodSecurityContext {
	if in == nil {
		return nil
	}
	out := new(PodSecurityContext)
	in.DeepCopyInto(out)
	return out
}

func (in *PodSecurityContext) DeepCopyInto(out *PodSecurityContext) {
	*out = *in
	if in.RunAsNonRoot != nil {
		in, out := &in.RunAsNonRoot, &out.RunAsNonRoot
		*out = new(bool)
		**out = **in
	}
	return
}

func (in *ResourceRequirements) DeepCopy() *ResourceRequirements {
	if in == nil {
		return nil
	}
	out := new(ResourceRequirements)
	in.DeepCopyInto(out)
	return out
}

func (in *ResourceRequirements) DeepCopyInto(out *ResourceRequirements) {
	*out = *in
	if in.Limits != nil {
		in, out := &in.Limits, &out.Limits
		*out = make(ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	if in.Requests != nil {
		in, out := &in.Requests, &out.Requests
		*out = make(ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	return
}

func (in *VolumeMount) DeepCopy() *VolumeMount {
	if in == nil {
		return nil
	}
	out := new(VolumeMount)
	in.DeepCopyInto(out)
	return out
}

func (in *VolumeMount) DeepCopyInto(out *VolumeMount) {
	*out = *in
	return
}