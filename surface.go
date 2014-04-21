package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"image"
	"runtime"
	"unsafe"
)

var csurftogosurf = map[surfaceType]func(*C.cairo_surface_t) (Surface, error){
	SurfaceTypeImage:      cNewImageSurface,
	SurfaceTypeSubsurface: cNewSubsurface,
}

//XtensionRegisterRawToSurface registers a factory to convert a libcairo
//surface into a properly formed cairo Surface of the appropriate type.
//
//It is mandatory for extensions defining new surface types to call this
//function during init, otherwise users will get random not implemented
//panics for your surface.
func XtensionRegisterRawToSurface(t surfaceType, f func(*C.cairo_surface_t) (Surface, error)) {
	csurftogosurf[t] = f
}

//XtensionRevivifySurface recreates a Go Surface of the proper type
//from a C surface.
//
//This is for extension writers only.
func XtensionRevivifySurface(s *C.cairo_surface_t) (S Surface, err error) {
	t := surfaceType(C.cairo_surface_get_type(s))
	f, ok := csurftogosurf[t]
	if !ok {
		panic("No C → Go surface converter registered for " + t.String())
	}
	id := surfaceGetID(s)
	if t == SurfaceTypeImage {
		_ = id
		//TODO query mapped surface registry with id, and if so change f
	}
	S, err = f(s)
	if err != nil {
		return nil, err
	}
	if err = S.Err(); err != nil {
		return nil, err
	}
	return S, nil
}

//Surface represents an image, either as the destination of a drawing
//operation or as the source when drawing onto another surface.
//
//To draw to a Surface, create a Context with the surface as the target.
//
//All methods are documented on XtensionSurface.
//
//Originally cairo_surface_t.
type Surface interface {
	CreateSimilar(c Content, w, h int) (Surface, error)
	CreateSimilarImage(f format, w, h int) (ImageSurface, error)
	CreateSubsurface(r Rectangle) (Subsurface, error)

	Err() error
	Close() error
	Flush() error

	Content() Content
	Device() (Device, error)

	FontOptions() *FontOptions

	SetDeviceOffset(Point)
	DeviceOffset() Point

	Type() surfaceType

	HasShowTextGlyphs() bool

	MapImage(r image.Rectangle) (mappedImageSurface, error)

	Equal(Surface) bool

	//XtensionRaw is ONLY for adding libcairo subsystems outside this package.
	//Otherwise just ignore.
	XtensionRaw() *C.cairo_surface_t
	//XtensionRegisterWriter is ONLY for adding libcairo subsystems outside
	//this package.
	//Otherwise just ignore.
	XtensionRegisterWriter(unsafe.Pointer)

	id() id
}

//BUG(jmf): Surface: how to handle mime stuff?

//VectorBacked is the set of methods available to a surface
//with a native vector backend.
//
//All methods are documented on XtensionVectorSurface.
type VectorBacked interface {
	SetFallbackResolution(xppi, yppi float64)
	FallbackResolution() (xppi, yppi float64)
}

//VectorSurface is a surface with a native vector backend.
//
//All VectorBacked methods are documented on XtensionVectorSurface.
type VectorSurface interface {
	Surface
	VectorBacked
}

//Paged is the set of methods available to a paged surface.
//
//All methods are documented on XtensionPagedSurface.
type Paged interface {
	ShowPage()
	CopyPage()
}

//PagedSurface is a surface that has the concept of pages.
//
//All Paged methods are documented on XtensionPagedSurface.
type PagedSurface interface {
	Surface
	Paged
}

//PagedVectorSurface is a surface that has the concept of pages
//and a native vector backend.
type PagedVectorSurface interface {
	Surface
	Paged
	VectorBacked
}

//XtensionSurface is the "base class" for cairo surfaces.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type XtensionSurface struct {
	s *C.cairo_surface_t
}

//NewXtensionSurface creates a base go surface from a c surface.
//
//This is only for extension builders.
func NewXtensionSurface(s *C.cairo_surface_t) (x *XtensionSurface) {
	surfaceSetID(s)
	x = &XtensionSurface{s}
	runtime.SetFinalizer(x, (*XtensionSurface).Close)
	return
}

