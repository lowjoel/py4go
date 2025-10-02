package python

// See:
//   https://docs.python.org/3/c-api/exceptions.html

import (
	"errors"
	"fmt"
)

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func HasException() bool {
	return C.PyErr_Occurred() != nil
}

func GetError() error {
	if exception := FetchException(); exception != nil {
		return exception
	} else {
		return errors.New("Python error without an exception")
	}
}

//
// Exception
//

type Exception struct {
	typeStr  string
	valueStr string
}

func FetchException() *Exception {
	var type_, value, traceback *C.PyObject
	C.PyErr_Fetch(&type_, &value, &traceback)
	if type_ != nil {
		defer C.PyErr_Restore(type_, value, traceback)

		var type__, value_, traceback_ *Reference

		if type_ != nil {
			type__ = NewWeakReference(type_)
		}

		if value != nil {
			value_ = NewWeakReference(value)
		}

		if traceback != nil {
			traceback_ = NewWeakReference(traceback)
		}

		return NewExceptionRaw(type__, value_, traceback_)
	} else {
		return nil
	}
}

func NewExceptionRaw(type_ *Reference, value *Reference, traceback *Reference) *Exception {
	gil := EnsureGilState()
	defer gil.Release()

	// Store the string versions of these objects because we do not own the exception value.
	res := &Exception{
		typeStr:  type_.String(),
		valueStr: value.String(),
	}
	return res
}

// error signature
func (self *Exception) Error() string {
	return fmt.Sprintf("%s: %s", self.typeStr, self.valueStr)
}
