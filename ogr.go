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

func init() {
	C.OGRRegisterAll()
}

// Clean up all OGR related resources
func CleanupOGR() {
	C.OGRCleanupAll()
}
/* -------------------------------------------------------------------- */
/*      Layer functions                                                 */
/* -------------------------------------------------------------------- */

type Layer struct {
	cval C.OGRLayerH
}

// Return the layer name
func (layer Layer) Name() string {
	name := C.OGR_L_GetName(layer.cval)
	return C.GoString(name)
}

// Return the layer geometry type
func (layer Layer) Type() GeometryType {
	gt := C.OGR_L_GetGeomType(layer.cval)
	return GeometryType(gt)
}

// Return the current spatial filter for this layer
func (layer Layer) SpatialFilter() Geometry {
	geom := C.OGR_L_GetSpatialFilter(layer.cval)
	return Geometry{geom}
}

// Set a new spatial filter for this layer
func (layer Layer) SetSpatialFilter(filter Geometry) {
	C.OGR_L_SetSpatialFilter(layer.cval, filter.cval)
}

// Set a new rectangular spatial filter for this layer
func (layer Layer) SetSpatialFilterRect(minX, minY, maxX, maxY float64) {
	C.OGR_L_SetSpatialFilterRect(
		layer.cval,
		C.double(minX), C.double(minY), C.double(maxX), C.double(maxY),
	)
}

// Set a new attribute query filter
func (layer Layer) SetAttributeFilter(filter string) error {
	cFilter := C.CString(filter)
	defer C.free(unsafe.Pointer(cFilter))
	return C.OGR_L_SetAttributeFilter(layer.cval, cFilter).Err()
}

// Reset reading to start on the first featre
func (layer Layer) ResetReading() {
	C.OGR_L_ResetReading(layer.cval)
}

// Fetch the next available feature from this layer
func (layer Layer) NextFeature() Feature {
	feature := C.OGR_L_GetNextFeature(layer.cval)
	return Feature{feature}
}

// Move read cursor to the provided index
func (layer Layer) SetNextByIndex(index int) error {
	return C.OGR_L_SetNextByIndex(layer.cval, C.GIntBig(index)).Err()
}

// Fetch a feature by its index
func (layer Layer) Feature(index int) Feature {
	feature := C.OGR_L_GetFeature(layer.cval, C.GIntBig(index))
	return Feature{feature}
}

// Rewrite the provided feature
func (layer Layer) SetFeature(feature Feature) error {
	return C.OGR_L_SetFeature(layer.cval, feature.cval).Err()
}

// Create and write a new feature within a layer
func (layer Layer) Create(feature Feature) error {
	return C.OGR_L_CreateFeature(layer.cval, feature.cval).Err()
}

// Delete indicated feature from layer
func (layer Layer) Delete(index int) error {
	return C.OGR_L_DeleteFeature(layer.cval, C.GIntBig(index)).Err()
}

// Fetch the schema information for this layer
func (layer Layer) Definition() FeatureDefinition {
	defn := C.OGR_L_GetLayerDefn(layer.cval)
	return FeatureDefinition{defn}
}

// Fetch the spatial reference system for this layer
func (layer Layer) SpatialReference() SpatialReference {
	sr := C.OGR_L_GetSpatialRef(layer.cval)
	return SpatialReference{sr}
}

// Fetch the feature count for this layer
func (layer Layer) FeatureCount(force bool) (count int, ok bool) {
	count = int(C.OGR_L_GetFeatureCount(layer.cval, BoolToCInt(force)))
	return count, count != -1
}

// Fetch the extent of this layer
func (layer Layer) Extent(force bool) (env Envelope, err error) {
	err = C.OGR_L_GetExtent(layer.cval, &env.cval, BoolToCInt(force)).Err()
	return
}

// Test if this layer supports the named capability
func (layer Layer) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_L_TestCapability(layer.cval, cString)
	return val != 0
}

// Create a new field on a layer
func (layer Layer) CreateField(fd FieldDefinition, approxOK bool) error {
	return C.OGR_L_CreateField(layer.cval, fd.cval, BoolToCInt(approxOK)).Err()
}

// Delete a field from the layer
func (layer Layer) DeleteField(index int) error {
	return C.OGR_L_DeleteField(layer.cval, C.int(index)).Err()
}

// Reorder all the fields of a layer
func (layer Layer) ReorderFields(layerMap []int) error {
	return C.OGR_L_ReorderFields(layer.cval, (*C.int)(unsafe.Pointer(&layerMap[0]))).Err()
}

// Reorder an existing field of a layer
func (layer Layer) ReorderField(oldIndex, newIndex int) error {
	return C.OGR_L_ReorderField(layer.cval, C.int(oldIndex), C.int(newIndex)).Err()
}

// Alter the definition of an existing field of a layer
func (layer Layer) AlterFieldDefn(index int, newDefn FieldDefinition, flags int) error {
	return C.OGR_L_AlterFieldDefn(layer.cval, C.int(index), newDefn.cval, C.int(flags)).Err()
}

