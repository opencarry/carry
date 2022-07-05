package meta

import (
	"fmt"

	metav1 "github.com/opencarry/carry/pkg/apis/carry.i/v1"
)

var errNotList = fmt.Errorf("object does not implement the List interfaces")

func ListAccessor(obj interface{}) (metav1.ListInterface, error) {
	switch t := obj.(type) {
	case metav1.ListInterface:
		return t, nil
	case metav1.ListMetaAccessor:
		if m := t.GetListMeta(); m != nil {
			return m, nil
		}
		return nil, errNotList
	default:
		return nil, errNotList
	}
}

var errNotObject = fmt.Errorf("object does not implement the Object interfaces")

func Accessor(obj interface{}) (metav1.Object, error) {
	switch t := obj.(type) {
	case metav1.Object:
		return t, nil
	case metav1.ObjectMetaAccessor:
		if m := t.GetObjectMeta(); m != nil {
			return m, nil
		}
		return nil, errNotObject
	default:
		return nil, errNotObject
	}
}
