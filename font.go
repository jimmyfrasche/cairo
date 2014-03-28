package cairo

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
import "C"

import (
	"runtime"
	"unsafe"
)

//FontOptions specify how fonts should be rendered.
//Most of the time the font options implied by a surface are just right
//and do not need changes, but, for pixel-based targets, tweaking font options
//may result in superior output on a particular display.
//
//Originally cairo_font_options_t.
type FontOptions struct {
	fo *C.cairo_font_options_t
}

//NewFontOptions creates a new FontOptions with the default values.
//
//Originally cairo_font_options_create.
func NewFontOptions() *FontOptions {
	fo := &FontOptions{C.cairo_font_options_create()}
	runtime.SetFinalizer(fo, (*FontOptions).Close)
	return fo
}

func initFontOptions(opt *C.cairo_font_options_t) *FontOptions {
	fo := &FontOptions{opt}
	runtime.SetFinalizer(fo, (*FontOptions).Close)
	return fo
}

//Close destroys the FontOptions. Close is idempotent.
//
//Originally cairo_font_options_destroy.
func (f *FontOptions) Close() error {
	if f == nil || f.fo == nil {
		return nil
	}
	C.cairo_font_options_destroy(f.fo)
	f.fo = nil
	runtime.SetFinalizer(f, nil)
	return nil
}

//Error queries f to see if there is an error.
//
//Originally cairo_font_options_status.
func (f *FontOptions) Err() error {
	return toerr(C.cairo_font_options_status(f.fo))
}

//Merge merges non-default options from o into f and return f.
//
//Originally cairo_font_options_merge.
func (f *FontOptions) Merge(o *FontOptions) *FontOptions {
	C.cairo_font_options_merge(f.fo, o.fo)
	return f
}

//Clone creates a new FontOptions with the same values as f.
//
//Originally cairo_font_options_copy.
func (f *FontOptions) Clone() *FontOptions {
	return &FontOptions{C.cairo_font_options_copy(f.fo)}
}

//Equal compares f with o.
//
//Originally cairo_font_options_equal.
func (f *FontOptions) Equal(o *FontOptions) bool {
	return C.cairo_font_options_equal(f.fo, o.fo) == 1
}

//SetAntialiasMode sets the antialiasing mode of f and returns f.
//
//Originally cairo_font_options_set_antialias.
func (f *FontOptions) SetAntialiasMode(a antialias) *FontOptions {
	C.cairo_font_options_set_antialias(f.fo, C.cairo_antialias_t(a))
	return f
}

//AntialiasMode reports the antialiasing mode of f.
//
//Originally cairo_font_topns_get_antialias.
func (f *FontOptions) AntialiasMode() antialias {
	return antialias(C.cairo_font_options_get_antialias(f.fo))
}

//SetSubpixelOrder sets the subpixel ordering of f and returns f.
//
//Originally cairo_font_options_set_subpixel_order.
func (f *FontOptions) SetSubpixelOrder(s subpixelOrder) *FontOptions {
	C.cairo_font_options_set_subpixel_order(f.fo, C.cairo_subpixel_order_t(s))
	return f
}

//SubpixelOrder reports the subpixel ordering of f.
//
//Originally cairo_font_options_get_subpixel_order.
func (f *FontOptions) SubpixelOrder() subpixelOrder {
	return subpixelOrder(C.cairo_font_options_get_subpixel_order(f.fo))
}

//SetHintStyle sets the hint style of f and returns f.
//
//Originally cairo_font_options_set_hint_style.
func (f *FontOptions) SetHintStyle(h hintStyle) *FontOptions {
	C.cairo_font_options_set_hint_style(f.fo, C.cairo_hint_style_t(h))
	return f
}

//HintStyle reports the hint style of f.
//
//Originally cairo_font_options_get_hint_style.
func (f *FontOptions) HintStyle() hintStyle {
	return hintStyle(C.cairo_font_options_get_hint_style(f.fo))
}

//SetHintMetrics sets the hint metrics of f and returns f.
//
//Originally cairo_font_options_set_hint_metrics.
func (f *FontOptions) SetHintMetrics(h hintMetrics) *FontOptions {
	C.cairo_font_options_set_hint_metrics(f.fo, C.cairo_hint_metrics_t(h))
	return f
}

//HintMetrics reports the hint metrics of f.
//
//Originally cairo_font_options_get_hint_metrics.
func (f *FontOptions) HintMetrics() hintMetrics {
	return hintMetrics(C.cairo_font_options_get_hint_metrics(f.fo))
}

//FontExtents stores metric information for a font.
//Values are given in the current user-space coordinate system.
//
//Because font metrics are in user-space coordinates, they are mostly,
//but not entirely, independent of the current transformation matrix.
//They will, however, change slightly due to hinting but otherwise remain
//unchanged.
//
//Originally cairo_font_extents_t.
type FontExtents struct {
	//Ascent is the distance the font extends above the baseline.
	Ascent float64
	//Descent is the distance the font extends below the baseline.
	Descent float64
	//Height is the recommended vertical distance between baselines
	//when setting consecutive lines of text with the font.
	Height float64
	//MaxAdvanceX is the maximum distance in the X direction that
	//the origin is advanced for any glyph in the font.
	//
	//Originally max_y_advance.
	MaxAdvanceX float64
	//MaxAdvanceY is the maximum distance in the Y direction that
	//the origin is advanced for any glyph in the font.
	//
	//This will be zero for most fonts used for horizontal writing.
	//
	//Originally max_x_advance.
	MaxAdvanceY float64
}