func (e *XtensionSurface) id() id {
	return surfaceGetID(e.s)
}

//Equal reports whether e and s are handles to the same surface.
func (e *XtensionSurface) Equal(s Surface) bool {
	return e.id() == s.id()
}

//XtensionRaw returns the raw cairo_surface_t pointer.
//
//XtensionRaw is only meant for creating new surface types and should NEVER
//be used directly.
func (e *XtensionSurface) XtensionRaw() *C.cairo_surface_t {
	return e.s
}

//MapImage returns an image surface that is the most efficient mechanism
//for modifying the backing store of this surface.
//
//If r is Empty, the entire surface is mapped, otherwise, just the region
//described by r is mapped.
//
//Note that r is an image.Rectangle and not a cairo.Rectangle.
//
//It is the callers responsibility to all Close on the returned surface
//in order to upload the content of the mapped image to this surface and
//destroys the image surface.
//
//The returned surface is an ImageSurface with a special Close method.
//
//Warning
//
//Using this surface as a target or source while mapped is undefined.
//
//The result of mapping a surface multiple times is undefined.
//
//Changing the device transform of either surface before the image surface
//is unmapped is undefined.
//
//Originally cairo_surface_map_to_image.
func (e *XtensionSurface) MapImage(r image.Rectangle) (mappedImageSurface, error) {
	var rect C.cairo_rectangle_int_t
	rect.x, rect.y = C.int(r.Min.X), C.int(r.Min.Y)
	rect.width, rect.height = C.int(r.Dx()), C.int(r.Dy())
	rp := &rect
	if r.Empty() {
		//use entire image
		rp = nil
	}
	return newMappedImageSurface(C.cairo_surface_map_to_image(e.s, rp), e.s)
}

//Err reports any errors on the surface.
//
//Originally cairo_surface_status.
func (e *XtensionSurface) Err() error {
	if e.s == nil {
		return ErrInvalidLibcairoHandle
	}
	return toerr_ided(C.cairo_surface_status(e.s), e)
}

//Close frees the resources used by this surface.
//
//Originally cairo_surface_destroy.
func (e *XtensionSurface) Close() error {
	if e.s == nil {
		return nil
	}
	err := e.Err()
	C.cairo_surface_destroy(e.s)
	e.s = nil
	runtime.SetFinalizer(e, nil)
	return err
}

//Flush performs any pending drawing and restores any temporary modifcations
//that libcairo has made to the surface's state.
//
//Originally cairo_surface_flush.
func (e *XtensionSurface) Flush() error {
	C.cairo_surface_flush(e.s)
	return e.Err()
}

//Type reports the type of this surface.
//
//Originally cairo_surface_get_type.
func (e *XtensionSurface) Type() surfaceType {
	return surfaceType(C.cairo_surface_get_type(e.s))
}

//Device reports the device of this surface.
//
//Originally cairo_surface_get_device.
func (e *XtensionSurface) Device() (Device, error) {
	return newCairoDevice(C.cairo_surface_get_device(e.s))
}

//Content reports the content of the surface.
//
//Originally cairo_surface_get_content.
func (e *XtensionSurface) Content() Content {
	return Content(C.cairo_surface_get_content(e.s))
}

//HasShowTextGlyphs reports whether this surface uses provided text and cluster
//data when called by a context's ShowTextGlyphs operation.
//
//Even if this method returns false, the ShowTextGlyphs operation will succeed,
//but the extra information will be ignored and the call will be equivalent
//to ShowGlyphs.
//
//Originally cairo_surface_has_show_text_glyphs.
func (e *XtensionSurface) HasShowTextGlyphs() bool {
	return C.cairo_surface_has_show_text_glyphs(e.s) == 1
}

