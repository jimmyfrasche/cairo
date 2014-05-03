package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"runtime"
	"unsafe"
)

var cdevtogodev = map[deviceType]func(*C.cairo_device_t) (Device, error){}

//XtensionRegisterRawToDevice registers a factory to convert a libcairo
//device into a properly formed Device of the proper underlying type.
//
//If no factory is registered a default XtensionDevice will be returned.
func XtensionRegisterRawToDevice(d deviceType, f func(*C.cairo_device_t) (Device, error)) {
	cdevtogodev[d] = f
}

//A Device abstracts the rendering backend of a cairo surface.
//
//Devices are created using custom functions specific to the rendering system
//you want to use.
//
//An important function that devices fulfill is sharing access to the rendering
//system between libcairo and your application.
//If you want to access a device directly that you used to draw to with Cairo,
//you must first call Flush to ensure that Cairo finishes all operations on the
//device and resets it to a clean state.
//
//Cairo also provides the Lock and Release methods to synchronize access
//to the rendering system in a multithreaded environment.
//This is done internally, but can also be used by applications.
//
//Putting this all together a function that works with devices should often look like:
//	func WithDevice(d Device, work func(Device) error) (err error) {
//		d.Flush()
//		if err = d.Lock(); err != nil {
//			return
//		}
//		defer d.Unlock()
//		return work(d)
//	}
//
//All methods are documented on XtensionDevice.
type Device interface {
	Type() deviceType
	Err() error
	Close() error
	Lock() error
	Unlock()
	Flush()
	Equal(Device) bool

	//only for writing extensions.
	XtensionRaw() *C.cairo_device_t
	XtensionRegisterWriter(w unsafe.Pointer)
	id() id
}

//XtensionDevice is the "base class" and default implementation for libcairo
//devices.
//
//Unless a particular type of device exposes special operations on the device,
//it will be an object of this type regardless of its deviceType.
//
//Originally cairo_device_t.
type XtensionDevice struct {
	d *C.cairo_device_t
}

func newCairoDevice(d *C.cairo_device_t) (Device, error) {
	t := deviceType(C.cairo_device_get_type(d))
	_ = deviceGetID(d) //panics if created outside of cairo without being registered
	f, ok := cdevtogodev[t]
	if !ok {
		D := NewXtensionDevice(d)
		return D, D.Err()
	}
	return f(d)
}

//NewXtensionDevice creates a plain device for a c surface.
//
//This is only for extension builders.
func NewXtensionDevice(d *C.cairo_device_t) *XtensionDevice {
	deviceSetID(d)
	o := &XtensionDevice{d: d}
	runtime.SetFinalizer(o, (*XtensionDevice).Close)
	return o
}

func (c *XtensionDevice) id() id {
	return deviceGetID(c.d)
}

//Equal reports whether c and d are handles of the same device.
func (c *XtensionDevice) Equal(d Device) bool {
	return c.id() == d.id()
}

//Lock acquires the device for the current thread.
//This method will block until no other thread has acquired the device.
//
//If err is nil, you successfully acquired the device.
//From now on your thread owns the device and no other thread will be able
//to acquire it until a matching call to Unlock.
//
//It is allowed to recursively acquire the device multiple times
//from the same thread.
//
//Note
//
//You must never acquire two different devices at the same time unless
//this is explicitly allowed.
//Otherwise, the possibility of deadlocks exist.
//
//As various libcairo functions can acquire devices when called,
//these may also cause deadlocks when you call them with an acquired device.
//So you must not have a device acquired when calling them.
//
//These functions are marked in the documentation.
//
//Orignally cairo_device_acquire.
func (c *XtensionDevice) Lock() (err error) {
	return toerr(C.cairo_device_acquire(c.d))
}

//Unlock releases the device previously acquired by Lock.
//
//Originally cairo_device_release.
func (c *XtensionDevice) Unlock() {
	C.cairo_device_release(c.d)
}

//Close releases the resources of this device.
//
//Originally cairo_device_destroy.
func (c *XtensionDevice) Close() error {
	if c == nil || c.d == nil {
		return nil
	}
	err := c.Err()
	C.cairo_device_destroy(c.d)
	c.d = nil
	runtime.SetFinalizer(c, nil)
	return err
}

//Err reports any error on this device.
//
//Originally cairo_device_status.
func (c *XtensionDevice) Err() error {
	if c.d == nil {
		return ErrInvalidLibcairoHandle
	}
	return toerrIded(C.cairo_device_status(c.d), c)
}

//Type reports the type of this device.
//
//Originally cairo_device_get_type.
func (c *XtensionDevice) Type() deviceType {
	return deviceType(C.cairo_device_get_type(c.d))
}

//Flush any pending operations and restore any temporary modification to the
//device state made by libcairo.
//
//This method must be called before switching from using the device with Cairo
//to operating on it directly with native APIs.
//If the device doesn't support direct access, then this does nothing.
//
//This may lock devices.
//
//Originally cairo_device_flush.
func (c *XtensionDevice) Flush() {
	C.cairo_device_flush(c.d)
}

//XtensionRaw returns the raw cairo_device_t pointer.
//
//XtensionRaw is only meant for creating new device types and should NEVER
//be used directly.
func (c *XtensionDevice) XtensionRaw() *C.cairo_device_t {
	return c.d
}
