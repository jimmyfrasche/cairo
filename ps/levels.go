package ps

//#cgo pkg-config: cairo
//#include <cairo/cairo-ps.h>
import "C"

import (
	"unsafe"
)

//cairo_ps_level_t
type level int

//The level type is used to describe the language level of the PostScript
//Language Reference that a generated PostScript file will conform to.
//
//While language level 3 supports additional features,
//such as gradient patterns, language level 2 printers cannot print PostScript
//containing language level 3 features.
//
//Note that when using Level3 the LanguageLevel DSC comment in the output may
//still indicate 2 if no level 3 features are used.
//
//Originally cairo_ps_level_t.
const (
	//Level2 is the language level 2 of the PostScript specification.
	Level2 level = C.CAIRO_PS_LEVEL_2
	//Level3 is the language level 3 of the PostScript specification.
	Level3 level = C.CAIRO_PS_LEVEL_3
)

func (l level) c() C.cairo_ps_level_t {
	return C.cairo_ps_level_t(l)
}

func (l level) String() string {
	v := C.cairo_ps_level_to_string(C.cairo_ps_level_t(l))
	if v == nil {
		return "unknown PS level"
	}
	return C.GoString(v)
}

//Levels reports the supported language levels.
//
//Originally cairo_ps_get_levels.
func Levels() (levels []level) {
	var lvls *C.cairo_ps_level_t
	var N C.int

	C.cairo_ps_get_levels(&lvls, &N)
	n := int(N)

	levels = make([]level, n)
	pseudoslice := (*[1 << 30]C.cairo_ps_level_t)(unsafe.Pointer(lvls))[:n:n]
	for i, v := range pseudoslice {
		levels[i] = level(v)
	}

	return
}
