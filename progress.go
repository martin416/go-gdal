package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"

#cgo linux  pkg-config: gdal
#cgo darwin pkg-config: gdal
#cgo windows LDFLAGS: -Lc:/gdal/release-1600-x64/lib -lgdal_i
#cgo windows CFLAGS: -IC:/gdal/release-1600-x64/include
*/
import "C"
import (
	"unsafe"
)

/* -------------------------------------------------------------------- */
/*      Callback "progress" function.                                   */
/* -------------------------------------------------------------------- */

type ProgressFunc func(complete float64, message string, progressArg interface{}) int

func DummyProgress(complete float64, message string, data interface{}) int {
	msg := C.CString(message)
	defer C.free(unsafe.Pointer(msg))

	retval := C.GDALDummyProgress(C.double(complete), msg, unsafe.Pointer(nil))
	return int(retval)
}

func TermProgress(complete float64, message string, data interface{}) int {
	msg := C.CString(message)
	defer C.free(unsafe.Pointer(msg))

	retval := C.GDALTermProgress(C.double(complete), msg, unsafe.Pointer(nil))
	return int(retval)
}

func ScaledProgress(complete float64, message string, data interface{}) int {
	msg := C.CString(message)
	defer C.free(unsafe.Pointer(msg))

	retval := C.GDALScaledProgress(C.double(complete), msg, unsafe.Pointer(nil))
	return int(retval)
}

func CreateScaledProgress(min, max float64, progress ProgressFunc, data unsafe.Pointer) unsafe.Pointer {
	panic("not implemented!")
	return nil
}

func DestroyScaledProgress(data unsafe.Pointer) {
	C.GDALDestroyScaledProgress(data)
}

type goGDALProgressFuncProxyArgs struct {
	progresssFunc ProgressFunc
	data          interface{}
}

//export goGDALProgressFuncProxyA
func goGDALProgressFuncProxyA(complete C.double, message *C.char, data unsafe.Pointer) int {
	arg := (*goGDALProgressFuncProxyArgs)(data)
	return arg.progresssFunc(
		float64(complete), C.GoString(message), arg.data,
	)
}
