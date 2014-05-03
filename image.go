package cairo

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
import "C"

import (
	"image"
	"runtime"
	"sync"
	"unsafe"
)

//An ImageSurface is an in-memory surface.
type ImageSurface struct {
	*XtensionSurface
	format                Format
	width, height, stride int
}

func newImg(s *C.cairo_surface_t, format Format, width, height, stride int) (ImageSurface, error) {
	S := ImageSurface{
		XtensionSurface: NewXtensionSurface(s),
		format:          format,
		width:           width,
		height:          height,
		stride:          stride,
	}
	return S, S.Err()
}

//NewImageSurface creates an image surface of the given width, height,
//and format.
//
//Originally cairo_image_surface_create.
func NewImageSurface(format Format, width, height int) (ImageSurface, error) {
	is := C.cairo_image_surface_create(format.c(), C.int(width), C.int(height))
	stride := int(C.cairo_image_surface_get_stride(is))
	return newImg(is, format, width, height, stride)
}

func cNewImageSurface(s *C.cairo_surface_t) (Surface, error) {
	format := Format(C.cairo_image_surface_get_format(s))
	width := int(C.cairo_image_surface_get_width(s))
	height := int(C.cairo_image_surface_get_height(s))
	stride := int(C.cairo_image_surface_get_stride(s))

	return newImg(s, format, width, height, stride)
}

//image surfaces will only ever have one key for image data.
var imgKey = &C.cairo_user_data_key_t{}

//big endian offsets for FromImage
var oA, oR, oG, oB = 0, 1, 2, 3

func init() {
	//flip offsets for little-endian
	t := uint32(1)
	if (*[4]byte)(unsafe.Pointer(&t))[0] == 1 {
		oB, oG, oR, oA = 0, 1, 2, 3
	}
}

//FromImage copies an image into a surface.
//
//The created image surface will have the same size as img,
//the optimal stride for img's width, and FormatARGB32.
//
//Originally cairo_image_surface_create_for_data and
//cairo_format_stride_for_width.
func FromImage(img image.Image) (ImageSurface, error) {
	f := FormatARGB32.c()
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	s := int(C.cairo_format_stride_for_width(f, C.int(w)))

	n := s * h
	data := (*C.uchar)(C.calloc(C.size_t(uintptr(n)), 1))
	pseudoslice := (*[1 << 30]C.uchar)(unsafe.Pointer(data))[:n:n]

	i := 0
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pseudoslice[i+oA] = C.uchar(a)
			pseudoslice[i+oR] = C.uchar(r)
			pseudoslice[i+oG] = C.uchar(g)
			pseudoslice[i+oB] = C.uchar(b)
			i += 4
		}
		i += 4 * (s/4 - w)
	}

	is := C.cairo_image_surface_create_for_data(data, f, C.int(w), C.int(h), C.int(s))
	C.cairo_surface_set_user_data(is, imgKey, unsafe.Pointer(data), free)

	return newImg(is, FormatARGB32, w, h, s)
}

//ToImage returns a copy of the surface as an image.
//
//Originally cairo_image_surface_get_data.
func (is ImageSurface) ToImage() (*image.RGBA, error) {
	if err := is.Err(); err != nil {
		return nil, err
	}
	C.cairo_surface_flush(is.s)

	data := C.cairo_image_surface_get_data(is.s)

	n := is.height * is.stride
	img := &image.RGBA{
		Pix:    make([]uint8, n),
		Stride: is.stride,
		Rect:   image.Rect(0, 0, is.width, is.height),
	}
	pseudoslice := (*[1 << 30]C.uchar)(unsafe.Pointer(data))[:n:n]

	for i := 0; i < n; i += 4 {
		img.Pix[i+0] = uint8(pseudoslice[i+oR])
		img.Pix[i+1] = uint8(pseudoslice[i+oG])
		img.Pix[i+2] = uint8(pseudoslice[i+oB])
		img.Pix[i+3] = uint8(pseudoslice[i+oA])
	}

	return img, nil
}

//Format reports the format of the surface.
//
//Originally cairo_image_surface_get_format.
func (is ImageSurface) Format() Format {
	return is.format
}

//Width reports the width of the surface in pixels.
//
//Originally cairo_image_surface_get_width.
func (is ImageSurface) Width() int {
	return is.width
}

//Height reports the height of the surface in pixels.
//
//Originally cairo_image_surface_get_height.
func (is ImageSurface) Height() int {
	return is.height
}

//Size returns the width and height of the image surface as a Point.
func (is ImageSurface) Size() Point {
	return Point{float64(is.width), float64(is.height)}
}

//Stride reports the stride of the image surface in number of bytes.
//
//Originally cairo_image_surface_get_stride.
func (is ImageSurface) Stride() int {
	return is.stride
}

var (
	mis    = map[id]*C.cairo_surface_t{} //the surface this surface was created from
	mismux = &sync.Mutex{}
)

//MappedImageSurface is a special ImageSurface created by Surface.Map.
type MappedImageSurface struct {
	ImageSurface
}

func registerImageSurface(m Surface, from *C.cairo_surface_t) {
	mismux.Lock()
	defer mismux.Unlock()
	//let original from be collected, we just need a general handle
	mis[m.id()] = C.cairo_surface_reference(from)
}

func newMappedImageSurface(s, from *C.cairo_surface_t) (m MappedImageSurface, err error) {
	im, err1 := cNewImageSurface(s)
	if err1 != nil {
		err = err1
		return
	}
	m = toMapped(im.(ImageSurface))
	err = m.Err()
	if err != nil {
		m.s = nil
	}
	registerImageSurface(m, from)
	return
}

func toMapped(i ImageSurface) MappedImageSurface {
	//Clear default finalizer so GC doesn't call surface_destroy.
	//We do not set a new finalizer on mappedImageSurface.Close,
	//because that would not do the right thing and the user is expected to unmap
	//manually when done.
	runtime.SetFinalizer(i.XtensionSurface, nil)
	return MappedImageSurface{
		ImageSurface: i,
	}
}

//Close always returns nil on a MappedImageSurface
func (m MappedImageSurface) Close() error {
	return nil
}

//Unmap uploads the content of the image to the target surface.
//Afterwards, the image is destroyed.
//
//Originally cairo_surface_unmap.
func (m MappedImageSurface) Unmap() error {
	err := m.Err()
	mismux.Lock()
	defer mismux.Unlock()
	from := mis[m.id()]
	C.cairo_surface_unmap_image(from, m.s)
	m.s = nil
	return err
}
