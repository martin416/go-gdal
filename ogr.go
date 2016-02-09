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
