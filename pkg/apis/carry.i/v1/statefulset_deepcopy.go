package v1

func (in *StatefulSet) DeepCopy() *StatefulSet {
	if in == nil {
		return nil
	}
	out := new(StatefulSet)
	in.DeepCopyInto(out)
	return out
}

func (in *StatefulSet) DeepCopyInto(out *StatefulSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

func (in *StatefulSetSpec) DeepCopy() *StatefulSetSpec {
	if in == nil {
		return nil
	}
	out := new(StatefulSetSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *StatefulSetSpec) DeepCopyInto(out *StatefulSetSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		if *in == nil {
			*out = nil
		} else {
			*out = new(int64)
			**out = **in
		}
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		if *in == nil {
			*out = nil
		} else {
			*out = new(LabelSelector)
			(*in).DeepCopyInto(*out)
		}
	}
	in.Template.DeepCopyInto(&out.Template)
	if in.RevisionHistoryLimit != nil {
		in, out := &in.RevisionHistoryLimit, &out.RevisionHistoryLimit
		if *in == nil {
			*out = nil
		} else {
			*out = new(int64)
			**out = **in
		}
	}
	return
}

func (in *StatefulSetStatus) DeepCopy() *StatefulSetStatus {
	if in == nil {
		return nil
	}
	out := new(StatefulSetStatus)
	in.DeepCopyInto(out)
	return out
}

func (in *StatefulSetStatus) DeepCopyInto(out *StatefulSetStatus) {
	*out = *in
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		if *in == nil {
			*out = nil
		} else {
			*out = new(int64)
			**out = **in
		}
	}
	if in.CollisionCount != nil {
		in, out := &in.CollisionCount, &out.CollisionCount
		if *in == nil {
			*out = nil
		} else {
			*out = new(int64)
			**out = **in
		}
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]StatefulSetCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *StatefulSetCondition) DeepCopy() *StatefulSetCondition {
	if in == nil {
		return nil
	}
	out := new(StatefulSetCondition)
	in.DeepCopyInto(out)
	return out
}

func (in *StatefulSetCondition) DeepCopyInto(out *StatefulSetCondition) {
	*out = *in
	return
}