//FontOptions retrieves the default font rendering options for this surface.
//
//Originally cairo_surface_get_font_options.
func (e *XtensionSurface) FontOptions() *FontOptions {
	f := &C.cairo_font_options_t{}
	C.cairo_surface_get_font_options(e.s, f)
	return initFontOptions(f)
}

//SetDeviceOffset sets the device offset of this surface.
//
//The device offset adds to the device coordinates determined by the coordinate
//transform matrix when drawing to a surface.
//
//One use case for this method is to create a surface that redirects a portion
//of drawing offscreen invisble to users of the Cairo api.
//Setting a transform is insufficent as queries such as DeviceToUser expose
//this offset.
//
//Note that the offset affects drawing to the surface as well
//as using the surface in a source pattern.
//
//The x and y components are in the unit of the surface's underlying device.
//
//Originally cairo_surface_set_device_offset.
func (e *XtensionSurface) SetDeviceOffset(vector Point) {
	C.cairo_surface_set_device_offset(e.s, C.double(vector.X), C.double(vector.Y))
}

//DeviceOffset reports the device offset set by SetDeviceOffset.
//
//Originally cairo_surface_get_device_offset.
func (e *XtensionSurface) DeviceOffset() (vector Point) {
	var x, y C.double
	C.cairo_surface_get_device_offset(e.s, &x, &y)
	return cPt(x, y)
}

//CreateSimilar is documented in the Surface interface.
//with e.
//For example, the new surface will have the same fallback resolution and font
//options as e.
//Generally, the new surface will also use the same backend as e, unless that
//is not possible for some reason.
//
//Initially the contents of the returned surface are all 0 (transparent if contents
//have transparency, black otherwise.)
//
//Originally cairo_surface_create_similar.
func (e *XtensionSurface) CreateSimilar(c Content, w, h int) (Surface, error) {
	s := C.cairo_surface_create_similar(e.s, c.c(), C.int(w), C.int(h))
	o := NewXtensionSurface(s)
	return o, o.Err()
}

//CreateSimilarImage creates a new surface that is as compatible as possible
//for uploading to and using in conjunction with existing surface.
//However, this surface can still be used like any normal image surface.
//
//Initially the contents of the returned surface are all 0 (transparent if contents
//have transparency, black otherwise.)
//
//Originally cairo_surface_create_similar_image.
func (e *XtensionSurface) CreateSimilarImage(f format, w, h int) (ImageSurface, error) {
	s := C.cairo_surface_create_similar_image(e.s, f.c(), C.int(w), C.int(h))
	stride := int(C.cairo_image_surface_get_stride(s))
	o := ImageSurface{
		XtensionSurface: NewXtensionSurface(s),
		format:          f,
		width:           w,
		height:          h,
		stride:          stride,
	}
	return o, o.Err()
}

//CreateSubsurface creates a window into e defined by r.
//All operations performed on s are clipped and translated to e.
//No operation on s performed outside the bounds of r are performed on e.
//
//This is useful for passing constrained child surfaces to routines that draw
//directly on the parent surface with no further allocations, double buffering,
//or copies.
//
//Warning
//
//The semantics of subsurfaces have not yet been finalized in libcairo, unless
//r is: in full device units, is contained within the extents of the target
//surface, and the target or subsurface's device transforms are not changed.
//
//Originally cairo_surface_create_for_rectangle.
func (e *XtensionSurface) CreateSubsurface(r Rectangle) (s Subsurface, err error) {
	r = r.Canon()
	x0, y0 := r.Min.c()
	x1 := C.double(r.Dx())
	y1 := C.double(r.Dy())
	ss := C.cairo_surface_create_for_rectangle(e.s, x0, y0, x1, y1)
	o := Subsurface{NewXtensionPagedVectorSurface(ss)}
	return o, o.Err()
}

func setFallbackResolution(s *C.cairo_surface_t, xppi, yppi float64) {
	C.cairo_surface_set_fallback_resolution(s, C.double(xppi), C.double(yppi))
}

func fallbackResolution(s *C.cairo_surface_t) (xppi, yppi float64) {
	var x, y C.double
	C.cairo_surface_get_fallback_resolution(s, &x, &y)
	return float64(x), float64(y)
}

