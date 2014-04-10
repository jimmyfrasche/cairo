//Package svg implements the SVG backend for libcairo rendering.
package svg

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//#include <cairo/cairo-svg.h>
import "C"

import (
	"github.com/jimmyfrasche/cairo"
)

//Surface is an SVG surface.
//
//Surface implements cairo.VectorSurface.
type Surface struct {
	cairo.XtensionVectorSurface
}

//New creates a new SVG surface that writes to writer.
//
//If writer needs to be flushed or closed, that is the responsibility
//of the caller.
//
//The parameters width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//
//Originally cairo_svg_surface_create_for_stream.
func New(w cairo.Writer, width, height float64) (Surface, error) {
	wp := cairo.XtensionWrapWriter(w)
	svg := C.cairo_svg_surface_create_for_stream(cairo.XtensionCairoWriteFuncT, wp, C.double(width), C.double(height))
	s := Surface{
		XtensionVectorSurface: cairo.NewXtensionVectorSurface(svg),
	}
	s.XtensionRegisterWriter(wp)
	return s, s.Err()
}

func cNew(s *C.cairo_surface_t) (cairo.Surface, error) {
	//Note that if the surface was created with a Writer we have no way of
	//getting it here but that's okay as long as the original reference lives on.
	S := Surface{
		XtensionVectorSurface: cairo.NewXtensionVectorSurface(s),
	}
	return S, S.Err()
}

func init() {
	cairo.XtensionRegisterRawToSurface(cairo.SurfaceTypeSVG, cNew)
}

//BUG(jmf): add documentation about mime type to New after I figure out what that entails.

//BUG(jmf): add special method for attaching url's as that's handled specially anyhoo.

//RestrictTo restricts the generated SVG file to the specified version.
//
//This method should only be called before any drawing operations have been
//performed on this surface.
//
//Originally cairo_svg_surface_restrict_to_version.
func (s Surface) RestrictTo(v version) error {
	C.cairo_svg_surface_restrict_to_version(s.XtensionRaw(), v.c())
	return s.Err()
}
