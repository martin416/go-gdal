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


type FieldType int

const (
	FT_Integer     = FieldType(C.OFTInteger)
	FT_IntegerList = FieldType(C.OFTIntegerList)
	FT_Real        = FieldType(C.OFTReal)
	FT_RealList    = FieldType(C.OFTRealList)
	FT_String      = FieldType(C.OFTString)
	FT_StringList  = FieldType(C.OFTStringList)
	FT_Binary      = FieldType(C.OFTBinary)
	FT_Date        = FieldType(C.OFTDate)
	FT_Time        = FieldType(C.OFTTime)
	FT_DateTime    = FieldType(C.OFTDateTime)
)

type Justification int

const (
	J_Undefined = Justification(C.OJUndefined)
	J_Left      = Justification(C.OJLeft)
	J_Right     = Justification(C.OJRight)
)

type FieldDefinition struct {
	cval C.OGRFieldDefnH
}

type Field struct {
	cval *C.OGRField
}

// Create a new field definition
func CreateFieldDefinition(name string, fieldType FieldType) FieldDefinition {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	fieldDef := C.OGR_Fld_Create(cName, C.OGRFieldType(fieldType))
	return FieldDefinition{fieldDef}
}

// Destroy the field definition
func (fd FieldDefinition) Destroy() {
	C.OGR_Fld_Destroy(fd.cval)
}

// Fetch the name of the field
func (fd FieldDefinition) Name() string {
	name := C.OGR_Fld_GetNameRef(fd.cval)
	return C.GoString(name)
}

// Set the name of the field
func (fd FieldDefinition) SetName(name string) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	C.OGR_Fld_SetName(fd.cval, cName)
}

// Fetch the type of this field
func (fd FieldDefinition) Type() FieldType {
	fType := C.OGR_Fld_GetType(fd.cval)
	return FieldType(fType)
}

// Set the type of this field
func (fd FieldDefinition) SetType(fType FieldType) {
	C.OGR_Fld_SetType(fd.cval, C.OGRFieldType(fType))
}

// Fetch the justification for this field
func (fd FieldDefinition) Justification() Justification {
	justify := C.OGR_Fld_GetJustify(fd.cval)
	return Justification(justify)
}

// Set the justification for this field
func (fd FieldDefinition) SetJustification(justify Justification) {
	C.OGR_Fld_SetJustify(fd.cval, C.OGRJustification(justify))
}

// Fetch the formatting width for this field
func (fd FieldDefinition) Width() int {
	width := C.OGR_Fld_GetWidth(fd.cval)
	return int(width)
}

// Set the formatting width for this field
func (fd FieldDefinition) SetWidth(width int) {
	C.OGR_Fld_SetWidth(fd.cval, C.int(width))
}

// Fetch the precision for this field
func (fd FieldDefinition) Precision() int {
	precision := C.OGR_Fld_GetPrecision(fd.cval)
	return int(precision)
}

// Set the precision for this field
func (fd FieldDefinition) SetPrecision(precision int) {
	C.OGR_Fld_SetPrecision(fd.cval, C.int(precision))
}

// Set defining parameters of field in a single call
func (fd FieldDefinition) Set(
	name string,
	fType FieldType,
	width, precision int,
	justify Justification,
) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	C.OGR_Fld_Set(
		fd.cval,
		cName,
		C.OGRFieldType(fType),
		C.int(width),
		C.int(precision),
		C.OGRJustification(justify),
	)
}

// Fetch whether this field should be ignored when fetching features
func (fd FieldDefinition) IsIgnored() bool {
	ignore := C.OGR_Fld_IsIgnored(fd.cval)
	return ignore != 0
}

// Set whether this field should be ignored when fetching features
func (fd FieldDefinition) SetIgnored(ignore bool) {
	C.OGR_Fld_SetIgnored(fd.cval, BoolToCInt(ignore))
}

// Fetch human readable name for the field type
func (ft FieldType) Name() string {
	name := C.OGR_GetFieldTypeName(C.OGRFieldType(ft))
	return C.GoString(name)
}


