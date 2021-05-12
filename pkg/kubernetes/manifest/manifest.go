package manifest

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metaunstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Annotated(o Object, annotations map[string]string) {
	merged := map[string]string{}

	for k, v := range o.GetAnnotations() {
		merged[k] = v
	}

	for k, v := range annotations {
		merged[k] = v
	}

	o.SetAnnotations(merged)
}

type Object = client.Object
type ObjectList = client.ObjectList

func ObjectFromRuntimeObject(ro runtime.Object) (Object, error) {
	o, err := meta.Accessor(ro)
	if err != nil {
		return nil, err
	}
	return o.(Object), nil
}

func ObjectListFromRuntimeObject(ro runtime.Object) (ObjectList, error) {
	o, err := meta.ListAccessor(ro)
	if err != nil {
		return nil, err
	}
	return o.(ObjectList), nil
}

func DeepCopy(m Object) (Object, error) {
	return ObjectFromRuntimeObject(m.DeepCopyObject())
}

func Identity(o Object) string {
	gvk := o.GetObjectKind().GroupVersionKind()

	if namespace := o.GetNamespace(); namespace == "" {
		return fmt.Sprintf("%s/%s/%s", gvk.GroupVersion(), gvk.Kind, o.GetName())
	}

	return fmt.Sprintf("%s/%s/%s/%s", gvk.GroupVersion(), gvk.Kind, o.GetNamespace(), o.GetName())
}

type List []Object

func (l List) DeepCopy() List {
	li := make(List, len(l))
	for i := range li {
		li[i], _ = DeepCopy(l[i])
	}
	return li
}
func (l List) Orphaned(expectedList List) List {
	list := List{}
	remains := map[string]bool{}

	for i := range expectedList {
		o := expectedList[i]
		remains[Identity(o)] = true
	}

	for i := range l {
		o := l[i]

		if _, ok := remains[Identity(o)]; ok {
			list = append(list, o)
		}
	}

	return list
}

func KindName(m Object) string {
	return m.GetObjectKind().GroupVersionKind().Kind + "/" + m.GetName()
}

func NewForGroupVersionKind(gvk schema.GroupVersionKind) (Object, error) {
	ro := &metaunstructured.Unstructured{}
	ro.GetObjectKind().SetGroupVersionKind(gvk)
	return ObjectFromRuntimeObject(ro)
}

func NewListForGroupVersionKind(gvk schema.GroupVersionKind) (ObjectList, error) {
	ro := &metaunstructured.UnstructuredList{}
	ro.GetObjectKind().SetGroupVersionKind(gvk)
	return ObjectListFromRuntimeObject(ro)
}
