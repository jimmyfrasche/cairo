package cairo

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//
//extern cairo_status_t go_write_callback(void*, unsigned char*, unsigned int);
//
//extern void go_write_callback_reaper(void*);
//
//extern cairo_status_t go_read_callback(void*, unsigned char*, unsigned int);
//
//extern void go_read_callback_reaper(void*);
//
//static cairo_status_t c_write_callback(void* w, unsigned char* data, unsigned int length) {
//	return go_write_callback(w, data, length);
//}
//
//typedef cairo_status_t (*wcallback_pass_back)(void*, unsigned char*, unsigned int);
//
///*This is required to expose the c callback wrapping the go callback back to Go*/
//static wcallback_pass_back callback_getter() {
//	return &c_write_callback;
//}
//
//static void c_write_callback_reaper(void *data) {
//	go_write_callback_reaper(data);
//}
//
//typedef void (*wreaper_pass_back)(void*);
//
//static wreaper_pass_back wreaper_getter() {
//	return &c_write_callback_reaper;
//}
import "C"

import (
	"io"
	"sync"
	"unsafe"
)

var (
	wmap = map[id]*writer{}
	mux  = new(sync.Mutex)
	//there will only ever be one writer per object.
	wkey = &C.cairo_user_data_key_t{}
)

type writer struct {
	w   io.Writer
	err error
	id  id
}

func (w *writer) write(p []byte) error {
	if w.err != nil {
		return w.err
	}
	n, err := w.w.Write(p)
	if err != nil {
		w.err = err
	}
	if n == len(p) {
		return w.err
	}
	w.err = io.ErrShortWrite
	return w.err
}

//XtensionCairoWriteFuncT is a cairo_write_func_t that expects as its closure
//argument the result of calling XtensionWrapWriter on a Writer.
//The surface or device created with this pair must be used to register
//the wrapped Writer with that objects XtensionRegisterWriter method.
//
//Anything less will cause at best memory leaks and at worst random errors.
//
//See XtensionWrapWriter for more information.
var XtensionCairoWriteFuncT = C.callback_getter()

//export go_write_callback
func go_write_callback(w unsafe.Pointer, data *C.uchar, length C.uint) C.cairo_status_t {
	W := (*writer)(w)

	bs := C.GoBytes(unsafe.Pointer(data), C.int(length))
	if err := W.write(bs); err == nil {
		return errSuccess
	}

	return errWriteError
}

//export go_write_callback_reaper
func go_write_callback_reaper(w unsafe.Pointer) {
	W := (*writer)(w)
	mux.Lock()
	defer mux.Unlock()
	delete(wmap, W.id)

	W.w = nil
	W.err = nil
}

func storeWriter(W *writer) {
	mux.Lock()
	defer mux.Unlock()
	wmap[W.id] = W
}

//XtensionRegisterWriter registers the writer wrapped by XtensionWrapWriter
//with the surface so that it does not get garbage collected until libcairo
//releases the surface.
//
//See XtensionWrapWriter for more information.
func (s *XtensionSurface) XtensionRegisterWriter(w unsafe.Pointer) {
	if err := s.Err(); err != nil {
		go_write_callback_reaper(w)
	}
	W := (*writer)(w)
	W.id = s.id()
	C.cairo_surface_set_user_data(s.s, wkey, w, C.wreaper_getter())
	storeWriter(W)
}

//XtensionRegisterWriter registers the writer wrapped by XtensionWrapWriter
//with the surface so that it does not get garbage collected until libcairo
//releases the device.
//
//See XtensionWrapWriter for more information.
func (d *XtensionDevice) XtensionRegisterWriter(w unsafe.Pointer) {
	if err := d.Err(); err != nil {
		go_write_callback_reaper(w)
	}
	W := (*writer)(w)
	W.id = d.id()
	C.cairo_device_set_user_data(d.d, wkey, w, C.wreaper_getter())
	storeWriter(W)
}

//XtensionWrapWriter wraps a writer in a special container to communicate
//with libcairo.
//
//It also stores the returned value so that it is not garbage collected.
//
//You must use this along with XtensionCairoWriteFuncT when wrapping any
//of libcairo's _create_for_stream factories.
//
//After the surface or device is created the returned pointer must
//be registered with the surface or device using its XtensionRegisterWriter
//method.
//
//Example
//
//Say you wanted to wrap an X surface created with
//cairo_X_surface_create_for_stream.
//
//In the factory for your Go surface, you need code like the following:
//	wrapped := cairo.XtensionWrapWriter(iowriter)
//	s := C.cairo_X_surface_create_for_stream(cairo.XtensionCairoWriteFuncT, wrapped)
//	S := cairo.NewXtensionSurface(s)
//	S.XtensionRegisterWriter(wrapped)
func XtensionWrapWriter(w io.Writer) (closure unsafe.Pointer) {
	W := &writer{w: w}
	return unsafe.Pointer(W)

}
