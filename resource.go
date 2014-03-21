package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

type Resource interface {
	//Originally cairo_X_acquire
	Lock() error
	//Originally cairo_X_release
	Unlock()
	//Originally cairo_X_destroy
	Close() error //BUG(jmf): do all lock/unlock resources have Destroy?
}

//BUG(jmf): document.
