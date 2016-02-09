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

//Safe array conversion
func IntSliceToCInt(data []int) []C.int {
	sliceSz := len(data)
	result := make([]C.int, sliceSz)
	for i := 0; i < sliceSz; i++ {
		result[i] = C.int(data[i])
	}
	return result
}

//Safe array conversion
func CIntSliceToInt(data []C.GUIntBig) []uint64 {
	sliceSz := len(data)
	result := make([]uint64, sliceSz)
	for i := 0; i < sliceSz; i++ {
		result[i] = uint64(data[i])
	}
	return result
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

// Types of color interpretation for raster bands.
type ColorInterp int

const (
	CI_Undefined      = ColorInterp(C.GCI_Undefined)
	CI_GrayIndex      = ColorInterp(C.GCI_GrayIndex)
	CI_PaletteIndex   = ColorInterp(C.GCI_PaletteIndex)
	CI_RedBand        = ColorInterp(C.GCI_RedBand)
	CI_GreenBand      = ColorInterp(C.GCI_GreenBand)
	CI_BlueBand       = ColorInterp(C.GCI_BlueBand)
	CI_AlphaBand      = ColorInterp(C.GCI_AlphaBand)
	CI_HueBand        = ColorInterp(C.GCI_HueBand)
	CI_SaturationBand = ColorInterp(C.GCI_SaturationBand)
	CI_LightnessBand  = ColorInterp(C.GCI_LightnessBand)
	CI_CyanBand       = ColorInterp(C.GCI_CyanBand)
	CI_MagentaBand    = ColorInterp(C.GCI_MagentaBand)
	CI_YellowBand     = ColorInterp(C.GCI_YellowBand)
	CI_BlackBand      = ColorInterp(C.GCI_BlackBand)
	CI_YCbCr_YBand    = ColorInterp(C.GCI_YCbCr_YBand)
	CI_YCbCr_CbBand   = ColorInterp(C.GCI_YCbCr_CbBand)
	CI_YCbCr_CrBand   = ColorInterp(C.GCI_YCbCr_CrBand)
	CI_Max            = ColorInterp(C.GCI_Max)
)

func (colorInterp ColorInterp) Name() string {
	return C.GoString(C.GDALGetColorInterpretationName(C.GDALColorInterp(colorInterp)))
}

func GetColorInterpretationByName(name string) ColorInterp {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return ColorInterp(C.GDALGetColorInterpretationByName(cName))
}

// Types of color interpretations for a GDALColorTable.
type PaletteInterp int

const (
	// Grayscale (in GDALColorEntry.c1)
	PI_Gray = PaletteInterp(C.GPI_Gray)
	// Red, Green, Blue and Alpha in (in c1, c2, c3 and c4)
	PI_RGB = PaletteInterp(C.GPI_RGB)
	// Cyan, Magenta, Yellow and Black (in c1, c2, c3 and c4)
	PI_CMYK = PaletteInterp(C.GPI_CMYK)
	// Hue, Lightness and Saturation (in c1, c2, and c3)
	PI_HLS = PaletteInterp(C.GPI_HLS)
)

func (paletteInterp PaletteInterp) Name() string {
	return C.GoString(C.GDALGetPaletteInterpretationName(C.GDALPaletteInterp(paletteInterp)))
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

type ColorTable struct {
	cval C.GDALColorTableH
}

type RasterAttributeTable struct {
	cval C.GDALRasterAttributeTableH
}

type AsyncReader struct {
	cval C.GDALAsyncReaderH
}

type ColorEntry struct {
	cval *C.GDALColorEntry
}

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

// -----------------------------------------------------------------------

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

/* ==================================================================== */
/*      GDAL_GCP                                                        */
/* ==================================================================== */

// Unimplemented: InitGCPs
// Unimplemented: DeinitGCPs
// Unimplemented: DuplicateGCPs
// Unimplemented: GCPsToGeoTransform
// Unimplemented: ApplyGeoTransform

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
/*      Color tables.                                                   */
/* ==================================================================== */

// Construct a new color table
func CreateColorTable(interp PaletteInterp) ColorTable {
	ct := C.GDALCreateColorTable(C.GDALPaletteInterp(interp))
	return ColorTable{ct}
}

// Destroy the color table
func (ct ColorTable) Destroy() {
	C.GDALDestroyColorTable(ct.cval)
}

// Make a copy of the color table
func (ct ColorTable) Clone() ColorTable {
	newCT := C.GDALCloneColorTable(ct.cval)
	return ColorTable{newCT}
}

// Fetch palette interpretation
func (ct ColorTable) PaletteInterpretation() PaletteInterp {
	pi := C.GDALGetPaletteInterpretation(ct.cval)
	return PaletteInterp(pi)
}

// Get number of color entries in table
func (ct ColorTable) EntryCount() int {
	count := C.GDALGetColorEntryCount(ct.cval)
	return int(count)
}

// Fetch a color entry from table
func (ct ColorTable) Entry(index int) ColorEntry {
	entry := C.GDALGetColorEntry(ct.cval, C.int(index))
	return ColorEntry{entry}
}

// Unimplemented: EntryAsRGB

// Set entry in color table
func (ct ColorTable) SetEntry(index int, entry ColorEntry) {
	C.GDALSetColorEntry(ct.cval, C.int(index), entry.cval)
}

// Create color ramp
func (ct ColorTable) CreateColorRamp(start, end int, startColor, endColor ColorEntry) {
	C.GDALCreateColorRamp(ct.cval, C.int(start), startColor.cval, C.int(end), endColor.cval)
}

/* ==================================================================== */
/*      Raster Attribute Table                                          */
/* ==================================================================== */

type RATFieldType int

const (
	GFT_Integer = RATFieldType(C.GFT_Integer)
	GFT_Real    = RATFieldType(C.GFT_Real)
	GFT_String  = RATFieldType(C.GFT_String)
)

type RATFieldUsage int

const (
	GFU_Generic    = RATFieldUsage(C.GFU_Generic)
	GFU_PixelCount = RATFieldUsage(C.GFU_PixelCount)
	GFU_Name       = RATFieldUsage(C.GFU_Name)
	GFU_Min        = RATFieldUsage(C.GFU_Min)
	GFU_Max        = RATFieldUsage(C.GFU_Max)
	GFU_MinMax     = RATFieldUsage(C.GFU_MinMax)
	GFU_Red        = RATFieldUsage(C.GFU_Red)
	GFU_Green      = RATFieldUsage(C.GFU_Green)
	GFU_Blue       = RATFieldUsage(C.GFU_Blue)
	GFU_Alpha      = RATFieldUsage(C.GFU_Alpha)
	GFU_RedMin     = RATFieldUsage(C.GFU_RedMin)
	GFU_GreenMin   = RATFieldUsage(C.GFU_GreenMin)
	GFU_BlueMin    = RATFieldUsage(C.GFU_BlueMin)
	GFU_AlphaMin   = RATFieldUsage(C.GFU_AlphaMin)
	GFU_RedMax     = RATFieldUsage(C.GFU_RedMax)
	GFU_GreenMax   = RATFieldUsage(C.GFU_GreenMax)
	GFU_BlueMax    = RATFieldUsage(C.GFU_BlueMax)
	GFU_AlphaMax   = RATFieldUsage(C.GFU_AlphaMax)
	GFU_MaxCount   = RATFieldUsage(C.GFU_MaxCount)
)

// Construct empty raster attribute table
func CreateRasterAttributeTable() RasterAttributeTable {
	rat := C.GDALCreateRasterAttributeTable()
	return RasterAttributeTable{rat}
}

// Destroy a RAT
func (rat RasterAttributeTable) Destroy() {
	C.GDALDestroyRasterAttributeTable(rat.cval)
}

// Fetch table column count
func (rat RasterAttributeTable) ColumnCount() int {
	count := C.GDALRATGetColumnCount(rat.cval)
	return int(count)
}

// Fetch the name of indicated column
func (rat RasterAttributeTable) NameOfCol(index int) string {
	name := C.GDALRATGetNameOfCol(rat.cval, C.int(index))
	return C.GoString(name)
}

// Fetch the usage of indicated column
func (rat RasterAttributeTable) UsageOfCol(index int) RATFieldUsage {
	rfu := C.GDALRATGetUsageOfCol(rat.cval, C.int(index))
	return RATFieldUsage(rfu)
}

// Fetch the type of indicated column
func (rat RasterAttributeTable) TypeOfCol(index int) RATFieldType {
	rft := C.GDALRATGetTypeOfCol(rat.cval, C.int(index))
	return RATFieldType(rft)
}

// Fetch column index for indicated usage
func (rat RasterAttributeTable) ColOfUsage(rfu RATFieldUsage) int {
	index := C.GDALRATGetColOfUsage(rat.cval, C.GDALRATFieldUsage(rfu))
	return int(index)
}

// Fetch row count
func (rat RasterAttributeTable) RowCount() int {
	count := C.GDALRATGetRowCount(rat.cval)
	return int(count)
}

// Fetch field value as string
func (rat RasterAttributeTable) ValueAsString(row, field int) string {
	cString := C.GDALRATGetValueAsString(rat.cval, C.int(row), C.int(field))
	return C.GoString(cString)
}

// Fetch field value as integer
func (rat RasterAttributeTable) ValueAsInt(row, field int) int {
	val := C.GDALRATGetValueAsInt(rat.cval, C.int(row), C.int(field))
	return int(val)
}

// Fetch field value as float64
func (rat RasterAttributeTable) ValueAsFloat64(row, field int) float64 {
	val := C.GDALRATGetValueAsDouble(rat.cval, C.int(row), C.int(field))
	return float64(val)
}

// Set field value from string
func (rat RasterAttributeTable) SetValueAsString(row, field int, val string) {
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	C.GDALRATSetValueAsString(rat.cval, C.int(row), C.int(field), cVal)
}

// Set field value from integer
func (rat RasterAttributeTable) SetValueAsInt(row, field, val int) {
	C.GDALRATSetValueAsInt(rat.cval, C.int(row), C.int(field), C.int(val))
}

// Set field value from float64
func (rat RasterAttributeTable) SetValueAsFloat64(row, field int, val float64) {
	C.GDALRATSetValueAsDouble(rat.cval, C.int(row), C.int(field), C.double(val))
}

// Set row count
func (rat RasterAttributeTable) SetRowCount(count int) {
	C.GDALRATSetRowCount(rat.cval, C.int(count))
}

// Create new column
func (rat RasterAttributeTable) CreateColumn(name string, rft RATFieldType, rfu RATFieldUsage) error {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	return C.GDALRATCreateColumn(rat.cval, cName, C.GDALRATFieldType(rft), C.GDALRATFieldUsage(rfu)).Err()
}

// Set linear binning information
func (rat RasterAttributeTable) SetLinearBinning(row0min, binsize float64) error {
	return C.GDALRATSetLinearBinning(rat.cval, C.double(row0min), C.double(binsize)).Err()
}

// Fetch linear binning information
func (rat RasterAttributeTable) LinearBinning() (row0min, binsize float64, exists bool) {
	success := C.GDALRATGetLinearBinning(rat.cval, (*C.double)(&row0min), (*C.double)(&binsize))
	return row0min, binsize, success != 0
}

// Initialize RAT from color table
func (rat RasterAttributeTable) FromColorTable(ct ColorTable) error {
	return C.GDALRATInitializeFromColorTable(rat.cval, ct.cval).Err()
}

// Translate RAT to a color table
func (rat RasterAttributeTable) ToColorTable(count int) ColorTable {
	ct := C.GDALRATTranslateToColorTable(rat.cval, C.int(count))
	return ColorTable{ct}
}

// Dump RAT in readable form to a file
// Unimplemented: DumpReadable

// Get row for pixel value
func (rat RasterAttributeTable) RowOfValue(val float64) (int, bool) {
	row := C.GDALRATGetRowOfValue(rat.cval, C.double(val))
	return int(row), row != -1
}

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
