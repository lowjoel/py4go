package python

// See:
//   https://docs.python.org/3/c-api/init.html

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import "errors"

func Initialize() {
	C.Py_Initialize()
}

func Finalize() error {
	var repanic *Exception
	if panicked := recover(); panicked != nil {
		if err, ok := panicked.(error); ok {
			errors.As(err, &repanic)
		}
	}

	finalizeErr := C.Py_FinalizeEx() != 0
	switch {
	case repanic != nil:
		panic(repanic)
	case finalizeErr:
		return GetError()
	default:
		return nil
	}
}

func Version() string {
	return C.GoString(C.Py_GetVersion())
}
