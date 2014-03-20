package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo-pdf.h>
//#include <cairo/cairo-ps.h>
//#include <cairo/cairo-svg.h>
import "C"

//cairo_antialias_t
type antialias int

const (
	AntialiasDefault antialias = C.CAIRO_ANTIALIAS_DEFAULT

	//Methods
	AntialiasNone     antialias = C.CAIRO_ANTIALIAS_NONE
	AntialiasGray     antialias = C.CAIRO_ANTIALIAS_GRAY
	AntialiasSubPixel antialias = C.CAIRO_ANTIALIAS_SUBPIXEL

	//Hints
	AntialiasFast antialias = C.CAIRO_ANTIALIAS_FAST
	AntialiasGood antialias = C.CAIRO_ANTIALIAS_GOOD
	AntialiasBest antialias = C.CAIRO_ANTIALIAS_BEST
)

//cairo_content_t
type content int

const (
	ContentColor      content = C.CAIRO_CONTENT_COLOR
	ContentAlpha      content = C.CAIRO_CONTENT_ALPHA
	ContentColorAlpha content = C.CAIRO_CONTENT_COLOR_ALPHA
)

//cairo_device_type_t
type deviceType int

const (
	DeviceTypeDRM   deviceType = C.CAIRO_DEVICE_TYPE_DRM
	DeviceTypeGl    deviceType = C.CAIRO_DEVICE_TYPE_GL
	DeviceTypeXCB   deviceType = C.CAIRO_DEVICE_TYPE_XCB
	DeviceTypeXLib  deviceType = C.CAIRO_DEVICE_TYPE_XLIB
	DeviceTypeXML   deviceType = C.CAIRO_DEVICE_TYPE_XML
	DeviceTypeCOGL  deviceType = C.CAIRO_DEVICE_TYPE_COGL
	DeviceTypeWin32 deviceType = C.CAIRO_DEVICE_TYPE_WIN32
)

//cairo_extend_t
type extend int

const (
	ExtendNone    extend = C.CAIRO_EXTEND_NONE
	ExtendRepeat  extend = C.CAIRO_EXTEND_REPEAT
	ExtendReflect extend = C.CAIRO_EXTEND_REFLECT
	ExtendPad     extend = C.CAIRO_EXTEND_PAD
)

//cairo_fill_rule_t
type fillRule int

const (
	FillRuleWinding fillRule = C.CAIRO_FILL_RULE_WINDING
	FillRuleEvenOdd fillRule = C.CAIRO_FILL_RULE_EVEN_ODD
)

//cairo_filter_t
type filter int

const (
	FilterFast     filter = C.CAIRO_FILTER_FAST
	FilterGood     filter = C.CAIRO_FILTER_GOOD
	FilterBest     filter = C.CAIRO_FILTER_BEST
	FilterNearest  filter = C.CAIRO_FILTER_NEAREST
	FilterBilinear filter = C.CAIRO_FILTER_BILINEAR
	FilterGaussian filter = C.CAIRO_FILTER_GAUSSIAN
)

//cairo_font_slant_t
type fontSlant int

//Specifies variants of a font face based on their slant.
const (
	FontSlantNormal  fontSlant = C.CAIRO_FONT_SLANT_NORMAL
	FontSlantItalic  fontSlant = C.CAIRO_FONT_SLANT_ITALIC
	FontSlantOblique fontSlant = C.CAIRO_FONT_SLANT_OBLIQUE
)

//cairo_font_type_t
type fontType int

const (
	FontTypeToy      fontType = C.CAIRO_FONT_TYPE_TOY
	FontTypeFreeType fontType = C.CAIRO_FONT_TYPE_FT
	FontTypeWin32    fontType = C.CAIRO_FONT_TYPE_WIN32
	FontTypeQuartz   fontType = C.CAIRO_FONT_TYPE_QUARTZ //previous CAIRO_FONT_TYPE_ATSUI
	FontTypeUser     fontType = C.CAIRO_FONT_TYPE_USER
)

//cairo_font_weight_t
type fontWeight int

//Specifies variants of a font face based on their weight.
const (
	FontWeightNormal fontWeight = C.CAIRO_FONT_WEIGHT_NORMAL
	FontWeightBold   fontWeight = C.CAIRO_FONT_WEIGHT_BOLD
)

//cairo_format_t
type format int

