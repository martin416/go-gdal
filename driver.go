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
	"fmt"
	"unsafe"
)

const (
	DMD_LONGNAME           = string(C.GDAL_DMD_LONGNAME)
	DMD_HELPTOPIC          = string(C.GDAL_DMD_HELPTOPIC)
	DMD_MIMETYPE           = string(C.GDAL_DMD_MIMETYPE)
	DMD_EXTENSION          = string(C.GDAL_DMD_EXTENSION)
	DMD_CREATIONOPTIONLIST = string(C.GDAL_DMD_CREATIONOPTIONLIST)
	DMD_CREATIONDATATYPES  = string(C.GDAL_DMD_CREATIONDATATYPES)

	DCAP_CREATE     = string(C.GDAL_DCAP_CREATE)
	DCAP_CREATECOPY = string(C.GDAL_DCAP_CREATECOPY)
	DCAP_VIRTUALIO  = string(C.GDAL_DCAP_VIRTUALIO)
)

type Driver struct {
	cval C.GDALDriverH
}

// Create a new dataset with this driver.
func (driver Driver) Create(
	filename string,
	xSize, ySize, bands int,
	dataType DataType,
	options []string,
) Dataset {
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	h := C.GDALCreate(
		driver.cval,
		name,
		C.int(xSize), C.int(ySize), C.int(bands),
		C.GDALDataType(dataType),
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	return Dataset{h}
}

// Create a copy of a dataset
func (driver Driver) CreateCopy(
	filename string,
	sourceDataset Dataset,
	strict int,
	options []string,
	progress ProgressFunc,
	data interface{},
) Dataset {
	name := C.CString(filename)
	defer C.free(unsafe.Pointer(name))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	var h C.GDALDatasetH

	if progress == nil {
		h = C.GDALCreateCopy(
			driver.cval, name,
			sourceDataset.cval,
			C.int(strict),
			(**C.char)(unsafe.Pointer(&opts[0])),
			nil,
			nil,
		)
	} else {
		arg := &goGDALProgressFuncProxyArgs{
			progress, data,
		}
		h = C.GDALCreateCopy(
			driver.cval, name,
			sourceDataset.cval,
			C.int(strict), (**C.char)(unsafe.Pointer(&opts[0])),
			C.goGDALProgressFuncProxyB(),
			unsafe.Pointer(arg),
		)
	}

	return Dataset{h}
}

// Return the driver needed to access the provided dataset name.
func IdentifyDriver(filename string, filenameList []string) Driver {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	length := len(filenameList)
	cFilenameList := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cFilenameList[i] = C.CString(filenameList[i])
		defer C.free(unsafe.Pointer(cFilenameList[i]))
	}
	cFilenameList[length] = (*C.char)(unsafe.Pointer(nil))

	driver := C.GDALIdentifyDriver(cFilename, (**C.char)(unsafe.Pointer(&cFilenameList[0])))
	return Driver{driver}
}

// Open an existing dataset
func Open(filename string, access Access) (Dataset, error) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	dataset := C.GDALOpen(cFilename, C.GDALAccess(access))
	if dataset == nil {
		return Dataset{nil}, fmt.Errorf("Error: dataset '%s' open error", filename)
	}
	return Dataset{dataset}, nil
}

// Open a shared existing dataset
func OpenShared(filename string, access Access) Dataset {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	dataset := C.GDALOpenShared(cFilename, C.GDALAccess(access))
	return Dataset{dataset}
}

// Unimplemented: DumpOpenDatasets

// Return the driver by short name
func GetDriverByName(driverName string) (Driver, error) {
	cName := C.CString(driverName)
	defer C.free(unsafe.Pointer(cName))

	driver := C.GDALGetDriverByName(cName)
	if driver == nil {
		return Driver{driver}, fmt.Errorf("Error: driver '%s' not found", driverName)
	}
	return Driver{driver}, nil
}

// Fetch the number of registered drivers.
func GetDriverCount() int {
	nDrivers := C.GDALGetDriverCount()
	return int(nDrivers)
}

// Fetch driver by index
func GetDriver(index int) Driver {
	driver := C.GDALGetDriver(C.int(index))
	return Driver{driver}
}

// Destroy a GDAL driver
func (driver Driver) Destroy() {
	C.GDALDestroyDriver(driver.cval)
}

// Registers a driver for use
func (driver Driver) Register() int {
	index := C.GDALRegisterDriver(driver.cval)
	return int(index)
}

// Reregister the driver
func (driver Driver) Deregister() {
	C.GDALDeregisterDriver(driver.cval)
}

// Destroy the driver manager
func DestroyDriverManager() {
	C.GDALDestroyDriverManager()
}

// Delete named dataset
func (driver Driver) DeleteDataset(name string) error {
	cDriver := driver.cval
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.GDALDeleteDataset(cDriver, cName).Err()
}

// Rename named dataset
func (driver Driver) RenameDataset(newName, oldName string) error {
	cDriver := driver.cval
	cNewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cNewName))
	cOldName := C.CString(oldName)
	defer C.free(unsafe.Pointer(cOldName))
	return C.GDALRenameDataset(cDriver, cNewName, cOldName).Err()
}

// Copy all files associated with the named dataset
func (driver Driver) CopyDatasetFiles(newName, oldName string) error {
	cDriver := driver.cval
	cNewName := C.CString(newName)
	defer C.free(unsafe.Pointer(cNewName))
	cOldName := C.CString(oldName)
	defer C.free(unsafe.Pointer(cOldName))
	return C.GDALCopyDatasetFiles(cDriver, cNewName, cOldName).Err()
}

// Get the short name associated with this driver
func (driver Driver) ShortName() string {
	cDriver := driver.cval
	return C.GoString(C.GDALGetDriverShortName(cDriver))
}

// Get the long name associated with this driver
func (driver Driver) LongName() string {
	cDriver := driver.cval
	return C.GoString(C.GDALGetDriverLongName(cDriver))
}
