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

// Unimplemented: InitGCPs
// Unimplemented: DeinitGCPs
// Unimplemented: DuplicateGCPs
// Unimplemented: GCPsToGeoTransform
// Unimplemented: ApplyGeoTransform

