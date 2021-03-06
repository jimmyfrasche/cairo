//Package pdf implements the PDF backend for libcairo rendering.
//
//Libcairo must be compiled with
//	CAIRO_HAS_PDF_SURFACE
//in addition to the requirements of cairo.
package pdf

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//#include <cairo/cairo-pdf.h>
import "C"

import (
	"io"

	"github.com/jimmyfrasche/cairo"
)

//Surface is a PDF surface.
//
//Surface implements cairo.PagedVectorSurface.
type Surface struct {
	cairo.XtensionPagedVectorSurface
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
//W is the Writer the PDF is written to.
//Width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//
//Originally cairo_pdf_surface_create_for_stream.
func New(w io.Writer, width, height float64) (Surface, error) {
	wp := cairo.XtensionWrapWriter(w)
	pdf := C.cairo_pdf_surface_create_for_stream(cairo.XtensionCairoWriteFuncT, wp, C.double(width), C.double(height))
	S, err := news(pdf)
	S.XtensionRegisterWriter(wp)
	return S, err
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
