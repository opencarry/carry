package validation

import (
	"github.com/opencarry/carry/pkg/api/meta"
	metav1 "github.com/opencarry/carry/pkg/apis/carry.i/v1"
	field "github.com/opencarry/carry/pkg/util/validation/field"
)

type ValidateNameFunc func(name string, prefix bool) []string

func ValidateObjectMeta(objMeta *metav1.ObjectMeta, requiresNamespace bool, nameFn ValidateNameFunc, fldPath *field.Path) field.ErrorList {
	metadata, err := meta.Accessor(objMeta)
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

	return allErrs
}
