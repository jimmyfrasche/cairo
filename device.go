package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"runtime"
)

type Device interface {
	Type() deviceType
	Err() error
	Close() error
	Lock() error
	Unlock()
}

type cairoDevice struct {
	d *C.cairo_device_t
}

func newCairoDevice(d *C.cairo_device_t) Device {
	if d == nil {
		return nil
	}
	o := &cairoDevice{d: d}
	runtime.SetFinalizer(o, (*cairoDevice).Close)
	return o
}

func (c *cairoDevice) Lock() error {
	return toerr(C.cairo_device_acquire(c.d))
}

func (c *cairoDevice) Unlock() {
	C.cairo_device_release(c.d)
}

func (c *cairoDevice) Close() error {
	if c == nil || c.d == nil {
		return nil
	}
	err := c.Err()
	C.cairo_device_destroy(c.d)
	c.d = nil
	runtime.SetFinalizer(c, nil)
	return err
}

func (c *cairoDevice) Err() error {
	return toerr(C.cairo_device_status(c.d))
}

func (c *cairoDevice) Type() deviceType {
	return deviceType(C.cairo_device_get_type(c.d))
}
