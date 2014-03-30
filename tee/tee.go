//Package tee implements a surface that multiplexes all operations performed
//on it to the one or more underlying surfaces.
package tee

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
//#include <cairo/cairo-tee.h>
import "C"

import (
	"github.com/jimmyfrasche/cairo"
)

//Surface is a tee surface.
//
//Surface is a paged vector surface.
//
//If any of the underlying surface are not paged or vector backed, those
//operations will be no-ops.
type Surface struct {
	cairo.XtensionPagedVectorSurface
}

//New creates a new tee surface.
//
//Originally cairo_tee_surface_create.
func New(masterSurface cairo.Surface, surfaces ...cairo.Surface) (Surface, error) {
	m := C.cairo_tee_surface_create(masterSurface.XtensionRaw())
	for _, s := range surfaces {
		C.cairo_tee_surface_add(m, s.XtensionRaw())
	}
	S := Surface{cairo.NewXtensionPagedVectorSurface(m)}
	return S, S.Err()
}

//Index returns the ith surface of this tee.
//
//The returned error is set if there's an error on the returned surface.
//If the index does not exist, Index returns (nil, nil).
//
//Originally cairo_tee_surface_index.
func (s Surface) Index(i int) (cairo.Surface, error) {
	if i < 0 {
		return nil, nil
	}
	sir := C.cairo_tee_surface_index(s.XtensionRaw(), C.uint(i))
	if C.cairo_surface_status(sir) == C.CAIRO_STATUS_INVALID_INDEX {
		return nil, nil
	}
	si, err := cairo.XtensionRevivifySurface(sir)
	if err != nil {
		return nil, err
	}
	return si, nil
}

//Remove removes rs surfaces from s.
//
//Originally cairo_tee_surface_remove.
func (s Surface) Remove(rs ...cairo.Surface) error {
	me := s.XtensionRaw()
	for _, r := range rs {
		C.cairo_tee_surface_remove(me, r.XtensionRaw())
	}
	return s.Err()
}

//Add adds as surfaces to s.
//
//Originally cairo_tee_surface_add.
func (s Surface) Add(as ...cairo.Surface) error {
	me := s.XtensionRaw()
	for _, a := range as {
		C.cairo_tee_surface_add(me, a.XtensionRaw())
	}
	return s.Err()
}
