package script

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
//#include <cairo/cairo-script.h>
import "C"

import "github.com/jimmyfrasche/cairo"

//Surface is a script surface.
type Surface struct {
	cairo.XtensionPagedVectorSurface
}

func cNewSurf(s *C.cairo_surface_t) (Surface, error) {
	S := Surface{
		XtensionPagedVectorSurface: cairo.NewXtensionPagedVectorSurface(s),
	}
	return S, S.Err()
}

func revivSurf(s *C.cairo_surface_t) (cairo.Surface, error) {
	S, err := cNewSurf(s)
	return S, err
}

//NewSurface creates a script surface that records to d.
//
//Originally cairo_script_surface_create.
func (d Device) NewSurface(content cairo.Content, width, height float64) (Surface, error) {
	c := C.cairo_content_t(content)
	w, h := C.double(width), C.double(height)
	return cNewSurf(C.cairo_script_surface_create(d.XtensionRaw(), c, w, h))
}

func init() {
	cairo.XtensionRegisterRawToSurface(cairo.SurfaceTypeScript, revivSurf)
}
