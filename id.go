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
//If s is a mapped image surface you must use XtensionRegisterAlienMappedSurface.
func XtensionRegisterAlienSurface(s *C.cairo_surface_t) (Surface, error) {
	surfaceSetID(s)
	return XtensionRevivifySurface(s)
}

//XtensionRegisterAlienMappedSurface registers a mapped image surface s.
//The surface s was mapped from is required.
//
//It is okay if from has been previously registered.
//
//This is only necessary if you are using another library that creates its own
//libcairo surface that will interact with this package.
func XtensionRegisterAlienMappedSurface(s, from *C.cairo_surface_t) (S, From Surface, err error) {
	From, err = XtensionRegisterAlienSurface(from)
	if err != nil {
		return
	}
	S, err = XtensionRegisterAlienSurface(s)
	if err != nil {
		return
	}
	S = toMapped(S.(ImageSurface))
	registerImageSurface(S, from)
	return
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
