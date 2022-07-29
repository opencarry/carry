package validation

import (
	"reflect"

	v1 "github.com/opencarry/carry/pkg/apis/carry.i/v1"
	"github.com/opencarry/carry/pkg/util/validation/field"
)

func ValidateStatefulSetName(name string, prefix bool) []string {
	return NameIsDNSSubdomain(name, prefix)
}

func ValidatePodTemplateSpecForStatefulSet(template *v1.PodTemplateSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if template == nil {
		allErrs = append(allErrs, field.Required(fldPath, ""))
	} else {
		allErrs = append(allErrs, ValidateLabels(template.Labels, fldPath.Child("labels"))...)
		allErrs = append(allErrs, ValidateAnnotations(template.Annotations, fldPath.Child("annotations"))...)
		allErrs = append(allErrs, ValidatePodSpecificAnnotations(template.Annotations, &template.Spec, fldPath.Child("annotations"))...)
	}
	return allErrs
}

// ValidateStatefulSetSpec tests if required fields in the StatefulSet spec are set.
func ValidateStatefulSetSpec(spec *v1.StatefulSetSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateNonnegativeField(int64(*spec.Replicas), fldPath.Child("replicas"))...)
	if spec.Selector == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("selector"), ""))
	} else {
		allErrs = append(allErrs, ValidateLabelSelector(spec.Selector, fldPath.Child("selector"))...)
		if len(spec.Selector.MatchLabels)+len(spec.Selector.MatchExpressions) == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("selector"), spec.Selector, "empty selector is not valid for statefulset."))
		}
	}

	if spec.Template.Spec.RestartPolicy != v1.RestartPolicyAlways {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("template", "spec", "restartPolicy"), spec.Template.Spec.RestartPolicy, []string{string(v1.RestartPolicyAlways)}))
	}

	return allErrs
}

// ValidateStatefulSet validates a StatefulSet.
func ValidateStatefulSet(statefulSet *v1.StatefulSet) field.ErrorList {
	allErrs := ValidateObjectMeta(&statefulSet.ObjectMeta, true, ValidateStatefulSetName, field.NewPath("metadata"))
	allErrs = append(allErrs, ValidateStatefulSetSpec(&statefulSet.Spec, field.NewPath("spec"))...)
	return allErrs
}

// ValidateStatefulSetUpdate tests if required fields in the StatefulSet are set.
func ValidateStatefulSetUpdate(statefulSet, oldStatefulSet *v1.StatefulSet) field.ErrorList {
	allErrs := ValidateObjectMetaUpdate(&statefulSet.ObjectMeta, &oldStatefulSet.ObjectMeta, field.NewPath("metadata"))

	// TODO: For now we're taking the safe route and disallowing all updates to
	// spec except for Replicas, for scaling, and Template.Spec.containers.image
	// for rolling-update. Enable others on a case by case basis.
	restoreReplicas := statefulSet.Spec.Replicas
	statefulSet.Spec.Replicas = oldStatefulSet.Spec.Replicas

	restoreContainers := statefulSet.Spec.Template.Spec.Containers
	statefulSet.Spec.Template.Spec.Containers = oldStatefulSet.Spec.Template.Spec.Containers

	if !reflect.DeepEqual(statefulSet.Spec, oldStatefulSet.Spec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "updates to statefulset spec for fields other than 'replicas' and 'containers' are forbidden."))
	}
	statefulSet.Spec.Replicas = restoreReplicas
	statefulSet.Spec.Template.Spec.Containers = restoreContainers

	allErrs = append(allErrs, ValidateNonnegativeField(int64(*statefulSet.Spec.Replicas), field.NewPath("spec", "replicas"))...)
	containerErrs, _ := ValidateContainerUpdates(statefulSet.Spec.Template.Spec.Containers, oldStatefulSet.Spec.Template.Spec.Containers, field.NewPath("spec").Child("template").Child("containers"))
	allErrs = append(allErrs, containerErrs...)
	return allErrs
}

// ValidateStatefulSetStatusUpdate tests if required fields in the StatefulSet are set.
func ValidateStatefulSetStatusUpdate(statefulSet, oldStatefulSet *v1.StatefulSet) field.ErrorList {
	allErrs := field.ErrorList{}
	allErrs = append(allErrs, ValidateObjectMetaUpdate(&statefulSet.ObjectMeta, &oldStatefulSet.ObjectMeta, field.NewPath("metadata"))...)
	// TODO: Validate status.
	return allErrs
}
