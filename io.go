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
//static cairo_user_data_key_t* new_user_key() {
//	return malloc(sizeof(cairo_user_data_key_t));
//}
//
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
//
//
//static cairo_status_t c_read_callback(void* w, unsigned char* data, unsigned int length) {
//	return go_read_callback(w, data, length);
//}
//
//typedef cairo_status_t (*rcallback_pass_back)(void*, unsigned char*, unsigned int);
//
//static rcallback_pass_back rcallback_getter() {
//	return &c_read_callback;
//}
//
//static void c_read_callback_reaper(void *data) {
//	go_read_callback_reaper(data);
//}
//
//typedef void (*rreaper_pass_back)(void*);
//
//static rreaper_pass_back rreaper_getter() {
//	return &c_read_callback_reaper;
//}
import "C"

import (
	"io"
	"sync"
	"unsafe"
)

//IOShutdowner provides a hook that cairo calls on io Readers and Writers
//that are passed through to libcairo to respond to being no longer needed.
//The error parameter is the error from the last read or write,
//which may be nil.
//
//This is entirely optional.
type IOShutdowner interface {
	IOShutdown(error)
}

var (
	wmap = map[*writer]struct{}{}
	rmap = map[*reader]struct{}{}
	mux  = new(sync.Mutex)
)

type writer struct {
	key *C.cairo_user_data_key_t
	w   io.Writer
	err error
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

type reader struct {
	key *C.cairo_user_data_key_t
	r   io.Reader
	err error
}

func (r *reader) read(p []byte) error {
	if r.err != nil {
		return r.err
	}
	_, err := r.r.Read(p)
	if err != nil {
		return nil
	}
	r.err = err
	return r.err
}

//XtensionCairoWriteFuncT is a cairo_write_func_t that expects as its closure
//argument the result of calling XtensionWrapWriter on an io.Writer.
//The surface or device created with this pair must be used to register
//the wrapped io.Wrapper with that objects XtensionRegisterWriter method.
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
	delete(wmap, W)

	if s, ok := W.w.(IOShutdowner); ok {
		s.IOShutdown(W.err)
	}

	W.w = nil
	W.err = nil
	C.free(unsafe.Pointer(W.key))
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
	C.cairo_surface_set_user_data(s.s, W.key, w, C.wreaper_getter())
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
	C.cairo_device_set_user_data(d.d, W.key, w, C.wreaper_getter())
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
	key := C.new_user_key()
	W := &writer{w: w, key: key}

	mux.Lock()
	defer mux.Unlock()
	wmap[W] = struct{}{}

	return unsafe.Pointer(W)
}

func wrapReader(r io.Reader) (closure unsafe.Pointer) {
	key := C.new_user_key()
	R := &reader{r: r, key: key}
	mux.Lock()
	defer mux.Unlock()
	rmap[R] = struct{}{}

	return unsafe.Pointer(R)
}

func (i ImageSurface) registerReader(r unsafe.Pointer) {
	if err := i.Err(); err != nil {
		go_read_callback_reaper(r)
	}
	R := (*reader)(r)
	C.cairo_surface_set_user_data(i.s, R.key, r, C.rreaper_getter())
}

//export go_read_callback
func go_read_callback(r unsafe.Pointer, data *C.uchar, length C.uint) C.cairo_status_t {
	R := (*reader)(r)

	bs := C.GoBytes(unsafe.Pointer(data), C.int(length))
	if err := R.read(bs); err == nil {
		return errSuccess
	}
	return errReadError
}

var cairoreadfunct = C.rcallback_getter()

//export go_read_callback_reaper
func go_read_callback_reaper(r unsafe.Pointer) {
	R := (*reader)(r)
	mux.Lock()
	defer mux.Unlock()
	delete(rmap, R)

	if s, ok := R.r.(IOShutdowner); ok {
		s.IOShutdown(R.err)
	}

	R.r = nil
	R.err = nil
	C.free(unsafe.Pointer(R.key))
}
