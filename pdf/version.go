package pdf

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
//#include <cairo/cairo-pdf.h>
import "C"

import (
	"unsafe"
)

//cairo_pdf_version_t
type version int

//The version type describes the version number of the PDF specification
//that a generated PDF file will conform to.
//
//Originally cairo_pdf_version_t.
const (
	//Version1_4 is the version 1.4 of the PDF specification.
	Version1_4 version = C.CAIRO_PDF_VERSION_1_4
	//Version1_5 is the version 1.5 of the PDF specification.
	Version1_5 version = C.CAIRO_PDF_VERSION_1_5
)

func (p version) String() string {
	v := C.cairo_pdf_version_to_string(C.cairo_pdf_version_t(p))
	if v == nil {
		return "unknown PDF version"
	}
	return C.GoString(v)
}

func Versions() (versions []version) {
	var vs *C.cairo_pdf_version_t
	var N C.int

	C.cairo_pdf_get_versions(&vs, &N)

	n := int(N)
	versions = make([]version, n)
	pseudoslice := (*[1 << 30]C.cairo_pdf_version_t)(unsafe.Pointer(vs))[:n:n]
	for i, v := range pseudoslice {
		versions[i] = version(v)
	}

	return
}