// Begin a transation on data sources which support it
func (layer Layer) StartTransaction() error {
	return C.OGR_L_StartTransaction(layer.cval).Err()
}

// Commit a transaction on data sources which support it
func (layer Layer) CommitTransaction() error {
	return C.OGR_L_CommitTransaction(layer.cval).Err()
}

// Roll back the current transaction on data sources which support it
func (layer Layer) RollbackTransaction() error {
	return C.OGR_L_RollbackTransaction(layer.cval).Err()
}

// Flush pending changes to the layer
func (layer Layer) Sync() error {
	return C.OGR_L_SyncToDisk(layer.cval).Err()
}

// Fetch the name of the FID column
func (layer Layer) FIDColumn() string {
	name := C.OGR_L_GetFIDColumn(layer.cval)
	return C.GoString(name)
}

// Fetch the name of the geometry column
func (layer Layer) GeometryColumn() string {
	name := C.OGR_L_GetGeometryColumn(layer.cval)
	return C.GoString(name)
}

// Set which fields can be ignored when retrieving features from the layer
func (layer Layer) SetIgnoredFields(names []string) error {
	length := len(names)
	cNames := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		cNames[i] = C.CString(names[i])
		defer C.free(unsafe.Pointer(cNames[i]))
	}
	cNames[length] = (*C.char)(unsafe.Pointer(nil))

	return C.OGR_L_SetIgnoredFields(layer.cval, (**C.char)(unsafe.Pointer(&cNames[0]))).Err()
}

// Return the intersection of two layers
// Unimplemented: Intersection
// Will be new in 2.0

// Return the union of two layers
// Unimplemented: Union
// Will be new in 2.0

// Return the symmetric difference of two layers
// Unimplemented: SymDifference
// Will be new in 2.0

// Identify features in this layer with ones from the provided layer
// Unimplemented: Identity
// Will be new in 2.0

// Update this layer with features from the provided layer
// Unimplemented: Update
// Will be new in 2.0

// Clip off areas that are not covered by the provided layer
// Unimplemented: Clip
// Will be new in 2.0

// Remove areas that are covered by the provided layer
// Unimplemented: Erase
// Will be new in 2.0

/* -------------------------------------------------------------------- */
/*      Data source functions                                           */
/* -------------------------------------------------------------------- */

type DataSource struct {
	cval C.OGRDataSourceH
}

// Open a file / data source with one of the registered drivers
func OpenDataSource(name string, update int) DataSource {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ds := C.OGROpen(cName, C.int(update), nil)
	return DataSource{ds}
}

// Open a shared file / data source with one of the registered drivers
func OpenSharedDataSource(name string, update int) DataSource {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ds := C.OGROpenShared(cName, C.int(update), nil)
	return DataSource{ds}
}

// Drop a reference to this datasource and destroy if reference is zero
func (ds DataSource) Release() error {
	return C.OGRReleaseDataSource(ds.cval).Err()
}

// Return the number of opened data sources
func OpenDataSourceCount() int {
	count := C.OGRGetOpenDSCount()
	return int(count)
}

// Return the i'th datasource opened
func OpenDataSourceByIndex(index int) DataSource {
	ds := C.OGRGetOpenDS(C.int(index))
	return DataSource{ds}
}

// Closes datasource and releases resources
func (ds DataSource) Destroy() {
	C.OGR_DS_Destroy(ds.cval)
}

// Fetch the name of the data source
func (ds DataSource) Name() string {
	name := C.OGR_DS_GetName(ds.cval)
	return C.GoString(name)
}

// Fetch the number of layers in this data source
func (ds DataSource) LayerCount() int {
	count := C.OGR_DS_GetLayerCount(ds.cval)
	return int(count)
}

// Fetch a layer of this data source by index
func (ds DataSource) LayerByIndex(index int) Layer {
	layer := C.OGR_DS_GetLayer(ds.cval, C.int(index))
	return Layer{layer}
}

// Fetch a layer of this data source by name
func (ds DataSource) LayerByName(name string) Layer {
	cString := C.CString(name)
	defer C.free(unsafe.Pointer(cString))
	layer := C.OGR_DS_GetLayerByName(ds.cval, cString)
	return Layer{layer}
}

// Delete the layer from the data source
func (ds DataSource) Delete(index int) error {
	return C.OGR_DS_DeleteLayer(ds.cval, C.int(index)).Err()
}

// Fetch the driver that the data source was opened with
func (ds DataSource) Driver() OGRDriver {
	driver := C.OGR_DS_GetDriver(ds.cval)
	return OGRDriver{driver}
}

// Create a new layer on the data source
func (ds DataSource) CreateLayer(
	name string,
	sr SpatialReference,
	geomType GeometryType,
	options []string,
) Layer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	layer := C.OGR_DS_CreateLayer(
		ds.cval,
		cName,
		sr.cval,
		C.OGRwkbGeometryType(geomType),
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	return Layer{layer}
}

