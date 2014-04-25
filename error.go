package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"errors"
)

const (
	errSuccess                 = C.CAIRO_STATUS_SUCCESS
	errNoMem                   = C.CAIRO_STATUS_NO_MEMORY
	errInvalidRestore          = C.CAIRO_STATUS_INVALID_RESTORE
	errInvalidPopGroup         = C.CAIRO_STATUS_INVALID_POP_GROUP
	errNoCurrentPoint          = C.CAIRO_STATUS_NO_CURRENT_POINT
	errInvalidMatrix           = C.CAIRO_STATUS_INVALID_MATRIX
	errInvalidStatus           = C.CAIRO_STATUS_INVALID_STATUS //seriously?
	errNullPointer             = C.CAIRO_STATUS_NULL_POINTER
	errInvalidString           = C.CAIRO_STATUS_INVALID_STRING
	errInvalidPathData         = C.CAIRO_STATUS_INVALID_PATH_DATA
	errReadError               = C.CAIRO_STATUS_READ_ERROR
	errWriteError              = C.CAIRO_STATUS_WRITE_ERROR
	errSurfaceFinished         = C.CAIRO_STATUS_SURFACE_FINISHED
	errSurfaceTypeMismatch     = C.CAIRO_STATUS_SURFACE_TYPE_MISMATCH
	errPatternTypeMismatch     = C.CAIRO_STATUS_PATTERN_TYPE_MISMATCH
	errInvalidContent          = C.CAIRO_STATUS_INVALID_CONTENT
	errInvalidFormat           = C.CAIRO_STATUS_INVALID_FORMAT
	errInvalidVisual           = C.CAIRO_STATUS_INVALID_VISUAL
	errFileNotFound            = C.CAIRO_STATUS_FILE_NOT_FOUND
	errInvalidDash             = C.CAIRO_STATUS_INVALID_DASH
	errInvalidDSCComment       = C.CAIRO_STATUS_INVALID_DSC_COMMENT
	errInvalidIndex            = C.CAIRO_STATUS_INVALID_INDEX
	errClipNotRepresentable    = C.CAIRO_STATUS_CLIP_NOT_REPRESENTABLE
	errTempFileError           = C.CAIRO_STATUS_TEMP_FILE_ERROR
	errInvalidStride           = C.CAIRO_STATUS_INVALID_STRIDE
	errFontTypeMismatch        = C.CAIRO_STATUS_FONT_TYPE_MISMATCH
	errUserFontImmutable       = C.CAIRO_STATUS_USER_FONT_IMMUTABLE
	errUserFontError           = C.CAIRO_STATUS_USER_FONT_ERROR
	errNegativeCount           = C.CAIRO_STATUS_NEGATIVE_COUNT
	errInvalidClusters         = C.CAIRO_STATUS_INVALID_CLUSTERS
	errInvalidSlant            = C.CAIRO_STATUS_INVALID_SLANT
	errInvalidWeight           = C.CAIRO_STATUS_INVALID_WEIGHT
	errInvalidSize             = C.CAIRO_STATUS_INVALID_SIZE
	errUserFontNotImplemented  = C.CAIRO_STATUS_USER_FONT_NOT_IMPLEMENTED
	errDeviceTypeMismatch      = C.CAIRO_STATUS_DEVICE_TYPE_MISMATCH
	errDeviceError             = C.CAIRO_STATUS_DEVICE_ERROR
	errInvalidMeshConstruction = C.CAIRO_STATUS_INVALID_MESH_CONSTRUCTION
	errDeviceFinished          = C.CAIRO_STATUS_DEVICE_FINISHED
	errLastStatus              = C.CAIRO_STATUS_LAST_STATUS
)

var (
	ErrInvalidPathData       = mkerr(errInvalidPathData)
	ErrInvalidDash           = mkerr(errInvalidDash)
	ErrInvalidLibcairoHandle = errors.New("invalid handle to libcairo resource")
)

func st2str(st C.cairo_status_t) string {
	return C.GoString(C.cairo_status_to_string(st))
}

func mkerr(st C.cairo_status_t) error {
	return errors.New(st2str(st))
}

func toerr(st C.cairo_status_t) error {
	return toerr_ided(st, nil)
}

func toerr_ided(st C.cairo_status_t, ider interface {
	id() id
}) error {
	switch int(st) {
	case errSuccess:
		return nil
	case errInvalidRestore, errInvalidPopGroup, errNoCurrentPoint, errInvalidMatrix, errInvalidString, errSurfaceFinished:
		panic(st2str(st))
	case errInvalidPathData:
		return ErrInvalidPathData
	case errInvalidDash:
		return ErrInvalidDash
	case errWriteError:
		mux.Lock()
		defer mux.Unlock()
		if w, ok := wmap[ider.id()]; ok {
			return w.err
		}
	}
	return errors.New(st2str(st))
}
