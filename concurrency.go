package python

// See:
//   https://docs.python.org/3/c-api/init.html

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"
import "runtime"

//
// ThreadState
//

type ThreadState struct {
	State *C.PyThreadState
}

func SaveThreadState() *ThreadState {
	return &ThreadState{C.PyEval_SaveThread()}
}

func (self *ThreadState) Restore() {
	C.PyEval_RestoreThread(self.State)
}

//
// GilState
//

type GilState struct {
	State C.PyGILState_STATE
}

func EnsureGilState() *GilState {
	runtime.LockOSThread()
	return &GilState{C.PyGILState_Ensure()}
}

func (self *GilState) Release() {
	C.PyGILState_Release(self.State)
	runtime.UnlockOSThread()
}
