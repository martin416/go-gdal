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
	"errors"
	"unsafe"
)

func init() {
	C.GDALAllRegister()
}

/* -------------------------------------------------------------------- */
/*      Significant constants.                                          */
/* -------------------------------------------------------------------- */

const (
	VERSION_MAJOR = int(C.GDAL_VERSION_MAJOR)
	VERSION_MINOR = int(C.GDAL_VERSION_MINOR)
	VERSION_REV   = int(C.GDAL_VERSION_REV)
	VERSION_BUILD = int(C.GDAL_VERSION_BUILD)
	VERSION_NUM   = int(C.GDAL_VERSION_NUM)
	RELEASE_DATE  = int(C.GDAL_RELEASE_DATE)
	RELEASE_NAME  = string(C.GDAL_RELEASE_NAME)
)

var (
	ErrDebug   = errors.New("Debug Error")
	ErrWarning = errors.New("Warning Error")
	ErrFailure = errors.New("Failure Error")
	ErrFatal   = errors.New("Fatal Error")
	ErrIllegal = errors.New("Illegal Error")
)

// Error handling.  The following is bare-bones, and needs to be replaced with something more useful.
func (err _Ctype_CPLErr) Err() error {
	switch err {
	case 0:
		return nil
	case 1:
		return ErrDebug
	case 2:
		return ErrWarning
	case 3:
		return ErrFailure
	case 4:
		return ErrFailure
	}
	return ErrIllegal
}

func (err _Ctype_OGRErr) Err() error {
	switch err {
	case 0:
		return nil
	case 1:
		return ErrDebug
	case 2:
		return ErrWarning
	case 3:
		return ErrFailure
	case 4:
		return ErrFailure
	}
	return ErrIllegal
}

// Pixel data types
type DataType int

const (
	Unknown  = DataType(C.GDT_Unknown)
	Byte     = DataType(C.GDT_Byte)
	UInt16   = DataType(C.GDT_UInt16)
	Int16    = DataType(C.GDT_Int16)
	UInt32   = DataType(C.GDT_UInt32)
	Int32    = DataType(C.GDT_Int32)
	Float32  = DataType(C.GDT_Float32)
	Float64  = DataType(C.GDT_Float64)
	CInt16   = DataType(C.GDT_CInt16)
	CInt32   = DataType(C.GDT_CInt32)
	CFloat32 = DataType(C.GDT_CFloat32)
	CFloat64 = DataType(C.GDT_CFloat64)
)

// Get data type size in bits.
func (dataType DataType) Size() int {
	return int(C.GDALGetDataTypeSize(C.GDALDataType(dataType)))
}

func (dataType DataType) IsComplex() int {
	return int(C.GDALDataTypeIsComplex(C.GDALDataType(dataType)))
}

func (dataType DataType) Name() string {
	return C.GoString(C.GDALGetDataTypeName(C.GDALDataType(dataType)))
}

func (dataType DataType) Union(dataTypeB DataType) DataType {
	return DataType(
		C.GDALDataTypeUnion(C.GDALDataType(dataType), C.GDALDataType(dataTypeB)),
	)
}

// status of the asynchronous stream
type AsyncStatusType int

const (
	AR_Pending  = AsyncStatusType(C.GARIO_PENDING)
	AR_Update   = AsyncStatusType(C.GARIO_UPDATE)
	AR_Error    = AsyncStatusType(C.GARIO_ERROR)
	AR_Complete = AsyncStatusType(C.GARIO_COMPLETE)
)

func (statusType AsyncStatusType) Name() string {
	return C.GoString(C.GDALGetAsyncStatusTypeName(C.GDALAsyncStatusType(statusType)))
}

func GetAsyncStatusTypeByName(statusTypeName string) AsyncStatusType {
	name := C.CString(statusTypeName)
	defer C.free(unsafe.Pointer(name))
	return AsyncStatusType(C.GDALGetAsyncStatusTypeByName(name))
}

// Flag indicating read/write, or read-only access to data.
type Access int

const (
	// Read only (no update) access
	ReadOnly = Access(C.GA_ReadOnly)
	// Read/write access.
	Update = Access(C.GA_Update)
)

// Read/Write flag for RasterIO() method
type RWFlag int

const (
	// Read data
	Read = RWFlag(C.GF_Read)
	// Write data
	Write = RWFlag(C.GF_Write)
)