// Duplicate an existing layer
func (ds DataSource) CopyLayer(
	source Layer,
	name string,
	options []string,
) Layer {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	layer := C.OGR_DS_CopyLayer(
		ds.cval,
		source.cval,
		cName,
		(**C.char)(unsafe.Pointer(&opts[0])),
	)
	return Layer{layer}
}

// Test if the data source has the indicated capability
func (ds DataSource) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_DS_TestCapability(ds.cval, cString)
	return val != 0
}

// Execute an SQL statement against the data source
func (ds DataSource) ExecuteSQL(sql string, filter Geometry, dialect string) Layer {
	cSQL := C.CString(sql)
	defer C.free(unsafe.Pointer(cSQL))
	cDialect := C.CString(dialect)
	defer C.free(unsafe.Pointer(cDialect))

	layer := C.OGR_DS_ExecuteSQL(ds.cval, cSQL, filter.cval, cDialect)
	return Layer{layer}
}

// Release the results of ExecuteSQL
func (ds DataSource) ReleaseResultSet(layer Layer) {
	C.OGR_DS_ReleaseResultSet(ds.cval, layer.cval)
}

// Flush pending changes to the data source
func (ds DataSource) Sync() error {
	return C.OGR_DS_SyncToDisk(ds.cval).Err()
}

/* -------------------------------------------------------------------- */
/*      Driver functions                                                */
/* -------------------------------------------------------------------- */

type OGRDriver struct {
	cval C.OGRSFDriverH
}

// Fetch name of driver (file format)
func (driver OGRDriver) Name() string {
	name := C.OGR_Dr_GetName(driver.cval)
	return C.GoString(name)
}

// Attempt to open file with this driver
func (driver OGRDriver) Open(filename string, update int) (newDS DataSource, ok bool) {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	ds := C.OGR_Dr_Open(driver.cval, cFilename, C.int(update))
	return DataSource{ds}, ds != nil
}

// Test if this driver supports the named capability
func (driver OGRDriver) TestCapability(capability string) bool {
	cString := C.CString(capability)
	defer C.free(unsafe.Pointer(cString))
	val := C.OGR_Dr_TestCapability(driver.cval, cString)
	return val != 0
}

// Create a new data source based on this driver
func (driver OGRDriver) Create(name string, options []string) (newDS DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	ds := C.OGR_Dr_CreateDataSource(driver.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	return DataSource{ds}, ds != nil
}

// Create a new datasource with this driver by copying all layers of the existing datasource
func (driver OGRDriver) Copy(source DataSource, name string, options []string) (newDS DataSource, ok bool) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
		defer C.free(unsafe.Pointer(opts[i]))
	}
	opts[length] = (*C.char)(unsafe.Pointer(nil))

	ds := C.OGR_Dr_CopyDataSource(driver.cval, source.cval, cName, (**C.char)(unsafe.Pointer(&opts[0])))
	return DataSource{ds}, ds != nil
}

// Delete a data source
func (driver OGRDriver) Delete(filename string) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	return C.OGR_Dr_DeleteDataSource(driver.cval, cFilename).Err()
}

// Add a driver to the list of registered drivers
func (driver OGRDriver) Register() {
	C.OGRRegisterDriver(driver.cval)
}

// Remove a driver from the list of registered drivers
func (driver OGRDriver) Deregister() {
	C.OGRDeregisterDriver(driver.cval)
}

// Fetch the number of registered drivers
func OGRDriverCount() int {
	count := C.OGRGetDriverCount()
	return int(count)
}

// Fetch the indicated driver by index
func OGRDriverByIndex(index int) OGRDriver {
	driver := C.OGRGetDriver(C.int(index))
	return OGRDriver{driver}
}

// Fetch the indicated driver by name
func OGRDriverByName(name string) OGRDriver {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	driver := C.OGRGetDriverByName(cName)
	return OGRDriver{driver}
}

/* -------------------------------------------------------------------- */
/*      Style manager functions                                         */
/* -------------------------------------------------------------------- */

type StyleMgr struct {
	cval C.OGRStyleMgrH
}

type StyleTool struct {
	cval C.OGRStyleToolH
}

type StyleTable struct {
	cval C.OGRStyleTableH
}

// Unimplemented: CreateStyleManager

// Unimplemented: Destroy

// Unimplemented: InitFromFeature

// Unimplemented: InitStyleString

// Unimplemented: PartCount

// Unimplemented: PartCount

// Unimplemented: AddPart

// Unimplemented: AddStyle

// Unimplemented: CreateStyleTool

// Unimplemented: Destroy

// Unimplemented: Type

// Unimplemented: Unit

// Unimplemented: SetUnit

// Unimplemented: ParamStr

// Unimplemented: ParamNum

// Unimplemented: ParamDbl

// Unimplemented: SetParamStr

// Unimplemented: SetParamNum

// Unimplemented: SetParamDbl

// Unimplemented: StyleString

// Unimplemented: RGBFromString

// Unimplemented: CreateStyleTable

// Unimplemented: Destroy

// Unimplemented: Save

// Unimplemented: Load

// Unimplemented: Find

// Unimplemented: ResetStyleStringReading

// Unimplemented: NextStyle

// Unimplemented: LastStyleName