const (
	FormatInvalid   format = C.CAIRO_FORMAT_INVALID
	FormatARGB32    format = C.CAIRO_FORMAT_ARGB32 //zero value
	FormatRGB24     format = C.CAIRO_FORMAT_RGB24
	FormatA8        format = C.CAIRO_FORMAT_A8
	FormatA1        format = C.CAIRO_FORMAT_A1
	FormatRGB16_565 format = C.CAIRO_FORMAT_RGB16_565
	FormatRGB30     format = C.CAIRO_FORMAT_RGB30
)

//cairo_hint_metrics_t
type hintMetrics int

//Specifies whether to hint font metrics; hinting font metrics means quantizing
//them so that they are integer values in device space. Doing this improves
//the consistency of letter and line spacing, however it also means that text
//will be laid out differently at different zoom factors.
const (
	//HintMetricsDefault use hint metrics in the default manner
	//for the font backend and target device.
	HintMetricsDefault hintMetrics = C.CAIRO_HINT_METRICS_DEFAULT
	//HintMetricsOff does not hint font metrics.
	HintMetricsOff hintMetrics = C.CAIRO_HINT_METRICS_OFF
	//HintMetricsOn hints font metrics.
	HintMetricsOn hintMetrics = C.CAIRO_HINT_METRICS_ON
)

//cairo_hint_style_t
type hintStyle int

const (
	HintStyleDefault hintStyle = C.CAIRO_HINT_STYLE_DEFAULT
	HintStyleNone    hintStyle = C.CAIRO_HINT_STYLE_NONE
	HintStyleSlight  hintStyle = C.CAIRO_HINT_STYLE_SLIGHT
	HintStyleMedium  hintStyle = C.CAIRO_HINT_STYLE_MEDIUM
	HintStyleFull    hintStyle = C.CAIRO_HINT_STYLE_FULL
)

//cairo_line_cap_t
type lineCap int

//Specifies how to render the endpoints of the path when stroking.
const (
	//LineCapButt starts(stops) the line exactly at the start(end) point.
	LineCapButt lineCap = C.CAIRO_LINE_CAP_BUTT
	//LineCapRound uses a round ending, the center of the circle is the end point.
	LineCapRound lineCap = C.CAIRO_LINE_CAP_ROUND
	//LineCapSquare uses a squared ending, the center of the square is
	//the end point.
	LineCapSquare lineCap = C.CAIRO_LINE_CAP_SQUARE
)

//cairo_line_join_t
type lineJoin int

const (
	LineJoinMiter lineJoin = C.CAIRO_LINE_JOIN_MITER
	LineJoinRound lineJoin = C.CAIRO_LINE_JOIN_ROUND
	LineJoinBevel lineJoin = C.CAIRO_LINE_JOIN_BEVEL
)

//cairo_operator_t
type op int

const (
	OpClear op = C.CAIRO_OPERATOR_CLEAR

	OpSource op = C.CAIRO_OPERATOR_SOURCE //default
	OpOver   op = C.CAIRO_OPERATOR_OVER
	OpIn     op = C.CAIRO_OPERATOR_IN
	OpOut    op = C.CAIRO_OPERATOR_OUT
	OpAtop   op = C.CAIRO_OPERATOR_ATOP

	OpDest     op = C.CAIRO_OPERATOR_DEST
	OpDestOver op = C.CAIRO_OPERATOR_DEST_OVER
	OpDestIn   op = C.CAIRO_OPERATOR_DEST_IN
	OpDestOut  op = C.CAIRO_OPERATOR_DEST_OUT
	OpDestAtop op = C.CAIRO_OPERATOR_DEST_ATOP

	OpXor      op = C.CAIRO_OPERATOR_XOR
	OpAdd      op = C.CAIRO_OPERATOR_ADD
	OpSaturate op = C.CAIRO_OPERATOR_SATURATE

	OpMultiply      op = C.CAIRO_OPERATOR_MULTIPLY
	OpScreen        op = C.CAIRO_OPERATOR_SCREEN
	OpOverlay       op = C.CAIRO_OPERATOR_OVERLAY
	OpDarken        op = C.CAIRO_OPERATOR_DARKEN
	OpLighten       op = C.CAIRO_OPERATOR_LIGHTEN
	OpColorDodge    op = C.CAIRO_OPERATOR_COLOR_DODGE
	OpColorBurn     op = C.CAIRO_OPERATOR_COLOR_BURN
	OpHardLight     op = C.CAIRO_OPERATOR_HARD_LIGHT
	OpSoftLight     op = C.CAIRO_OPERATOR_SOFT_LIGHT
	OpDifference    op = C.CAIRO_OPERATOR_DIFFERENCE
	OpExclusion     op = C.CAIRO_OPERATOR_EXCLUSION
	OpHueHSL        op = C.CAIRO_OPERATOR_HSL_HUE
	OpSaturationHSL op = C.CAIRO_OPERATOR_HSL_SATURATION
	OpColorHSL      op = C.CAIRO_OPERATOR_HSL_COLOR
	OpLuminosityHSL op = C.CAIRO_OPERATOR_HSL_LUMINOSITY
)

