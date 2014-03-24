package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"
import (
	"image/color"
	"runtime"
)

//Pattern is a pattern used for drawing.
//
//Originally cairo_pattern_t.
type Pattern interface {
	Type() patternType
	Err() error
	Close() error

	SetExtend(extend)
	Extend() extend
	SetFilter(filter)
	Filter() filter
	SetMatrix(Matrix)
	Matrix() Matrix
}

type pattern struct {
	p *C.cairo_pattern_t
}

func newPattern(p *C.cairo_pattern_t) *pattern {
	P := &pattern{p}
	runtime.SetFinalizer(p, (*pattern).Close)
	return P
}

//Type returns the type of the pattern.
//
//Originally cairo_pattern_get_type.
func (p *pattern) Type() patternType {
	return patternType(C.cairo_pattern_get_type(p.p))
}

//Err reports any error on this pattern.
//
//Originally cairo_pattern_status.
func (p *pattern) Err() error {
	return toerr(C.cairo_pattern_status(p.p))
}

//Close releases the pattern's resources.
//
//Originally cairo_pattern_destroy.
func (p *pattern) Close() error {
	if p.p == nil {
		return nil
	}
	runtime.SetFinalizer(p, nil)
	C.cairo_pattern_destroy(p.p)
	p.p = nil
	return nil
}

//SetExtend sets the mode used for drawing outside the area of this pattern.
//
//Originally cairo_pattern_set_extend.
func (p *pattern) SetExtend(e extend) {
	C.cairo_pattern_set_extend(p.p, e.c())
}

//Extend reports the mode used for drawing outside the area of this pattern.
//
//Originally cairo_pattern_get_extend.
func (p *pattern) Extend() extend {
	return extend(C.cairo_pattern_get_extend(p.p))
}

//SetFilter sets the filter used when resizing this pattern.
//
//Originally cairo_pattern_set_filter.
func (p *pattern) SetFilter(f filter) {
	C.cairo_pattern_set_filter(p.p, f.c())
}

//Filter returns the filter used when resizing this pattern.
//
//Originally cairo_pattern_get_filter.
func (p *pattern) Filter() filter {
	return filter(C.cairo_pattern_get_filter(p.p))
}

//SetMatrix sets the pattern's transformation matrix.
//This matrix is a transformation from user space to pattern space.
//
//When a pattern is first created it always has the identity matrix for its
//transformation matrix, which means that pattern space is initially identical
//to user space.
//
//Important
//
//Please note that the direction of this transformation matrix is from user
//space to pattern space.
//This means that if you imagine the flow from a pattern to user space
//(and on to device space), then coordinates in that flow will be transformed
//by the inverse of the pattern matrix.
//
//Originally cairo_pattern_set_matrix.
func (p pattern) SetMatrix(m Matrix) {
	C.cairo_pattern_set_matrix(p.p, &m.m)
}

//Matrix returns this patterns transformation matrix.
//
//Originally cairo_pattern_get_matrix.
func (p pattern) Matrix() Matrix {
	var m C.cairo_matrix_t
	C.cairo_pattern_get_matrix(p.p, &m)
	return Matrix{m}
}

//SolidPattern is a Pattern corresponding to a single translucent color.
type SolidPattern struct {
	*pattern
	c color.Color
}

//NewSolidPattern creates a solid pattern of color c.
//
//Originally cairo_pattern_create_rgba.
func NewSolidPattern(c color.Color) SolidPattern {
	c = AlphaColorModel.Convert(c).(AlphaColor).Canon()
	r, g, b, a := c.(AlphaColor).c()
	p := C.cairo_pattern_create_rgba(r, g, b, a)
	return SolidPattern{
		pattern: newPattern(p),
		c:       c,
	}
}

//Color returns the color this pattern was created with.
//Regardless of the type of color this pattern was created
//with the returned color will always be an AlphaColor.
//
//Originally cairo_pattern_get_rgba.
func (s SolidPattern) Color() color.Color {
	return s.c
}

//SurfacePattern is a Pattern backed by a Surface.
type SurfacePattern struct {
	*pattern
	s Surface
}

