package mimepattern

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <string.h>
//#include <cairo/cairo.h>
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
	"bytes"
	"errors"
	"image"
	"io"
	"net/url"
	"unsafe"

	"github.com/jimmyfrasche/cairo"
)

//URL creates a Pattern similar to New, except that, instead of
//storing the uncompressed image data, a url to the image is stored.
//
//This is only useful on an SVG surface.
//
//Originally cairo_surface_set_mime_data.
func URL(URL *url.URL, r io.Reader) (cairo.Pattern, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	is, err := cairo.FromImage(img)
	if err != nil {
		return nil, err
	}
	defer is.Close()

	bs := []byte(URL.String())
	if err = embed(is, uri, bs); err != nil {
		return nil, err
	}

	return newPattern(is)
}

//New creates a pattern of the image data in r of type mime.
//
//The returned pattern is essentially a SurfacePattern but the underlying
//surface is irrecoverable as any modification will lead to corruption.
//
//New will fail if there is not an image decoder for the relevant mime type
//installed with the image package.
//See the image package documentation for more information.
//It is up to you to ensure that the image data in r is of the mime type
//specified by mime.
//
//Originally cairo_surface_set_mime_data.
func New(mime mime, r io.Reader) (cairo.Pattern, error) {
	if !mime.valid() {
		return nil, errors.New("invalid or unsupported mime type: " + mime.s())
	}

	buf := bytes.NewBuffer(make([]byte, 0, 4096))
	img, _, err := image.Decode(io.TeeReader(r, buf))
	if err != nil {
		return nil, err
	}

	is, err := cairo.FromImage(img)
	if err != nil {
		return nil, err
	}
	defer is.Close()

	if err = embed(is, mime, buf.Bytes()); err != nil {
		return nil, err
	}

	return newPattern(is)
}

type mimepattern struct {
	*cairo.XtensionPattern
}

//newPattern creates a SurfacePattern without a Surface method.
func newPattern(is cairo.ImageSurface) (cairo.Pattern, error) {
	p, err := cairo.NewSurfacePattern(is)
	if err != nil {
		return nil, err
	}

	return mimepattern{p.XtensionPattern}, nil
}

func embed(s cairo.Surface, mime mime, bs []byte) error {
	raw := s.XtensionRaw()
	len := C.size_t(uintptr(len(bs)))

	//We make this copy as bs will be GC'd.
	//We could track the object lifetime and cache the []byte in this package
	//until it is destroyed, which is a fine optimization but necessitates a lot
	//of code.
	data := (*C.uchar)(C.malloc(len))
	C.memcpy(unsafe.Pointer(data), unsafe.Pointer(&bs[0]), len)

	st := C.cairo_surface_set_mime_data(raw, mime.c(), data, C.ulong(len), C.gocairo_free_get(), nil)
	if st != C.CAIRO_STATUS_SUCCESS {
		return errors.New("could not set mime data")
	}

	return nil
}
