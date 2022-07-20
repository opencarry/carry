package v1

func IsControlledBy(obj Object, owner Object) bool {
	ref := GetControllerOf(obj)
	if ref == nil {
		return false
	}
	return ref.UID == owner.GetUID()
}

// GetControllerOf
// @Description:
// Parameters:
// @param Controllee 受控
// @return *OwnerReference
func GetControllerOf(Controllee Object) *OwnerReference {
	for _, ref := range Controllee.GetOwnerReferences() {
		if ref.Controller != nil && *ref.Controller {
			return &ref
		}
	}
	return nil
}

func NewControllerRef(owner Object, groupVersion, kind string) *OwnerReference {
	isController := true
	return &OwnerReference{
		APIVersion: groupVersion,
		Kind:       kind,
		Name:       owner.GetName(),
		UID:        owner.GetUID(),
		Controller: &isController,
	}
}
