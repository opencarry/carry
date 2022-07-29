package validation

import (
	"fmt"
	"reflect"
	"strings"

	metav1 "github.com/opencarry/carry/pkg/apis/carry.i/v1"
	"github.com/opencarry/carry/pkg/util/validation"
	"github.com/opencarry/carry/pkg/util/validation/field"
)

func ValidateObjectMeta(objMeta *metav1.ObjectMeta, requiresNamespace bool, nameFn ValidateNameFunc, fldPath *field.Path) field.ErrorList {
	metadata, err := metav1.Accessor(objMeta)
	if err != nil {
		allErrs := field.ErrorList{}
		allErrs = append(allErrs, field.Invalid(fldPath, objMeta, err.Error()))
		return allErrs
	}
	return ValidateObjectMetaAccessor(metadata, requiresNamespace, nameFn, fldPath)
}

func ValidateObjectMetaAccessor(meta metav1.Object, requiresNamespace bool, nameFn ValidateNameFunc, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(meta.GetGenerateName()) != 0 {
		for _, msg := range nameFn(meta.GetGenerateName(), true) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("generate_name"), meta.GetGenerateName(), msg))
		}
	}

	if len(meta.GetName()) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "name or generate_name is required"))
	} else {
		for _, msg := range nameFn(meta.GetName(), true) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), meta.GetName(), msg))
		}
	}

	if requiresNamespace {
		if len(meta.GetNamespace()) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("namespace"), ""))
		} else {
			for _, msg := range ValidateNamespaceName(meta.GetNamespace(), false) {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("namespace"), meta.GetNamespace(), msg))
			}
		}
	} else {
		if len(meta.GetNamespace()) == 0 {
			allErrs = append(allErrs, field.Required(fldPath.Child("namespace"), "not allow on this type"))
		}
	}

	allErrs = append(allErrs, ValidateLabels(meta.GetLabels(), fldPath.Child("labels"))...)
	allErrs = append(allErrs, ValidateAnnotations(meta.GetAnnotations(), fldPath.Child("annotations"))...)
	allErrs = append(allErrs, ValidateOwnerReferences(meta.GetOwnerReferences(), fldPath.Child("owner_references"))...)
	return allErrs
}

const FieldImmutableErrorMsg string = `field is immutable`

const totalAnnotationSizeLimitB int = 256 * (1 << 10) // 256 kB

// BannedOwners is a black list of object that are not allowed to be owners.
var BannedOwners = map[string]struct{}{
	"event": {},
}

// ValidateAnnotations validates that a set of annotations are correctly defined.
func ValidateAnnotations(annotations map[string]string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	var totalSize int64
	for k, v := range annotations {
		for _, msg := range validation.IsQualifiedName(strings.ToLower(k)) {
			allErrs = append(allErrs, field.Invalid(fldPath, k, msg))
		}
		totalSize += (int64)(len(k)) + (int64)(len(v))
	}
	if totalSize > (int64)(totalAnnotationSizeLimitB) {
		allErrs = append(allErrs, field.TooLong(fldPath, "", totalAnnotationSizeLimitB))
	}
	return allErrs
}

func validateOwnerReference(ownerReference metav1.OwnerReference, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(ownerReference.APIVersion) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVersion"), ownerReference.APIVersion, "version must not be empty"))
	}
	if len(ownerReference.Kind) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("kind"), ownerReference.Kind, "kind must not be empty"))
	}
	if len(ownerReference.Name) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), ownerReference.Name, "name must not be empty"))
	}
	if len(ownerReference.UID) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("uid"), ownerReference.UID, "uid must not be empty"))
	}
	if _, ok := BannedOwners[ownerReference.Kind]; ok {
		allErrs = append(allErrs, field.Invalid(fldPath, ownerReference, fmt.Sprintf("%s is disallowed from being an owner", ownerReference.Kind)))
	}
	return allErrs
}

func ValidateOwnerReferences(ownerReferences []metav1.OwnerReference, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	controllerName := ""
	for _, ref := range ownerReferences {
		allErrs = append(allErrs, validateOwnerReference(ref, fldPath)...)
		if ref.Controller != nil && *ref.Controller {
			if controllerName != "" {
				allErrs = append(allErrs, field.Invalid(fldPath, ownerReferences,
					fmt.Sprintf("Only one reference can have Controller set to true. Found \"true\" in references for %v and %v", controllerName, ref.Name)))
			} else {
				controllerName = ref.Name
			}
		}
	}
	return allErrs
}

func ValidateObjectMetaUpdate(newMeta, oldMeta *metav1.ObjectMeta, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	newMetadata, err := metav1.Accessor(newMeta)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, newMeta, err.Error()))
		return allErrs
	}
	oldMetadata, err := metav1.Accessor(oldMeta)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath, oldMeta, err.Error()))
		return allErrs
	}

	return ValidateObjectMetaAccessorUpdate(newMetadata, oldMetadata, fldPath)
}

// ValidateObjectMetaAccessorUpdate validates an object's metadata when updated.
func ValidateObjectMetaAccessorUpdate(newMeta, oldMeta metav1.Object, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if !oldMeta.GetDeletionTime().IsZero() {
		// 已删除状态下，某些字段不允许更新
		// return
	}

	// Reject updates that don't specify a resource version
	if len(newMeta.GetResourceVersion()) == 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resource_version"), newMeta.GetResourceVersion(), "must be specified for an update"))
	}

	// if newMeta.GetUID() != oldMeta.GetUID() {
	// 	allErrs = append(allErrs, field.Invalid(fldPath.Child("uid"), newMeta.GetUID(), FieldImmutableErrorMsg))
	// }
	allErrs = append(allErrs, ValidateImmutableField(newMeta.GetName(), oldMeta.GetName(), fldPath.Child("name"))...)
	allErrs = append(allErrs, ValidateImmutableField(newMeta.GetNamespace(), oldMeta.GetNamespace(), fldPath.Child("namespace"))...)
	allErrs = append(allErrs, ValidateImmutableField(newMeta.GetUID(), oldMeta.GetUID(), fldPath.Child("uid"))...)
	allErrs = append(allErrs, ValidateImmutableField(newMeta.GetCreationTime(), oldMeta.GetCreationTime(), fldPath.Child("creation_time"))...)
	allErrs = append(allErrs, ValidateImmutableField(newMeta.GetDeletionTime(), oldMeta.GetDeletionTime(), fldPath.Child("deletion_time"))...)
	allErrs = append(allErrs, ValidateImmutableField(newMeta.GetDeletionGracePeriodSeconds(), oldMeta.GetDeletionGracePeriodSeconds(), fldPath.Child("deletion_grace_period_seconds"))...)

	allErrs = append(allErrs, ValidateLabels(newMeta.GetLabels(), fldPath.Child("labels"))...)
	allErrs = append(allErrs, ValidateAnnotations(newMeta.GetAnnotations(), fldPath.Child("annotations"))...)
	allErrs = append(allErrs, ValidateOwnerReferences(newMeta.GetOwnerReferences(), fldPath.Child("owner_references"))...)

	return allErrs
}

// ValidateImmutableField validates the new value and the old value are deeply equal.
func ValidateImmutableField(newVal, oldVal interface{}, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	// todo
	// func(a, b resource.Quantity) bool
	// func(a, b labels.Selector) bool {
	// func(a, b fields.Selector) bool {
	if !reflect.DeepEqual(oldVal, newVal) {
		allErrs = append(allErrs, field.Invalid(fldPath, newVal, FieldImmutableErrorMsg))
	}
	return allErrs
}
