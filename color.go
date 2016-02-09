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

type ColorTable struct {
	cval C.GDALColorTableH
}

type ColorEntry struct {
	cval *C.GDALColorEntry
}

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
