//Package ps implements the PostScript backend for libcairo rendering.
package ps

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//#include <cairo/cairo-ps.h>
import "C"

import (
	"io"
	"unsafe"

	"github.com/jimmyfrasche/cairo"
)

//Surface is a PostScript surface.
//
//Surface implements cairo.PagedVectorSurface
type Surface struct {
	cairo.ExtensionPagedVectorSurface
	//w is used in NewWriter to ensure a reference to the writer lives as long as we do
	w   io.Writer
	eps bool
}

func news(s *C.cairo_surface_t, eps bool) (Surface, error) {
	S := Surface{
		ExtensionPagedVectorSurface: cairo.ExtensionNewPagedVectorSurface(s),
		eps: eps,
	}
	return S, S.Err()
}

func cNew(s *C.cairo_surface_t) (cairo.Surface, error) {
	eps := C.cairo_ps_surface_get_eps(s) == 1
	//Note that if the surface was created with an io.Writer we have no way of
	//getting it here but that's okay as long as the original reference lives on.
	S, err := news(s, eps)
	return S, err
}

func init() {
	cairo.ExtensionRegisterRawToSurface(cairo.SurfaceTypePS, cNew)
}

func cfgSurf(ps *C.cairo_surface_t, eps bool, header, setup Comments) {
	var ceps C.cairo_bool_t
	if eps {
		ceps = 1
	}
	C.cairo_ps_surface_set_eps(ps, ceps)

	for _, c := range header {
		s := C.CString(c.String())
		C.cairo_ps_surface_dsc_comment(ps, s)
		C.free(unsafe.Pointer(s))
	}

	if len(setup) > 0 {
		C.cairo_ps_surface_dsc_begin_setup(ps)

		for _, c := range setup {
			s := C.CString(c.String())
			C.cairo_ps_surface_dsc_comment(ps, s)
			C.free(unsafe.Pointer(s))
		}
	}

	//ensure calls to AddComments apply to page even if no drawing has been performed.
	C.cairo_ps_surface_dsc_begin_page_setup(ps)
}

func errChk(header, setup Comments) (err error) {
	if err = header.Err(); err != nil {
		return
	}
	return setup.Err()
}

//New creates a new PostScript of the specified size.
//
//W is the io.Writer the PostScript is written to.
//Width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//Eps specifies whether this will be Encapsulated PostScript.
//Header is any DSC comments to apply to the header section.
//Setup is any DSC comment to apply to the setup section.
//
//Originally cairo_ps_surface_create_for_stream
//and cairo_ps_surface_set_eps and cairo_ps_surface_dsc_comment
//and cairo_ps_surface_dsc_begin_setup and
//cairo_ps_surface_dsc_begin_page_setup.
func New(w io.Writer, width, height float64, eps bool, header, setup Comments) (S Surface, err error) {
	if err = errChk(header, setup); err != nil {
		return
	}

	wp := unsafe.Pointer(&w)
	ps := C.cairo_ps_surface_create_for_stream(cairo.ExtensionCairoWriteFuncT, wp, C.double(width), C.double(height))

	cfgSurf(ps, eps, header, setup)

	return news(ps, eps)
}

//NewFile creates a new PostScript of the specified size.
//
//Filename is the file the PostScript is written to.
//Width and height are in the unit of a typographical point
//(1 point = 1/72 inch).
//Eps specifies whether this will be Encapsulated PostScript.
//Header is any DSC comments to apply to the header section.
//Setup is any DSC comment to apply to the setup section.
//
//Originally cairo_ps_surface_create
//and cairo_ps_surface_set_eps and cairo_ps_surface_dsc_comment
//and cairo_ps_surface_dsc_begin_setup and
//cairo_ps_surface_dsc_begin_page_setup.
func NewFile(filename string, width, height float64, eps bool, header, setup Comments) (S Surface, err error) {
	if err = errChk(header, setup); err != nil {
		return
	}

	s := C.CString(filename)
	ps := C.cairo_ps_surface_create(s, C.double(width), C.double(height))
	C.free(unsafe.Pointer(s))

	cfgSurf(ps, eps, header, setup)

	return news(ps, eps)
}

//BUG(jmf): filename's CString not freed. Does cairo take control of it or strdup it?

//RestrictTo restricts the generated PostScript to the specified level.
//
//This method should only be called before any drawing operations have been
//performed on this surface.
//
//Originally cairo_ps_surface_restrict_to_level.
func (s Surface) RestrictTo(level level) error {
	C.cairo_ps_surface_restrict_to_level(s.ExtensionRaw(), level.c())
	return s.Err()
}

//SetSize changes the size of the PostScript surface for the current
//and subsequent pages.
//
//This method should only be called before any drawing operations have
//been performed on the current page.
//
//Originally cairo_ps_surface_set_size.
func (s Surface) SetSize(width, height float64) error {
	C.cairo_ps_surface_set_size(s.ExtensionRaw(), C.double(width), C.double(height))
	return s.Err()
}

//EPS reports whether s is Encapsulated PostScript.
//
//Originally cairo_ps_surface_get_eps.
func (s Surface) EPS() bool {
	return s.eps
}

//AddComments adds comments to the PageSetup sections.
//
//Originally cairo_ps_surface_dsc_comment.
func (s Surface) AddComments(comments Comments) (err error) {
	if err = comments.Err(); err != nil {
		return
	}

	for _, c := range comments {
		str := C.CString(c.String())
		C.cairo_ps_surface_dsc_comment(s.ExtensionRaw(), str)
		C.free(unsafe.Pointer(str))
	}

	return s.Err()
}
