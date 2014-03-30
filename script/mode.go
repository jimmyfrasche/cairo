package script

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
//#include <cairo/cairo-script.h>
import "C"

//cairo_script_mode_t
type mode int

//The mode type specifies the output mode of a script.
const (
	//ASCII is a human readable mode.
	ASCII mode = C.CAIRO_SCRIPT_MODE_ASCII
	//Binary specifies a byte-code based mode.
	Binary mode = C.CAIRO_SCRIPT_MODE_BINARY
)

func (m mode) c() C.cairo_script_mode_t {
	return C.cairo_script_mode_t(m)
}

func (m mode) String() (s string) {
	switch m {
	case ASCII:
		s = "ASCII"
	case Binary:
		s = "Binary"
	default:
		s = "unknown"
	}
	return s + " script mode"
}
