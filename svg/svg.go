//Package svg implements the SVG backend for libcairo rendering.
package svg

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
//#include <cairo/cairo-svg.h>
import "C"

import (
	"io"
	"unsafe"

	"github.com/jimmyfrasche/cairo"
)

//Surface is an SVG surface.
//
//Surface implements cairo.VectorSurface.
//
//Originally cairo_svg_surface_t.
type Surface struct {
	cairo.ExtensionVectorSurface
	//w is used in NewWriter to ensure a reference to the writer lives as long as we do
	w io.Writer
}

//NewFile creates a new SVG surface that writes to filename.
//
//The parameters width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//
//Originally cairo_svg_surface_create.
func NewFile(filename string, width, height float64) (Surface, error) {
	nm := C.CString(filename)
	svg := C.cairo_svg_surface_create(nm, C.double(width), C.double(height))
	s := Surface{
		ExtensionVectorSurface: cairo.ExtensionNewVectorSurface(svg),
	}
	return s, s.Err()
}

//BUG(jmf): filename's CString not freed. Does cairo take control of it or strdup it?

//New creates a new SVG surface that writes to writer.
//
//If writer needs to be flushed or closed, that is the responsibility
//of the caller.
//
//The parameters width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//
//Originally cairo_svg_surface_create_for_stream.
func New(w io.Writer, width, height float64) (Surface, error) {
	wp := unsafe.Pointer(&w)
	svg := C.cairo_svg_surface_create_for_stream(cairo.ExtensionCairoWriteFuncT, wp, C.double(width), C.double(height))
	s := Surface{
		ExtensionVectorSurface: cairo.ExtensionNewVectorSurface(svg),
		w: w,
	}
	return s, s.Err()
}

func cNew(s *C.cairo_surface_t) (cairo.Surface, error) {
	//Note that if the surface was created with an io.Writer we have no way of
	//getting it here but that's okay as long as the original reference lives on.
	S := Surface{
		ExtensionVectorSurface: cairo.ExtensionNewVectorSurface(s),
	}
	return S, S.Err()
}

func init() {
	cairo.ExtensionRegisterRawToSurface(cairo.SurfaceTypeSVG, cNew)
}

//BUG(jmf): add documentation about mime type to New after I figure out what that entails.

//BUG(jmf): add special method for attaching url's as that's handled specially anyhoo.

//RestrictTo restricts the generated SVG file to the specified version.
//
//Originally cairo_svg_surface_restrict_to_version.
func (s Surface) RestrictTo(v version) error {
	C.cairo_svg_surface_restrict_to_version(s.ExtensionRaw(), v.c())
	return s.Err()
}
