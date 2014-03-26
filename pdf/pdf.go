package pdf

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//#include <cairo/cairo-pdf.h>
import "C"

import (
	"io"
	"unsafe"

	"github.com/jimmyfrasche/cairo"
)

type Surface struct {
	cairo.ExtensionPagedVectorSurface
	//w is used in NewWriter to ensure a reference to the writer lives as long as we do
	w io.Writer
}

func news(s *C.cairo_surface_t) (Surface, error) {
	S := Surface{
		ExtensionPagedVectorSurface: cairo.ExtensionNewPagedVectorSurface(s),
	}
	return S, S.Err()
}

func cNew(s *C.cairo_surface_t) (cairo.Surface, error) {
	return news(s)
}

func init() {
	cairo.ExtensionRegisterRawToSurface(cairo.SurfaceTypePDF, cNew)
}

func New(w io.Writer, width, height float64) (Surface, error) {
	wp := unsafe.Pointer(&w)
	pdf := C.cairo_pdf_surface_create_for_stream(cairo.ExtensionCairoWriteFuncT, wp, C.double(width), C.double(height))
	return news(pdf)
}

func NewFile(filename string, width, height float64) (Surface, error) {
	s := C.CString(filename)
	pdf := C.cairo_pdf_surface_create(s, C.double(width), C.double(height))
	C.free(unsafe.Pointer(s))
	return news(pdf)
}
