package mimepattern

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

type mime string

//These are the MIME types libcairo supports for embedding.
const (
	//PNG is the Portable Network Graphics image file format (ISO/IEC 15948).
	//
	//Originally CAIRO_MIME_TYPE_PNG.
	PNG mime = "image/png"
	//JPEG is the Joint Photographic Experts Group (JPEG) image coding standard (ISO/IEC 10918-1).
	//
	//Originally CAIRO_MIME_TYPE_JPEG.
	JPEG mime = "image/jpeg"
	//JP2 is the Joint Photographic Experts Group (JPEG) 2000 image coding standard (ISO/IEC 15444-1).
	//
	//Originally CAIRO_MIME_TYPE_JP2.
	JP2 mime = "image/jp2"

	uri mime = "text/x-uri"
)

func (m mime) valid() bool {
	switch m {
	case PNG, JPEG, JP2:
		return true
	}
	return false
}

func (m mime) s() string {
	return string(m)
}

var m2c = map[mime]*C.char{
	PNG:  C.CString(PNG.s()),
	JPEG: C.CString(JPEG.s()),
	JP2:  C.CString(JP2.s()),
	uri:  C.CString(uri.s()),
}

func (m mime) c() (cs *C.char) {
	return m2c[m]
}
