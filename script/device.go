//Package script implements a device and surface for writing drawing operations
//to a file for debugging purposes.
package script

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//#include <cairo/cairo-script.h>
import "C"

import (
	"io"
	"unsafe"
	"github.com/jimmyfrasche/cairo"
	"github.com/jimmyfrasche/cairo/recording"
)

//Device is a pseudo-Device that records operations performed on it.
type Device struct {
	*cairo.XtensionDevice
	mode mode
	//w is used to ensure a reference to the writer lives as long as we do
	w io.Writer
}

func cNew(d *C.cairo_device_t, m mode, w io.Writer) (Device, error) {
	D := Device{
		XtensionDevice: cairo.NewXtensionDevice(d),
		mode:           m,
		w:              w,
	}
	return D, D.Err()
}

//New creates a script device from writer in mode.
//
//Warning
//
//It is the caller's responsibility to keep a reference to w for the lifetime
//of this surface.
//As it is passed to libcairo, the Go garbage collector will otherwise find
//no reference to it.
//
//Originally cairo_script_create_for_stream.
func New(w io.Writer, mode mode) (Device, error) {
	wp := unsafe.Pointer(&w)
	d := C.cairo_script_create_for_stream(cairo.XtensionCairoWriteFuncT, wp)
	return cNew(d, mode, w)
}

//NewFile creates a device, in mode, that writes to filename.
//
//Originally cairo_script_create.
func NewFile(filename string, mode mode) (Device, error) {
	s := C.CString(filename)
	d := C.cairo_script_create(s)
	C.free(unsafe.Pointer(s))
	return cNew(d, mode, nil)
}

//FromRecordingSurface outputs the record operations in rs to d.
//
//Originally cairo_script_from_recording_surface.
func (d Device) FromRecordingSurface(rs recording.Surface) error {
	//we grab the status from d.Err
	_ = C.cairo_script_from_recording_surface(d.XtensionRaw(), rs.XtensionRaw())
	return d.Err()
}

func reviv(d *C.cairo_device_t) (cairo.Device, error) {
	m := mode(C.cairo_script_get_mode(d))
	return cNew(d, m, nil)
}

func init() {
	cairo.XtensionRegisterRawToDevice(cairo.DeviceTypeScript, reviv)
}

//Comment adds a comment to the script.
//
//Originally cairo_script_write_comment.
func (d Device) Comment(c string) {
	s := C.CString(c)
	C.cairo_script_write_comment(d.XtensionRaw(), s, C.int(len(c)))
	C.free(unsafe.Pointer(s))
}

//Mode reports the recording mode of this device.
//
//Originally cairo_script_get_mode.
func (d Device) Mode() mode {
	return d.mode
}

//Proxy creates a script surface p that renders to s and records to d.
//
//Originally cairo_script_surface_create_for_target.
func (d Device) Proxy(s cairo.Surface) (p Surface, err error) {
	sr := C.cairo_script_surface_create_for_target(d.XtensionRaw(), s.XtensionRaw())
	return cNewSurf(sr)
}