//NewSurfacePattern creates a Pattern from a Surface.
//
//Originally cairo_pattern_create_for_surface.
func NewSurfacePattern(s Surface) (sp SurfacePattern, err error) {
	if err = s.Err(); err != nil {
		return
	}
	r := s.ExtensionRaw()
	p := C.cairo_pattern_create_for_surface(r)
	sp = SurfacePattern{
		pattern: newPattern(p),
		s:       s,
	}
	return sp, sp.Err()
}

//BUG(jmf): (potentially) assuming pattern returned by cairo_pattern_create_for_surface
//is the same as the pattern put into it. If this is not true, things could get messy.

//Surface returns the Surface of this Pattern.
//
//Originally cairo_pattern_get_surface.
func (s SurfacePattern) Surface() Surface {
	return s.s
}

//Gradient is a linear or radial gradient.
type Gradient interface {
	Pattern
	AddColorStop(float64, color.Color)
	ColorStops() int
	ColorStop(int) (float64, color.Color, error)
}

type patternGradient struct {
	*pattern
}

//AddColorStop adds a color stop to the gradient.
//
//The offset specifies the location along the gradient's control vector.
//
//If two (or more) stops are specified with identical offset values,
//they will be sorted according to the order in which the stops are added.
//Stops added earlier will compare less than stops added later.
//This can be useful for reliably making sharp color transitions
//instead of the typical blend.
//
//Originally cairo_pattern_add_color_stop_rgb if c has no alpha channel or
//cairo_pattern_add_color_stop_rgba otherwise.
func (p patternGradient) AddColorStop(offset float64, c color.Color) {
	o := C.double(clamp01(offset))
	if _, _, _, a := c.RGBA(); a == 0xffff {
		c := ColorModel.Convert(c).(Color)
		r, g, b := c.Canon().c()
		C.cairo_pattern_add_color_stop_rgb(p.p, o, r, g, b)
	} else {
		c := AlphaColorModel.Convert(c).(AlphaColor)
		r, g, b, a := c.Canon().c()
		C.cairo_pattern_add_color_stop_rgba(p.p, o, r, g, b, a)
	}
}

//ColorStops returns the number of color stops in this gradient.
//
//Originally cairo_pattern_get_color_stop_count.
func (p patternGradient) ColorStops() int {
	var n C.int
	//only returns error if not a gradient, but disallowed by construction.
	_ = C.cairo_pattern_get_color_stop_count(p.p, &n)
	return int(n)
}

//ColorStop returns the nth color stop of this gradient.
//
//Originally cairo_pattern_get_color_stop_rgba.
func (p patternGradient) ColorStop(n int) (offset float64, color AlphaColor, err error) {
	var o, r, g, b, a C.double
	err = toerr(C.cairo_pattern_get_color_stop_rgba(p.p, C.int(n), &o, &r, &g, &b, &a))
	if err != nil {
		return
	}
	offset = float64(o)
	color = AlphaColor{
		R: float64(r),
		G: float64(g),
		B: float64(b),
		A: float64(a),
	}
	return
}

//LinearGradient is a linear gradient pattern.
type LinearGradient struct {
	patternGradient
	start, end Point
}

//NewLinearGradient creates a new linear gradient, from start to end.
//
//Originally cairo_pattern_create_linear.
func NewLinearGradient(start, end Point) LinearGradient {
	x0, y0 := start.c()
	x1, y1 := end.c()
	p := C.cairo_pattern_create_linear(x0, y0, x1, y1)
	return LinearGradient{
		patternGradient: patternGradient{
			pattern: newPattern(p),
		},
		start: start,
		end:   end,
	}
}

//Line returns the start and end points of this linear gradient.
//
//Originally cairo_pattern_get_linear_points.
func (l LinearGradient) Line() (start, end Point) {
	return l.start, l.end
}

//RadialGradient is a radial gradient pattern.
type RadialGradient struct {
	patternGradient
	start, end Circle
}

//NewRadialGradient creates a new radial gradient between the circles
//start and end.
//
//Originally cairo_pattern_create_radial.
func NewRadialGradient(start, end Circle) RadialGradient {
	x0, y0, r0 := start.c()
	x1, y1, r1 := end.c()
	p := C.cairo_pattern_create_radial(x0, y0, r0, x1, y1, r1)
	return RadialGradient{
		patternGradient: patternGradient{
			pattern: newPattern(p),
		},
		start: start,
		end:   end,
	}
}
