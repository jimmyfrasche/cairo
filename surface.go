package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

//Surface represents an image, either as the destination of a drawing
//operation or as the source when drawing onto another surface.
//
//To draw to a Surface, create a Context with the surface as the target.
//
//All methods are documented on ExtensionSurface.
type Surface interface {
	CreateSimilar(c content, w, h int) (Surface, error)
	CreateSimilarImage(f format, w, h int) (ImageSurface, error)
	CreateSubsurface(r Rectangle) (Surface, error)

	Err() error
	Close() error
	Flush() error

	Content() content
	Device() Device

	FontOptions() *FontOptions

	SetDeviceOffset(Point)
	DeviceOffset() Point

	Type() surfaceType

	HasShowTextGlyphs() bool

	//ExtensionRaw is ONLY for adding libcairo subsystems outside this package.
	//Otherwise just ignore.
	ExtensionRaw() *C.cairo_surface_t
}

//BUG(jmf): Surface: how to handle map/unmap

//BUG(jmf): Surface: how to handle mime stuff?

//VectorBacked is the set of methods available to a surface
//with a native vector backend.
//
//All methods are documented on ExtensionVectorSurface.
type VectorBacked interface {
	SetFallbackResolution(xppi, yppi float64)
	FallbackResolution() (xppi, yppi float64)
}

//VectorSurface is a surface with a native vector backend.
//
//All VectorBacked methods are documented on ExtensionVectorSurface.
type VectorSurface interface {
	Surface
	VectorBacked
}

//Paged is the set of methods available to a paged surface.
//
//All methods are documented on ExtensionPagedSurface.
type Paged interface {
	ShowPage()
	CopyPage()
}

//PagedSurface is a surface that has the concept of pages.
//
//All Paged methods are documented on ExtensionPagedSurface.
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

//ExtensionSurface is the "base class" for cairo surfaces.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type ExtensionSurface struct {
	s *C.cairo_surface_t
}

//ExtensionNewSurface creates a base go surface from a c surface.
//
//This is only for extension builders.
func ExtensionNewSurface(s *C.cairo_surface_t) ExtensionSurface {
	return ExtensionSurface{s}
}

//ExtensionRaw returns the raw cairo_surface_t pointer.
//
//ExtensionRaw is only meant for creating new surface types and should NEVER
//be used directly.
func (e ExtensionSurface) ExtensionRaw() *C.cairo_surface_t {
	return e.s
}

//Err reports any errors on the surface.
//
//Originally cairo_surface_status.
func (e ExtensionSurface) Err() error {
	return toerr(C.cairo_surface_status(e.s))
}

//Close frees the resources used by this surface.
//
//Originally cairo_surface_destroy.
func (e ExtensionSurface) Close() error {
	if e.s == nil {
		return nil
	}
	C.cairo_surface_destroy(e.s)
	e.s = nil
	return nil
}

//Flush performs any pending drawing and restores any temporary modifcations
//that libcairo has made to the surface's state.
//
//Originally cairo_surface_flush.
func (e ExtensionSurface) Flush() error {
	C.cairo_surface_flush(e.s)
	return e.Err()
}

//Type reports the type of this surface.
//
//Originally cairo_surface_get_type.
func (e ExtensionSurface) Type() surfaceType {
	return surfaceType(C.cairo_surface_get_type(e.s))
}

//Device reports the device of this surface.
//
//Originally cairo_surface_get_device.
func (e ExtensionSurface) Device() Device {
	return newCairoDevice(C.cairo_surface_get_device(e.s))
}

//Content reports the content of the surface.
//
//Originally cairo_surface_get_content.
func (e ExtensionSurface) Content() content {
	return content(C.cairo_surface_get_content(e.s))
}

//HasShowTextGlyphs reports whether this surface uses provided text and cluster
//data when called by a context's ShowTextGlyphs operation.
//
//Even if this method returns false, the ShowTextGlyphs operation will succeed,
//but the extra information will be ignored and the call will be equivalent
//to ShowGlyphs.
//
//Originally cairo_surface_has_show_text_glyphs.
func (e ExtensionSurface) HasShowTextGlyphs() bool {
	return C.cairo_surface_has_show_text_glyphs(e.s) == 1
}

