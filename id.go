package cairo

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <stdint.h>
//#include <cairo/cairo.h>
//
//static void* cgo_id_malloc(uint64_t id) {
//	uint64_t* p = (uint64_t*) malloc(sizeof(uint64_t));
//	*p = id;
//	return p;
//}
//
//static void gocairo_free(void* data) {
//	free(data);
//}
//
//static cairo_destroy_func_t gocairo_free_get() {
//	return &gocairo_free;
//}
import "C"

import (
	"sync"
	"unsafe"
)

var free = C.gocairo_free_get()

type (
	id        uint64
	subtypeID uint64
)

var (
	nextID uint64
	idmux  = &sync.Mutex{}
	idkey  = &C.cairo_user_data_key_t{}

	subtypes = map[string]subtypeID{}
)

func generateID() unsafe.Pointer {
	idmux.Lock()
	defer idmux.Unlock()

	p := C.cgo_id_malloc(C.uint64_t(nextID))
	nextID++
	return unsafe.Pointer(p)
}

func ctoint(p unsafe.Pointer) uint64 {
	return uint64(*(*C.uint64_t)(p))
}

func surfaceSetID(s *C.cairo_surface_t) {
	if C.cairo_surface_get_user_data(s, idkey) != nil {
		return
	}
	C.cairo_surface_set_user_data(s, idkey, generateID(), free)
}

func surfaceGetID(s *C.cairo_surface_t) id {
	p := C.cairo_surface_get_user_data(s, idkey)
	if p == nil {
		panic("surface does not have ID - created outside of cairo binding and not registered")
	}
	return id(ctoint(p))
}

//XtensionRegisterAlienSurface registers a surface created outside of this
//package and creates a Surface of the proper type.
//
//This is only necessary if you are using another library that creates its own
//libcairo surface that will interact with this package.
//
//There is no way to tell if an image surface is mapped or not, so if s is
//an image surface you must specify whether it is.
func XtensionRegisterAlienSurface(s *C.cairo_surface_t, mapped bool) (Surface, error) {
	surfaceSetID(s)
	if surfaceType(C.cairo_surface_get_type(s)) == SurfaceTypeImage && mapped {
		//surfaceGetID(s)
		//TODO add to mapped registry
	}
	return XtensionRevivifySurface(s)
}

func deviceSetID(d *C.cairo_device_t) {
	if C.cairo_device_get_user_data(d, idkey) != nil {
		return
	}
	C.cairo_device_set_user_data(d, idkey, generateID(), free)
}

func deviceGetID(d *C.cairo_device_t) id {
	p := C.cairo_device_get_user_data(d, idkey)
	if p == nil {
		panic("device does not have ID - created outside of cairo binding and not registered")
	}
	return id(ctoint(p))
}

//XtensionRegisterAlienDevice registers a device created outside of this
//package and creates a Device of the proper type.
//
//This is only necessary if you are using another library that creates its own
//libcairo device that will interact with this package.
func XtensionRegisterAlienDevice(d *C.cairo_device_t) (Device, error) {
	deviceSetID(d)
	return newCairoDevice(d)
}

func createSubtypeID(name string) subtypeID {
	idmux.Lock()
	defer idmux.Unlock()

	if _, ok := subtypes[name]; ok {
		panic("subtype " + name + " previously registered")
	}

	id := subtypeID(nextID)
	subtypes[name] = id
	nextID++
	return id
}
