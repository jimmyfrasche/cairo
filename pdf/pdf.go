//Package pdf implements the PDF backend for libcairo rendering.
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

//Surface is a PDF surface.
//
//Surface implements cairo.PagedVectorSurface.
type Surface struct {
	cairo.XtensionPagedVectorSurface
	//w is used in NewWriter to ensure a reference to the writer lives as long as we do
	w io.Writer
}

func news(s *C.cairo_surface_t) (Surface, error) {
	S := Surface{
		XtensionPagedVectorSurface: cairo.NewXtensionPagedVectorSurface(s),
	}
	return S, S.Err()
}

func cNew(s *C.cairo_surface_t) (cairo.Surface, error) {
	return news(s)
}

func init() {
	cairo.XtensionRegisterRawToSurface(cairo.SurfaceTypePDF, cNew)
}

//New creates a new PDF of the specified size.
//W is the io.Writer the PDF is written to.
//Width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//
//Warning
//
//It is the caller's responsibility to keep a reference to w for the lifetime
//of this surface.
//As it is passed to libcairo, the Go garbage collector will otherwise find
//no reference to it.
//
//Originally cairo_pdf_surface_create_for_stream.
func New(w io.Writer, width, height float64) (Surface, error) {
	wp := unsafe.Pointer(&w)
	pdf := C.cairo_pdf_surface_create_for_stream(cairo.XtensionCairoWriteFuncT, wp, C.double(width), C.double(height))
	return news(pdf)
}

//NewFile creates a new PDF of the specified size.
//Filename specifies the file the PDF is written to.
//Width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//
//Originally cairo_pdf_surface_create.
func NewFile(filename string, width, height float64) (Surface, error) {
	s := C.CString(filename)
	pdf := C.cairo_pdf_surface_create(s, C.double(width), C.double(height))
	C.free(unsafe.Pointer(s))
	return news(pdf)
}

//RestrictTo restricts the generated PDF to the specified verison.
//
//This method should only be called before any drawing operations have been
//performed on this surface.
//
//Originally cairo_pdf_surface_restrict_to_version.
func (s Surface) RestrictTo(v version) error {
	C.cairo_pdf_surface_restrict_to_version(s.XtensionRaw(), C.cairo_pdf_version_t(v))
	return s.Err()
}

//SetSize changes the size of the PDF surface for the current and subsequent
//pages.
//
//This method should only be called before any drawing operations have
//been performed on the current page.
//
//Originally cairo_pdf_surface_set_size.
func (s Surface) SetSize(width, height float64) error {
	C.cairo_pdf_surface_set_size(s.XtensionRaw(), C.double(width), C.double(height))
	return s.Err()
}