//FontOptions retrieves the default font rendering options for this surface.
//
//Originally cairo_surface_get_font_options.
func (e ExtensionSurface) FontOptions() *FontOptions {
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
func (e ExtensionSurface) SetDeviceOffset(vector Point) {
	C.cairo_surface_set_device_offset(e.s, C.double(vector.X), C.double(vector.Y))
}

//DeviceOffset reports the device offset set by SetDeviceOffset.
//
//Originally cairo_surface_get_device_offset.
func (e ExtensionSurface) DeviceOffset() (vector Point) {
	var x, y C.double
	C.cairo_surface_get_device_offset(e.s, &x, &y)
	return Pt(float64(x), float64(y))
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
func (e ExtensionSurface) CreateSimilar(c content, w, h int) (Surface, error) {
	s := C.cairo_surface_create_similar(e.s, c.c(), C.int(w), C.int(h))
	o := ExtensionSurface{s}
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
func (e ExtensionSurface) CreateSimilarImage(f format, w, h int) (ImageSurface, error) {
	s := C.cairo_surface_create_similar_image(e.s, f.c(), C.int(w), C.int(h))
	stride := int(C.cairo_image_surface_get_stride(s))
	o := ImageSurface{
		ExtensionSurface: ExtensionSurface{
			s,
		},
		format: f,
		width:  w,
		height: h,
		stride: stride,
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
func (e ExtensionSurface) CreateSubsurface(r Rectangle) (s Surface, err error) {
	x0, y0, x1, y1 := r.Canon().c()
	ss := C.cairo_surface_create_for_rectangle(e.s, x0, y0, x1, y1)
	o := ExtensionSurface{ss}
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

//ExtensionVectorSurface is the "base class" for cairo surfaces
//that have native support for vector graphics.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type ExtensionVectorSurface struct {
	ExtensionSurface
}

//ExtensionNewVectorSurface creates a base go vector surface from a c surface.
//
//This is only for extension builders.
func ExtensionNewVectorSurface(s *C.cairo_surface_t) ExtensionVectorSurface {
	return ExtensionVectorSurface{ExtensionNewSurface(s)}
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
func (e ExtensionVectorSurface) SetFallbackResolution(xppi, yppi float64) {
	setFallbackResolution(e.s, xppi, yppi)
}

//FallbackResolution reports the fallback resolution.
//
//Originally cairo_surface_get_fallback_resolution.
func (e ExtensionVectorSurface) FallbackResolution() (xppi, yppi float64) {
	return fallbackResolution(e.s)
}

func copyPage(s *C.cairo_surface_t) {
	C.cairo_surface_copy_page(s)
}

func showPage(s *C.cairo_surface_t) {
	C.cairo_surface_show_page(s)
}

//ExtensionPagedSurface is the "base class" for cairo surfaces
//that are paged.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type ExtensionPagedSurface struct {
	ExtensionSurface
}

//ExtensionNewPagedSurface creates a base go paged surface from a c surface.
//
//This is only for extension builders.
func ExtensionNewPagedSurface(s *C.cairo_surface_t) ExtensionPagedSurface {
	return ExtensionPagedSurface{ExtensionNewSurface(s)}
}

//CopyPage emits the current page, but does not clear it.
//The contents of the current page will be retained for the next page.
//
//Use ShowPage to emit the current page and clear it.
//
//Originally cairo_surface_copy_page.
func (e ExtensionPagedSurface) CopyPage() {
	copyPage(e.s)
}

//ShowPage emits and clears the current page.
//
//Use CopyPage if you want to emit the current page but not clear it.
//
//Originally cairo_surface_show_page.
func (e ExtensionPagedSurface) ShowPage() {
	showPage(e.s)
}

//ExtensionPagedVectorSurface is the "base class" for cairo surfaces
//that are paged and have native support for vector graphics.
//
//It is meant only for embedding in new surface types and should NEVER
//be used directly.
type ExtensionPagedVectorSurface struct {
	ExtensionSurface
}

//ExtensionNewPagedVectorSurface creates a base go paged vector surface
//from a c surface.
//
//This is only for extension builders.
func ExtensionNewPagedVectorSurface(s *C.cairo_surface_t) ExtensionPagedVectorSurface {
	return ExtensionPagedVectorSurface{ExtensionNewSurface(s)}
}

//SetFallbackResolution is documented on ExtensionVectorSurface.
func (e ExtensionPagedVectorSurface) SetFallbackResolution(xppi, yppi float64) {
	setFallbackResolution(e.s, xppi, yppi)
}

//FallbackResolution is documented on ExtensionVectorSurface.
func (e ExtensionPagedVectorSurface) FallbackResolution() (xppi, yppi float64) {
	return fallbackResolution(e.s)
}

//CopyPage is documented on ExtensionPagedSurface.
func (e ExtensionPagedVectorSurface) CopyPage() {
	copyPage(e.s)
}

//ShowPage is documented on ExtensionPagedSurface.
func (e ExtensionPagedVectorSurface) ShowPage() {
	showPage(e.s)
}