package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"runtime"
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
func NewFontOptions() FontOptions {
	fo := FontOptions{C.cairo_font_options_create()}
	runtime.SetFinalizer(fo, FontOptions.Close)
	return fo
}

//Close destroys the FontOptions. Close is idempotent.
//
//Originally cairo_font_options_destroy.
func (f FontOptions) Close() error {
	if f.fo == nil {
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
func (f FontOptions) Error() error {
	return toerr(C.cairo_font_options_status(f.fo))
}

//Merge merges non-default options from o into f and return f.
//
//Originally cairo_font_options_merge.
func (f FontOptions) Merge(o FontOptions) FontOptions {
	C.cairo_font_options_merge(f.fo, o.fo)
	return f
}

//Clone creates a new FontOptions with the same values as f.
//
//Originally cairo_font_options_copy.
func (f FontOptions) Clone() FontOptions {
	return FontOptions{C.cairo_font_options_copy(f.fo)}
}

//Equal compares f with o.
//
//Originally cairo_font_options_equal.
func (f FontOptions) Equal(o FontOptions) bool {
	return C.cairo_font_options_equal(f.fo, o.fo) == 1
}

//SetAntialiasMode sets the antialiasing mode of f and returns f.
//
//Originally cairo_font_options_set_antialias.
func (f FontOptions) SetAntialiasMode(a antialias) FontOptions {
	C.cairo_font_options_set_antialias(f.fo, C.cairo_antialias_t(a))
	return f
}

//AntialiasMode reports the antialiasing mode of f.
//
//Originally cairo_font_topns_get_antialias.
func (f FontOptions) AntialiasMode() antialias {
	return antialias(C.cairo_font_options_get_antialias(f.fo))
}

//SetSubpixelOrder sets the subpixel ordering of f and returns f.
//
//Originally cairo_font_options_set_subpixel_order.
func (f FontOptions) SetSubpixelOrder(s subpixelOrder) FontOptions {
	C.cairo_font_options_set_subpixel_order(f.fo, C.cairo_subpixel_order_t(s))
	return f
}

//SubpixelOrder reports the subpixel ordering of f.
//
//Originally cairo_font_options_get_subpixel_order.
func (f FontOptions) SubpixelOrder() subpixelOrder {
	return subpixelOrder(C.cairo_font_options_get_subpixel_order(f.fo))
}

//SetHintStyle sets the hint style of f and returns f.
//
//Originally cairo_font_options_set_hint_style.
func (f FontOptions) SetHintStyle(h hintStyle) FontOptions {
	C.cairo_font_options_set_hint_style(f.fo, C.cairo_hint_style_t(h))
	return f
}

//HintStyle reports the hint style of f.
//
//Originally cairo_font_options_get_hint_style.
func (f FontOptions) HintStyle() hintStyle {
	return hintStyle(C.cairo_font_options_get_hint_style(f.fo))
}

//SetHintMetrics sets the hint metrics of f and returns f.
//
//Originally cairo_font_options_set_hint_metrics.
func (f FontOptions) SetHintMetrics(h hintMetrics) FontOptions {
	C.cairo_font_options_set_hint_metrics(f.fo, C.cairo_hint_metrics_t(h))
	return f
}

//HintMetrics reports the hint metrics of f.
//
//Originally cairo_font_options_get_hint_metrics.
func (f FontOptions) HintMetrics() hintMetrics {
	return hintMetrics(C.cairo_font_options_get_hint_metrics(f.fo))
}
