package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
//
//extern cairo_status_t go_write_callback(void*, unsigned char*, unsigned int);
//
//static cairo_status_t c_write_callback(void* w, unsigned char* data, unsigned int length) {
//	return go_write_callback(w, data, length);
//}
//
//typedef cairo_status_t (*callback_pass_back)(void*, unsigned char*, unsigned int);
//
///*This is required to expose the c callback wrapping the go callback back to Go*/
//static callback_pass_back callback_getter() {
//	return &c_write_callback;
//}
import "C"

import (
	"io"
	"unsafe"
)

//XtensionCairoWriteFuncT is a cairo_write_func_t that expects an *io.Writer
//as its closure argument.
//
//It is the caller's responsibility to hold a reference to the io.Writer so that it does
//not become garbage collected.
//
//See cairo/svg for an example.
var XtensionCairoWriteFuncT = C.callback_getter()

//export go_write_callback
func go_write_callback(writer unsafe.Pointer, data *C.uchar, length C.uint) C.cairo_status_t {
	w := *(*io.Writer)(writer)
	len := int(length)

	bs := C.GoBytes(unsafe.Pointer(data), C.int(len))
	n, err := w.Write(bs)
	if err != nil || n != len {
		return errWriteError
	}

	return errSuccess
}
