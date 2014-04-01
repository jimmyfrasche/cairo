package cairo

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
import "C"

import (
	"runtime"
	"unsafe"
)

type ImageSurface struct {
	*XtensionSurface
	format                format
	width, height, stride int
}

func newImg(s *C.cairo_surface_t, format format, width, height, stride int) (ImageSurface, error) {
	S := ImageSurface{
		XtensionSurface: NewXtensionSurface(s),
		format:          format,
		width:           width,
		height:          height,
		stride:          stride,
	}
	return S, S.Err()
}

//Originally cairo_image_surface_create.
func NewImageSurface(format format, width, height int) (ImageSurface, error) {
	is := C.cairo_image_surface_create(format.c(), C.int(width), C.int(height))
	stride := int(C.cairo_image_surface_get_stride(is))
	return newImg(is, format, width, height, stride)
}

func cNewImageSurface(s *C.cairo_surface_t) (Surface, error) {
	format := format(C.cairo_image_surface_get_format(s))
	width := int(C.cairo_image_surface_get_width(s))
	height := int(C.cairo_image_surface_get_height(s))
	stride := int(C.cairo_image_surface_get_stride(s))

	return newImg(s, format, width, height, stride)
}

//NewImageSurfaceFromPNG creates a new image surface and initalizes
//it with the given PNG file.
//
//Originally cairo_image_surface_create_from_png_stream.
func NewImageSurfaceFromPNG(r Reader) (ImageSurface, error) {
	rp := wrapReader(r)
	is := C.cairo_image_surface_create_from_png_stream(cairoreadfunct, rp)
	s, err := cNewImageSurface(is)
	S := s.(ImageSurface)
	if err != nil {
		return S, err
	}
	S.registerReader(rp)
	return S, nil
}

//NewImageSurfaceFromPNGFile creates a new image surface and initializes
//it with the contents of the given PNG file.
//
//Originally cairo_image_surface_create_from_png.
func NewImageSurfaceFromPNGFile(filename string) (ImageSurface, error) {
	f := C.CString(filename)
	is := C.cairo_image_surface_create_from_png(f)
	C.free(unsafe.Pointer(f))
	s, err := cNewImageSurface(is)
	return s.(ImageSurface), err
}

//BUG(jmf): ImageSurface: need safe wrapper around get_data

//BUG(jmf): ImageSurface: need image_surface_create_for_data analog(s)

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

type mappedImageSurface struct {
	ImageSurface
	from *C.cairo_surface_t
}

func newMappedImageSurface(s, from *C.cairo_surface_t) (m mappedImageSurface, err error) {
	im, err1 := cNewImageSurface(s)
	if err1 != nil {
		err = err1
		return
	}
	//Clear default finalizer so GC doesn't call surface_destroy.
	//We do not set a new finalizer on mappedImageSurface.Close,
	//because that would not do the right thing and the user is expected to close
	//manually when done.
	runtime.SetFinalizer(m.XtensionSurface, nil)
	m = mappedImageSurface{
		ImageSurface: im.(ImageSurface),
		from:         from,
	}
	err = m.Err()
	if err != nil {
		m.s = nil
	}
	return
}

func (m mappedImageSurface) Close() error {
	err := m.Err()
	C.cairo_surface_unmap_image(m.from, m.s)
	m.from = nil
	m.s = nil
	return err
}