//XtensionVectorSurface is the "base class" for cairo surfaces
//that have native support for vector graphics.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type XtensionVectorSurface struct {
	*XtensionSurface
}

//NewXtensionVectorSurface creates a base go vector surface from a c surface.
//
//This is only for extension builders.
func NewXtensionVectorSurface(s *C.cairo_surface_t) XtensionVectorSurface {
	return XtensionVectorSurface{NewXtensionSurface(s)}
}

//SetFallbackResolution sets the horizontal and vertical resolution
//for image fallbacks.
//
//When certain operations aren't supported natively by a backend, cairo
//will fallback by rendering operations to an image and then overlaying
//that image onto the output.
//
//If not called the default for x and y is 300 pixels per inch.
//
//Originally cairo_surface_set_fallback_resolution.
func (e XtensionVectorSurface) SetFallbackResolution(xppi, yppi float64) {
	setFallbackResolution(e.s, xppi, yppi)
}

//FallbackResolution reports the fallback resolution.
//
//Originally cairo_surface_get_fallback_resolution.
func (e XtensionVectorSurface) FallbackResolution() (xppi, yppi float64) {
	return fallbackResolution(e.s)
}

func copyPage(s *C.cairo_surface_t) {
	C.cairo_surface_copy_page(s)
}

func showPage(s *C.cairo_surface_t) {
	C.cairo_surface_show_page(s)
}

//XtensionPagedSurface is the "base class" for cairo surfaces
//that are paged.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type XtensionPagedSurface struct {
	*XtensionSurface
}

//NewXtensionPagedSurface creates a base go paged surface from a c surface.
//
//This is only for extension builders.
func NewXtensionPagedSurface(s *C.cairo_surface_t) XtensionPagedSurface {
	return XtensionPagedSurface{NewXtensionSurface(s)}
}

//CopyPage emits the current page, but does not clear it.
//The contents of the current page will be retained for the next page.
//
//Use ShowPage to emit the current page and clear it.
//
//Originally cairo_surface_copy_page.
func (e XtensionPagedSurface) CopyPage() {
	copyPage(e.s)
}

//ShowPage emits and clears the current page.
//
//Use CopyPage if you want to emit the current page but not clear it.
//
//Originally cairo_surface_show_page.
func (e XtensionPagedSurface) ShowPage() {
	showPage(e.s)
}

//XtensionPagedVectorSurface is the "base class" for cairo surfaces
//that are paged and have native support for vector graphics.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type XtensionPagedVectorSurface struct {
	*XtensionSurface
}

//NewXtensionPagedVectorSurface creates a base go paged vector surface
//from a c surface.
//
//This is only for extension builders.
func NewXtensionPagedVectorSurface(s *C.cairo_surface_t) XtensionPagedVectorSurface {
	return XtensionPagedVectorSurface{NewXtensionSurface(s)}
}

//SetFallbackResolution is documented on XtensionVectorSurface.
func (e XtensionPagedVectorSurface) SetFallbackResolution(xppi, yppi float64) {
	setFallbackResolution(e.s, xppi, yppi)
}

//FallbackResolution is documented on XtensionVectorSurface.
func (e XtensionPagedVectorSurface) FallbackResolution() (xppi, yppi float64) {
	return fallbackResolution(e.s)
}

//CopyPage is documented on XtensionPagedSurface.
func (e XtensionPagedVectorSurface) CopyPage() {
	copyPage(e.s)
}

//ShowPage is documented on XtensionPagedSurface.
func (e XtensionPagedVectorSurface) ShowPage() {
	showPage(e.s)
}

//Subsurface is a buffer or restrained copy of surface.
//
//It must be created with CreateSubsurface.
type Subsurface struct {
	XtensionPagedVectorSurface
}

func cNewSubsurface(c *C.cairo_surface_t) (Surface, error) {
	s := Subsurface{NewXtensionPagedVectorSurface(c)}
	return s, s.Err()
}
