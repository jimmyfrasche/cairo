//Package recording provides a surface that records drawing operations
//performed on it to later be replayed on another surface.
//
//Libcairo must be compiled with
//	CAIRO_HAS_RECORDING_SURFACE
//in addition to the requirements of cairo.
package recording

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"github.com/jimmyfrasche/cairo"
)

func rect(x, y, w, h C.double) cairo.Rectangle {
	return cairo.RectWH(float64(x), float64(y), float64(w), float64(h))
}

//Surface records all drawing operations performed against it.
//This recording can then be "replayed" against a target surface by using
//it as a source surface.
type Surface struct {
	cairo.XtensionPagedVectorSurface
	extents cairo.Rectangle
}

//New creates a new recording surface.
//
//If extents.Empty() is true, the recording surface is unbounded.
//
//Originally cairo_recording_surface_create.
func New(content cairo.Content, extents cairo.Rectangle) Surface {
	con := C.cairo_content_t(content)
	extents = extents.Canon()
	var s *C.cairo_surface_t
	if extents.Empty() {
		s = C.cairo_recording_surface_create(con, nil)
	} else {
		r := C.cairo_rectangle_t{
			x:      C.double(extents.Min.X),
			y:      C.double(extents.Min.Y),
			width:  C.double(extents.Dx()),
			height: C.double(extents.Dy()),
		}
		s = C.cairo_recording_surface_create(con, &r)
	}
	return cNew(s, extents)
}

func cNew(s *C.cairo_surface_t, e cairo.Rectangle) Surface {
	return Surface{
		cairo.NewXtensionPagedVectorSurface(s),
		e,
	}
}

func reviv(s *C.cairo_surface_t) (cairo.Surface, error) {
	var r C.cairo_rectangle_t
	C.cairo_recording_surface_get_extents(s, &r)
	e := rect(r.x, r.y, r.width, r.height)
	S := cNew(s, e)
	return S, S.Err()
}

func init() {
	cairo.XtensionRegisterRawToSurface(cairo.SurfaceTypeRecording, reviv)
}

//Extents reports the extents of this surface.
//If the surface is unbounded, then extents.Empty() is true.
//
//Originally cairo_recording_surface_get_extents.
func (s Surface) Extents() (extents cairo.Rectangle) {
	return s.extents
}

//Unbounded reports whether this surface is unbounded.
func (s Surface) Unbounded() bool {
	return s.extents.Empty()
}

//InkExtents measures the extents of the operations recorded on this surface.
//
//This is useful to compute the required size of a surface to replay
//the recorded drawing operations on.
//
//Originally cairo_recording_surface_ink_extents.
func (s Surface) InkExtents() cairo.Rectangle {
	var x, y, w, h C.double
	C.cairo_recording_surface_ink_extents(s.XtensionRaw(), &x, &y, &w, &h)
	return rect(x, y, w, h)
}
