package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

//Version returns the version of libcairo.
func Version() string {
	return C.GoString(C.cairo_version_string())
}
