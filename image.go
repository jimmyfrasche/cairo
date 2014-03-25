package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

type ImageSurface struct {
	ExtensionSurface
	format                format
	width, height, stride int
}

//Originally cairo_image_surface_create.
func NewImageSurface(format format, width, height int) (ImageSurface, error) {
	is := C.cairo_image_surface_create(format.c(), C.int(width), C.int(height))
	stride := int(C.cairo_image_surface_get_stride(is))
	s := ImageSurface{
		ExtensionSurface: ExtensionNewSurface(is),
		format:           format,
		width:            width,
		height:           height,
		stride:           stride,
	}
	return s, s.Err()
}

func cNewImageSurface(s *C.cairo_surface_t) (Surface, error) {
	format := format(C.cairo_image_surface_get_format(s))
	width := int(C.cairo_image_surface_get_width(s))
	height := int(C.cairo_image_surface_get_height(s))
	stride := int(C.cairo_image_surface_get_stride(s))
	S := ImageSurface{
		ExtensionSurface: ExtensionNewSurface(s),
		format:           format,
		width:            width,
		height:           height,
		stride:           stride,
	}

	return S, S.Err()
}

//BUG(jmf): need safe wrapper around get_data

//BUG(jmf): need image_surface_create_for_data analog(s)

//Format reports the format of the surface.
//
//Originally cairo_image_surface_get_format.
func (i ImageSurface) Format() format {
	return i.format
}

//Width reports the width of the surface in pixels.
//
//Originally cairo_image_surface_get_width.
func (i ImageSurface) Width() int {
	return i.width
}

//Height reports the height of the surface in pixels.
//
//Originally cairo_image_surface_get_height.
func (i ImageSurface) Height() int {
	return i.height
}

//Stride reports the stride of the image surface in number of bytes.
//
//Originally cairo_image_surface_get_stride.
func (i ImageSurface) Stride() int {
	return i.stride
}
