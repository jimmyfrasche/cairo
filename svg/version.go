package svg

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//#include <cairo/cairo-svg.h>
import "C"

import "unsafe"

type version int

//A version describes the version number of the SVG specification
//that a generated SVG file conforms to.
//
//Originally cairo_svg_version_t.
const (
	//SVG1_1 is the version 1.1 SVG specfication.
	SVG1_1 version = C.CAIRO_SVG_VERSION_1_1
	//SVG1_2 is the version 1.2 SVG specfication.
	SVG1_2 version = C.CAIRO_SVG_VERSION_1_2
)

func (v version) c() C.cairo_svg_version_t {
	return C.cairo_svg_version_t(v)
}

func (v version) String() string {
	V := C.cairo_svg_version_to_string(C.cairo_svg_version_t(v))
	if V == nil {
		return "unknown SVG version"
	}
	return C.GoString(V)
}

//Versions returns the SVG versions that libcairo supports.
//
//Originally cairo_svg_get_versions.
func Versions() (versions []version) {
	var vs *C.cairo_svg_version_t
	var N C.int

	C.cairo_svg_get_versions(&vs, &N)
	defer C.free(unsafe.Pointer(vs))

	n := int(N)
	pseudoslice := (*[1 << 30]C.cairo_svg_version_t)(unsafe.Pointer(vs))[:n:n]
	for _, v := range pseudoslice {
		versions = append(versions, version(v))
	}

	return
}
