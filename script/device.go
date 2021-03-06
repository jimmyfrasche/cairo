//Package script implements a device and surface for writing drawing operations
//to a file for debugging purposes.
//
//Libcairo must be compiled with
//	CAIRO_HAS_SCRIPT_SURFACE
//in addition to the requirements of cairo.
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
}

func cNew(d *C.cairo_device_t, m mode) (Device, error) {
	D := Device{
		XtensionDevice: cairo.NewXtensionDevice(d),
		mode:           m,
	}
	return D, D.Err()
}

//New creates a script device from writer in mode.
//
//Originally cairo_script_create_for_stream.
func New(w io.Writer, mode mode) (Device, error) {
	wp := cairo.XtensionWrapWriter(w)
	d := C.cairo_script_create_for_stream(cairo.XtensionCairoWriteFuncT, wp)
	D, err := cNew(d, mode)
	D.XtensionRegisterWriter(wp)
	return D, err
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
	return cNew(d, m)
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
