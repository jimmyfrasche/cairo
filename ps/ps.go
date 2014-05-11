//Package ps implements the PostScript backend for libcairo rendering.
//
//EPS
//
//EPS files must contain only one page.
//
//Note that libcairo does not include the device independent preview.
//
//The Encapsulated PostScript Specfication† describes how to embed EPS
//files into PostScript files
//
//† https://partners.adobe.com/public/developer/en/ps/5002.EPSF_Spec.pdf
//
//Fonts
//
//The PostScript surface natively supports Type 1 and TrueType fonts.
//All other font types (including OpenType/PS) are converted to Type 1.
//
//Fonts are always subsetted and embedded.
//
//Fallback images
//
//Cairo will ensure that the PostScript output looks the same as an image
//surface for the same set of operations.
//When cairo drawing operations are performed that cannot be represented
//natively in PostScript, the drawing is rasterized and embedded in the output.
//
//The rasterization of unsupported operations is limited to the smallest
//rectangle, or set of rectangles, required to draw the unsupported operations.
//
//Fallback images are the main cause of large file sizes
//and slow printing times.
//Fallback images have a comment containing the size and location
//of the fallback.
//	> grep 'Fallback' output.ps
//	output.ps:% Fallback Image: x=100, y=100, w=299, h=50 res=300dpi size=783750
//	output.ps:% Fallback Image: x=100, y=150, w=350, h=250 res=300dpi size=4560834
//	output.ps:% Fallback Image: x=150, y=400, w=299, h=50 res=300dpi size=783750
//
//Supported Features
//
//The following tables lists all features natively supported by the PostScript
//surface.
//
//	FEATURE                          NOTES
//	Paint/Fill/Stroke/ShowGlyphs     depending on pattern
//	Fonts                            some may be converted to Type 1
//	Opaque colors
//	Images
//	Linear gradients                 level 3 only
//	Radial gradients                 level 3 and when one circle is inside
//	                                 the other and extent is ExtendNone or
//	                                 ExtendPad only.
//	PushGroup/CreateSimilar
//	OpSource
//	OpOver
//
//
//Requirements
//
//Libcairo must be compiled with
//	CAIRO_HAS_PS_SURFACE
//in addition to the requirements of cairo.
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
	cairo.XtensionPagedVectorSurface
	eps bool
}

func news(s *C.cairo_surface_t, eps bool) (Surface, error) {
	S := Surface{
		XtensionPagedVectorSurface: cairo.NewXtensionPagedVectorSurface(s),
		eps: eps,
	}
	return S, S.Err()
}

func cNew(s *C.cairo_surface_t) (cairo.Surface, error) {
	eps := C.cairo_ps_surface_get_eps(s) == 1
	//Note that if the surface was created with a Writer we have no way of
	//getting it here but that's okay as long as the original reference lives on.
	S, err := news(s, eps)
	return S, err
}

func init() {
	cairo.XtensionRegisterRawToSurface(cairo.SurfaceTypePS, cNew)
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
//W is the Writer the PostScript is written to.
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

	wp := cairo.XtensionWrapWriter(w)
	ps := C.cairo_ps_surface_create_for_stream(cairo.XtensionCairoWriteFuncT, wp, C.double(width), C.double(height))

	cfgSurf(ps, eps, header, setup)

	S, err = news(ps, eps)

	S.XtensionRegisterWriter(wp)

	return S, err
}

//RestrictTo restricts the generated PostScript to the specified level.
//The default is Level3.
//
//This method should only be called before any drawing operations have been
//performed on this surface.
//
//Originally cairo_ps_surface_restrict_to_level.
func (s Surface) RestrictTo(level level) error {
	C.cairo_ps_surface_restrict_to_level(s.XtensionRaw(), level.c())
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
	C.cairo_ps_surface_set_size(s.XtensionRaw(), C.double(width), C.double(height))
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
		C.cairo_ps_surface_dsc_comment(s.XtensionRaw(), str)
		C.free(unsafe.Pointer(str))
	}

	return s.Err()
}

//AddComment is shorthand for adding a single comment.
func (s Surface) AddComment(key, value string) error {
	return s.AddComments(Comments{Comment(key, value)})
}

//AddCommentf is shorthand for adding a single formatted comment.
func (s Surface) AddCommentf(key, value string, vars ...interface{}) error {
	return s.AddComments(Comments{Comment(key, value)})
}