//cairo_path_data_type_t
type pathDataType int

const (
	PathDataTypeMoveTo    pathDataType = C.CAIRO_PATH_MOVE_TO
	PathDataTypeLineTo    pathDataType = C.CAIRO_PATH_LINE_TO
	PathDataTypeCurveTo   pathDataType = C.CAIRO_PATH_CURVE_TO
	PathDataTypeClosePath pathDataType = C.CAIRO_PATH_CLOSE_PATH
)

//cairo_pdf_version_t
type pdfVersion int

const (
	PDFVersion1_4 pdfVersion = C.CAIRO_PDF_VERSION_1_4
	PDFVersion1_5 pdfVersion = C.CAIRO_PDF_VERSION_1_5
)

//cairo_ps_level_t
type psLevel int

const (
	PSLevel2 psLevel = C.CAIRO_PS_LEVEL_2
	PSLevel3 psLevel = C.CAIRO_PS_LEVEL_3
)

//cairo_status_t is handled in error.go
//BUG(jmf): make this not a lie ^

//cairo_subpixel_order_t
type subpixelOrder int

const (
	SubpixelOrderDefault subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_DEFAULT
	SubpixelOrderRGB     subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_RGB
	SubpixelOrderBGR     subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_BGR
	SubpixelOrderVRGB    subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_VRGB
	SubpixelOrderVBGR    subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_VBGR
)

//cairo_surface_type_t
type surfaceType int

const (
	SurfaceTypeImage         surfaceType = C.CAIRO_SURFACE_TYPE_IMAGE
	SurfaceTypePDF           surfaceType = C.CAIRO_SURFACE_TYPE_PDF
	SurfaceTypePS            surfaceType = C.CAIRO_SURFACE_TYPE_PS
	SurfaceTypeXLib          surfaceType = C.CAIRO_SURFACE_TYPE_XLIB
	SurfaceTypeXCB           surfaceType = C.CAIRO_SURFACE_TYPE_XCB
	SurfaceTypeGlitz         surfaceType = C.CAIRO_SURFACE_TYPE_GLITZ
	SurfaceTypeQuartz        surfaceType = C.CAIRO_SURFACE_TYPE_QUARTZ
	SurfaceTypeWin32         surfaceType = C.CAIRO_SURFACE_TYPE_WIN32
	SurfaceTypeBeOS          surfaceType = C.CAIRO_SURFACE_TYPE_BEOS
	SurfaceTypeDirectFB      surfaceType = C.CAIRO_SURFACE_TYPE_DIRECTFB
	SurfaceTypeSVG           surfaceType = C.CAIRO_SURFACE_TYPE_SVG
	SurfaceTypeOS2           surfaceType = C.CAIRO_SURFACE_TYPE_OS2
	SurfaceTypeWin32Printing surfaceType = C.CAIRO_SURFACE_TYPE_WIN32_PRINTING
	SurfaceTypeQuartzImage   surfaceType = C.CAIRO_SURFACE_TYPE_QUARTZ_IMAGE
	SurfaceTypeQT            surfaceType = C.CAIRO_SURFACE_TYPE_QT
	SurfaceTypeRecording     surfaceType = C.CAIRO_SURFACE_TYPE_RECORDING
	SurfaceTypeVG            surfaceType = C.CAIRO_SURFACE_TYPE_VG
	SurfaceTypeGL            surfaceType = C.CAIRO_SURFACE_TYPE_GL
	SurfaceTypeDRM           surfaceType = C.CAIRO_SURFACE_TYPE_DRM
	SurfaceTypeTee           surfaceType = C.CAIRO_SURFACE_TYPE_TEE
	SurfaceTypeXML           surfaceType = C.CAIRO_SURFACE_TYPE_XML
	SurfaceTypeSkia          surfaceType = C.CAIRO_SURFACE_TYPE_SKIA
	SurfaceTypeSubsurface    surfaceType = C.CAIRO_SURFACE_TYPE_SUBSURFACE
	SurfaceTypeCOGL          surfaceType = C.CAIRO_SURFACE_TYPE_COGL
)

//BUG(jmf): document all "enums"
//BUG(jmf): add String methods for all enums