//Safe array conversion
func intSliceToCInt(data []int) []C.int {
	sliceSz := len(data)
	result := make([]C.int, sliceSz)
	for i := 0; i < sliceSz; i++ {
		result[i] = C.int(data[i])
	}
	return result
}

//Safe array conversion
func cIntSliceToInt(data []C.GUIntBig) []uint64 {
	sliceSz := len(data)
	result := make([]uint64, sliceSz)
	for i := 0; i < sliceSz; i++ {
		result[i] = uint64(data[i])
	}
	return result
}

// "well known" metadata items.
const (
	MD_AREA_OR_POINT = string(C.GDALMD_AREA_OR_POINT)
	MD_AOP_AREA      = string(C.GDALMD_AOP_AREA)
	MD_AOP_POINT     = string(C.GDALMD_AOP_POINT)
)

/* -------------------------------------------------------------------- */
/*      Define handle types related to various internal classes.        */
/* -------------------------------------------------------------------- */

type MajorObject struct {
	cval C.GDALMajorObjectH
}

type AsyncReader struct {
	cval C.GDALAsyncReaderH
}

/* ==================================================================== */
/*      major objects (dataset, and, driver, drivermanager).            */
/* ==================================================================== */

// Fetch object description
func (object MajorObject) Description() string {
	cObject := object.cval
	desc := C.GoString(C.GDALGetDescription(cObject))
	return desc
}

// Set object description
func (object MajorObject) SetDescription(desc string) {
	cObject := object.cval
	cDesc := C.CString(desc)
	defer C.free(unsafe.Pointer(cDesc))
	C.GDALSetDescription(cObject, cDesc)
}

// Fetch metadata
func (object MajorObject) Metadata(domain string) []string {
	panic("not implemented!")
	return nil
}

// Set metadata
func (object MajorObject) SetMetadata(metadata []string, domain string) {
	panic("not implemented!")
	return
}

// Fetch a single metadata item
func (object MajorObject) MetadataItem(name, domain string) string {
	panic("not implemented!")
	return ""
}

// Set a single metadata item
func (object MajorObject) SetMetadataItem(name, value, domain string) {
	panic("not implemented!")
	return
}

// TODO: Make korrekt class hirerarchy via interfaces

func (object *RasterBand) SetMetadataItem(name, value, domain string) error {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_value))

	c_domain := C.CString(domain)
	defer C.free(unsafe.Pointer(c_domain))

	return C.GDALSetMetadataItem(
		C.GDALMajorObjectH(unsafe.Pointer(object.cval)),
		c_name, c_value, c_domain,
	).Err()
}

// TODO: Make korrekt class hirerarchy via interfaces

func (object *Dataset) SetMetadataItem(name, value, domain string) error {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_value := C.CString(value)
	defer C.free(unsafe.Pointer(c_value))

	c_domain := C.CString(domain)
	defer C.free(unsafe.Pointer(c_domain))

	return C.GDALSetMetadataItem(
		C.GDALMajorObjectH(unsafe.Pointer(object.cval)),
		c_name, c_value, c_domain,
	).Err()
}

// Fetch single metadata item.
func (object *Driver) MetadataItem(name, domain string) string {
	c_name := C.CString(name)
	defer C.free(unsafe.Pointer(c_name))

	c_domain := C.CString(domain)
	defer C.free(unsafe.Pointer(c_domain))

	return C.GoString(
		C.GDALGetMetadataItem(
			C.GDALMajorObjectH(unsafe.Pointer(object.cval)),
			c_name, c_domain,
		),
	)
}

// Generate downsampled overviews
// Unimplemented: RegenerateOverviews

/* ==================================================================== */
/*     GDALAsyncReader                                                  */
/* ==================================================================== */

// Unimplemented: GetNextUpdatedRegion
// Unimplemented: LockBuffer
// Unimplemented: UnlockBuffer

/* ==================================================================== */
/*      GDAL Cache Management                                           */
/* ==================================================================== */

// Set maximum cache memory
func SetCacheMax(bytes int) {
	C.GDALSetCacheMax64(C.GIntBig(bytes))
}

// Get maximum cache memory
func GetCacheMax() int {
	bytes := C.GDALGetCacheMax64()
	return int(bytes)
}

// Get cache memory used
func GetCacheUsed() int {
	bytes := C.GDALGetCacheUsed64()
	return int(bytes)
}

// Try to flush one cached raster block
func FlushCacheBlock() bool {
	flushed := C.GDALFlushCacheBlock()
	return flushed != 0
}
