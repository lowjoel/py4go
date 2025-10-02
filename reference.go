package python

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import "runtime"

//
// Reference
//

type Reference struct {
	Object *C.PyObject
}

func py_DecRef(object *C.PyObject) {
	state := EnsureGilState()
	defer state.Release()

	C.Py_DecRef(object)
}

func NewWeakReference(pyObject *C.PyObject) *Reference {
	return &Reference{pyObject}
}

func NewReference(pyObject *C.PyObject) *Reference {
	r := NewWeakReference(pyObject)
	C.Py_IncRef(pyObject)
	runtime.AddCleanup(r, py_DecRef, r.Object)

	return r
}

func (self *Reference) Type() *Type {
	return NewType(self.Object.ob_type)
}
