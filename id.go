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

func (s subtypeID) c() unsafe.Pointer {
	return C.cgo_id_malloc(C.uint64_t(s))
}

var (
	nextID uint64
	idmux  = &sync.Mutex{}
	idkey  = &C.cairo_user_data_key_t{}
	stkey  = &C.cairo_user_data_key_t{}

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

func fontSetSubtypeID(f *C.cairo_font_face_t, s subtypeID) {
	if fontType(C.cairo_font_face_get_type(f)) != FontTypeUser {
		panic("font is not a user font")
	}
	if C.cairo_font_face_get_user_data(f, stkey) != nil {
		panic("font already has subtype set")
	}
	C.cairo_font_face_set_user_data(f, stkey, s.c(), free)
}

func fontGetSubtypeID(f *C.cairo_font_face_t) subtypeID {
	if fontType(C.cairo_font_face_get_type(f)) != FontTypeUser {
		panic("font is not a user font")
	}
	p := C.cairo_font_face_get_user_data(f, stkey)
	if p == nil {
		panic("no subtype set: font not registered")
	}
	return subtypeID(ctoint(p))
}

func patternSetSubtypeID(p *C.cairo_pattern_t, s subtypeID) {
	if fontType(C.cairo_pattern_get_type(p)) != PatternTypeRasterSource {
		panic("pattern is not a raster pattern")
	}
	if C.cairo_pattern_get_user_data(p, stkey) != nil {
		panic("pattern already has subtype set")
	}
	C.cairo_pattern_set_user_data(p, stkey, s.c(), free)
}

func patternGetSubtypeID(p *C.cairo_pattern_t) subtypeID {
	if fontType(C.cairo_pattern_get_type(p)) != PatternTypeRasterSource {
		panic("pattern is not a raster pattern")
	}
	ptr := C.cairo_pattern_get_user_data(p, stkey)
	if ptr == nil {
		panic("no subtype set: pattern not registered")
	}
	return subtypeID(ctoint(ptr))
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

type fontsubfac struct {
	name string
	fac  func(*C.cairo_font_face_t) (Font, error)
}

var (
	fontsubtypenames = map[string]subtypeID{}
	fontsubtypes     = map[subtypeID]*fontsubfac{}
)

//XtensionRegisterAlienUserFontSubtype registers a factory to create a Go wrapper
//around an existing libcairo user font and associates it with a unique name,
//retrievable via Subtype.
//
//After the subtype is registered, instances MUST be registered with
//XtensionRegisterAlienUserFont.
func XtensionRegisterAlienUserFontSubtype(name string, fac func(*C.cairo_font_face_t) (Font, error)) {
	if name == "" {
		panic("user font subtype name must not be empty string")
	}
	if fac == nil {
		panic("user font factory must not be nil")
	}
	idmux.Lock()
	defer idmux.Unlock()

	if _, ok := fontsubtypenames[name]; ok {
		panic("subtype " + name + " previously registered")
	}

	id := subtypeID(nextID)
	fontsubtypenames[name] = id
	fontsubtypes[id] = &fontsubfac{
		name: name,
		fac:  fac,
	}
	nextID++
}

//XtensionRegisterAlienUserFont registers a libcairo user font with cairo.
//
//The subtype must be registered XtensionRegisterAlienUserFontSubtype.
func XtensionRegisterAlienUserFont(subtype string, f *C.cairo_font_face_t) (Font, error) {
	id := fontsubtypenames[subtype]
	fontSetSubtypeID(f, id)
	return fontsubtypes[id].fac(f)
}

type rastersubfac struct {
	name string
	fac  func(*C.cairo_pattern_t) (Pattern, error)
}

var (
	rastersubtypenames = map[string]subtypeID{}
	rastersubtypes     = map[subtypeID]*rastersubfac{}
)

//XtensionRegisterAlienRasterPatternSubtype registers a factory to create
//a Go wrapper an existing libcairo raster pattern and associates it with
//a unique name, retrievable via Subtype.
//
//After the subtype is registered, instances MUST be registered with
//XtensionRegisterAlienRasterPattern.
func XtensionRegisterAlienRasterPatternSubtype(name string, fac func(*C.cairo_pattern_t) (Pattern, error)) {
	if name == "" {
		panic("raster pattern subtype name must not be empty string")
	}
	if fac == nil {
		panic("raster pattern factory must not be nil")
	}
	idmux.Lock()
	defer idmux.Unlock()

	if _, ok := rastersubtypenames[name]; ok {
		panic("subtype " + name + "previously registered")
	}

	id := subtypeID(nextID)
	rastersubtypenames[name] = id
	rastersubtypes[id] = &rastersubfac{
		name: name,
		fac:  fac,
	}
	nextID++
}

//XtensionRegisterAlienRasterPattern registers a libcairo user font with cairo.
//
//The subtype must be registered
//with XtensionRegisterAlienRasterPatternSubtype.
func XtensionRegisterAlienRasterPattern(subtype string, p *C.cairo_pattern_t) (Pattern, error) {
	id := rastersubtypenames[subtype]
	patternSetSubtypeID(p, id)
	return rastersubtypes[id].fac(p)
}