func newFontExtents(fe C.cairo_font_extents_t) FontExtents {
	return FontExtents{
		float64(fe.ascent),
		float64(fe.descent),
		float64(fe.height),
		float64(fe.max_x_advance),
		float64(fe.max_y_advance),
	}
}

//ExternalLeading reports the difference between the Height and the sum
//of the Ascent and Descent. Also known as "line spacing".
func (f FontExtents) ExternalLeading() float64 {
	return f.Height - (f.Ascent + f.Descent)
}

//TextExtents stores the extents of a single glyph or string of glyphs
//in user-space coordinates.
//Because text extents are in user-space coordinates, they are mostly,
//but not entirely, independent of the current transformation matrix.
//They will, however, change slightly due to hinting.
//
//Originally cairo_text_extents_t.
type TextExtents struct {
	//The horizontal distance from the origin to the leftmost part of the glyhps
	//as drawn.
	//Positive if the glyphs lie entirely to the right of the origin.
	//
	//Originally x_bearing.
	BearingX float64
	//The vertical distance from the origin to the topmost part of the glyhps
	//as drawn.
	//Positive if the glyphs lie entirely below the origin.
	//
	//Originally y_bearing.
	BearingY float64
	//Width of the glyphs as drawn.
	Width float64
	//Height of the glyphs as drawn.
	Height float64
	//AdvanceX is the distance in the X direction to advance after drawing
	//these glyphs.
	//
	//Originally x_advance.
	AdvanceX float64
	//AdvanceY is the distance in the Y direction to advance after drawing
	//these glyphs.
	//
	//This will be zero for most fonts used for horizontal writing.
	//
	//Originally y_advance.
	AdvanceY float64
}

func newTextExtents(te C.cairo_text_extents_t) TextExtents {
	return TextExtents{
		BearingX: float64(te.x_bearing),
		BearingY: float64(te.y_bearing),
		Width:    float64(te.width),
		Height:   float64(te.height),
		AdvanceX: float64(te.x_advance),
		AdvanceY: float64(te.y_advance),
	}
}

//Glyph holds information about a single glyph when drawing or measuring text.
//A font is (in simple terms) a collection of shapes used to draw text.
//A glyph is one of these shapes.
//There can be multiple glyphs for a single character (alternates to be used
//in different contexts, for example), or a glyph can be a ligature of multiple
//characters.
//Cairo doesn't expose any way of converting input text into glyphs,
//so in order to use the Cairo interfaces that take arrays of glyphs, you must
//directly access the appropriate underlying font system.
//
//Note that the offsets given by Point.X and Point.Y are not cumulative.
//When drawing or measuring text, each glyph is individually positioned with respect to the overall origin
//
//Originally cairo_glyph_t.
type Glyph struct {
	Index uint64
	Point Point
}

func cGlyph(g C.cairo_glyph_t) Glyph {
	return Glyph{
		Index: uint64(g.index),
		Point: cPt(g.x, g.y),
	}
}

func cGlyphs(glyphs *C.cairo_glyph_t, N C.int) []Glyph {
	n := int(N)
	gs := (*[1 << 30]C.cairo_glyph_t)(unsafe.Pointer(glyphs))[:n:n]
	out := make([]Glyph, n)
	for i, v := range gs {
		out[i] = cGlyph(v)
	}
	return out
}

func glyphsC(gs []Glyph) (glyphs *C.cairo_glyph_t, N C.int) {
	n := len(gs)
	N = C.int(n)
	var t C.cairo_glyph_t
	glyphs = (*C.cairo_glyph_t)(C.malloc(C.size_t(uintptr(n) * unsafe.Sizeof(t))))
	iter := (*[1 << 30]C.cairo_glyph_t)(unsafe.Pointer(glyphs))[:n:n]
	for i, g := range gs {
		iter[i] = g.c()
	}
	return
}

func (g Glyph) c() C.cairo_glyph_t {
	out := C.cairo_glyph_t{}
	out.index = C.ulong(g.Index)
	out.x, out.y = g.Point.c()
	return out
}

type Font interface {
	Type() fontType
	Close() error
	Err() error
	//TODO
	XtensionRaw() *C.cairo_font_face_t
}

type XtensionFont struct {
	f *C.cairo_font_face_t
}

func cFont(f *C.cairo_font_face_t) (Font, error) {
	return nil, nil //TODO
}

type ScaledFont struct {
	//TODO
	f *C.cairo_scaled_font_t
}

func cNewScaledFont(f *C.cairo_scaled_font_t) *ScaledFont {
	return nil //TODO
}

//Originally cairo_text_cluster_t.
type TextCluster struct {
	RuneLength, NumGlyphs int
}

func clustersC(tcs []TextCluster) (*C.cairo_text_cluster_t, C.int) {
	return nil, 0 //TODO
}
