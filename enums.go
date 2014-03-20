package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo-pdf.h>
//#include <cairo/cairo-ps.h>
//#include <cairo/cairo-svg.h>
import "C"

//cairo_antialias_t
type antialias int

//Specifies the type of antialiasing to do when rendering text or shapes.
//
//As it is not necessarily clear from the above what advantages a particular
//antialias method provides, since libcairo 1.12, there is also a set of hints:
//	AntialiasFast
//		Allow the backend to degrade raster quality for speed
//	AntialiasGood
//		A balance between speed and quality
//	AntialiasBest
//		A high-fidelity, but potentially slow, raster mode
//
//These make no guarantee on how the backend will perform its rasterisation
//(if it even rasterises!), nor that they have any differing effect other
//than to enable some form of antialiasing. In the case of glyph rendering,
//AntialiasFast and AntialiasGood will be mapped to AntialiasGray,
//with AntialiasBest being equivalent to AntialiasSubpixel.
//
//The interpretation of AntialiasDefault is left entirely up to the backend,
//typically this will be similar to AntialiasGood.
//
//Originally cairo_antialias_t.
const (
	//AntialiasDefault uses the default antialiasing for the subsystem
	//and target device.
	AntialiasDefault antialias = C.CAIRO_ANTIALIAS_DEFAULT

	//AntialiasNone uses a bilevel alpha mask.
	AntialiasNone antialias = C.CAIRO_ANTIALIAS_NONE
	//AntialiasGray performs single-color antialiasing (using shades of gray
	//for black text on white background, for example).
	AntialiasGray antialias = C.CAIRO_ANTIALIAS_GRAY
	//AntialiasSubpixel performs antialiasing by taking advantage of the order
	//of subpixel elements on devices such as LCD panels.
	AntialiasSubpixel antialias = C.CAIRO_ANTIALIAS_SUBPIXEL

	//AntialiasFast is a hint that the backend should perform some antialiasing
	//but prefer speed over quality.
	AntialiasFast antialias = C.CAIRO_ANTIALIAS_FAST
	//AntialiasGood is a hint that the backend should balance quality against
	//performance.
	AntialiasGood antialias = C.CAIRO_ANTIALIAS_GOOD
	//AntialiasBest is a hint that the backend should render at the highest
	//quality, sacrificing speed if necessary.
	AntialiasBest antialias = C.CAIRO_ANTIALIAS_BEST
)

//BUG(jmf): need to check cairo source to make sure the antialias String method isn't crazy.

func (a antialias) String() string {
	if a == AntialiasNone {
		return "No antialiasing"
	}
	s := ""
	switch {
	case a&^AntialiasGray == 0:
		s += "Gray "
	case a&^AntialiasSubpixel == 0:
		s += "Supixel "
	}
	switch {
	case a&^AntialiasFast == 0:
		s += "Fast "
	case a&^AntialiasGood == 0:
		s += "Good "
	case a&^AntialiasBest == 0:
		s += "Best "
	}
	return s + "antialiasing"
}

//cairo_content_t
type content int

//Content is used to describe the content that a surface will contain, whether
//color information, alpha (translucence vs. opacity), or both.
//
//Originally cairo_content_t.
const (
	//ContentColor specifies that the surface will hold color content only.
	ContentColor content = C.CAIRO_CONTENT_COLOR
	//ContentAlpha specifies that the surface will hold alpha content only.
	ContentAlpha content = C.CAIRO_CONTENT_ALPHA
	//ContentColorAlpha specifies that the surface will hold color and alpha
	//content.
	ContentColorAlpha content = C.CAIRO_CONTENT_COLOR_ALPHA
)

func (c content) String() string {
	switch c {
	case ContentColor:
		return "Color content only"
	case ContentAlpha:
		return "Alpha content only"
	case ContentColorAlpha:
		return "Color and alpha content"
	}
	return "unknown content"
}

//cairo_device_type_t
type deviceType int

//A deviceType describes the type of a given device, also known as a "backend".
//
//A deviceType value has the following methods, in addition to String, which all return bool:
//
//	Native
//		A native device (win32, xcb, etc).
//
//	GL
//		OpenGL or COGL.
//
//	Pseudo
//		A device that doesn't output to a screen of some kind (XML).
//
//Originally cairo_device_type_t.
const (
	//DeviceTypeDRM is a Direct Render Manager device.
	DeviceTypeDRM deviceType = C.CAIRO_DEVICE_TYPE_DRM
	//DeviceTypeGL is an OpenGL device.
	DeviceTypeGL deviceType = C.CAIRO_DEVICE_TYPE_GL
	//DeviceTypeXCB is an XCB device.
	DeviceTypeXCB deviceType = C.CAIRO_DEVICE_TYPE_XCB
	//DeviceTypeXLib is an X lib device.
	DeviceTypeXLib deviceType = C.CAIRO_DEVICE_TYPE_XLIB
	//DeviceTypeXML is an XML device.
	DeviceTypeXML deviceType = C.CAIRO_DEVICE_TYPE_XML
	//DeviceTypeCOGL is a COGL device.
	DeviceTypeCOGL deviceType = C.CAIRO_DEVICE_TYPE_COGL
	//DeviceTypeWin32 is a Win32 device.
	DeviceTypeWin32 deviceType = C.CAIRO_DEVICE_TYPE_WIN32
)

var devstr = map[deviceType]string{
	DeviceTypeDRM:   "DRM",
	DeviceTypeGL:    "OpenGL",
	DeviceTypeXCB:   "XCB",
	DeviceTypeXLib:  "Xlib",
	DeviceTypeXML:   "XML",
	DeviceTypeCOGL:  "COGL",
	DeviceTypeWin32: "Win32",
}

func (d deviceType) String() string {
	s := devstr[d]
	if s == "" {
		s = "unknown"
	}
	return s + " device"
}

//Native returns true if the device type is a native platform type.
func (d deviceType) Native() bool {
	switch d {
	case DeviceTypeDRM, DeviceTypeXCB, DeviceTypeXLib, DeviceTypeWin32:
		return true
	}
	return false
}

//GL returns true if the device type is an OpenGL or COGL device
func (d deviceType) GL() bool {
	return d == DeviceTypeGL || d == DeviceTypeCOGL
}

//Pseudo returns true for pseudodevices (eg, XML).
func (d deviceType) Pseudo() bool {
	return d == DeviceTypeXML
}

//cairo_extend_t
type extend int

//The extend type describes how pattern color/alpha will be determined
//for areas "outside" the pattern's natural area, (for example, outside
//the surface bounds or outside the gradient geometry).
//
//Originally cairo_extend_t.
const (
	//ExtendNone makes pixels outside of the source pattern are fully transparent.
	ExtendNone extend = C.CAIRO_EXTEND_NONE
	//ExtendRepeat means the pattern is tiled by repeating.
	ExtendRepeat extend = C.CAIRO_EXTEND_REPEAT
	//ExtendReflect means the pattern is tiled by reflecting at the edges.
	ExtendReflect extend = C.CAIRO_EXTEND_REFLECT
	//ExtendPad means pixels outside of the pattern copy the closest pixel
	//from the source.
	ExtendPad extend = C.CAIRO_EXTEND_PAD
)

func (e extend) String() string {
	var s string
	switch e {
	case ExtendNone:
		s = "No"
	case ExtendRepeat:
		s = "Repeat"
	case ExtendReflect:
		s = "Relfect"
	case ExtendPad:
		s = "Pad"
	default:
		s = "unknown"
	}
	return s + " extend"
}

//cairo_fill_rule_t
type fillRule int

//The fillRule type is used to select how paths are filled.
//For both fill rules, whether or not a point is included in the fill
//is determined by taking a ray from that point to infinity and looking
//at intersections with the path.
//The ray can be in any direction, as long as it doesn't pass through
//the end point of a segment or have a tricky intersection
//such as intersecting tangent to the path.
//(Note that filling is not actually implemented in this way.
//This is just a description of the rule that is applied.)
//
//The default fillRule is FillRuleWinding.
//
//Originally cairo_fill_rule_t.
const (
	//FillRuleWinding works as follows:
	//If the path crosses the ray from left-to-right, counts +1.
	//If the path crosses the ray from right to left, counts -1.
	//(Left and right are determined from the perspective of looking along
	//the ray from the starting point.) If the total count is non-zero,
	//the point will be filled.
	FillRuleWinding fillRule = C.CAIRO_FILL_RULE_WINDING //default

	//FillRuleEvenOdd counts the total number of intersections,
	//without regard to the orientation of the contour.
	//If the total number of intersections is odd, the point will be filled.
	FillRuleEvenOdd fillRule = C.CAIRO_FILL_RULE_EVEN_ODD
)

func (f fillRule) String() string {
	switch f {
	case FillRuleWinding:
		return "winding fill rule"
	case FillRuleEvenOdd:
		return "even-odd fill rule"
	}
	return "unknown fill rule"
}

//cairo_filter_t
type filter int

//NB CAIRO_FILTER_GAUSSIAN is left off as the docs say it is currently unimplemented

//The filter type indicates what filtering should be applied when reading pixel
//values from patterns.
//
//Originally cairo_filter_t.
const (
	//FilterFast is a high performance filter with quality similar to
	//FilterNearest.
	FilterFast filter = C.CAIRO_FILTER_FAST

	//FilterGood is a reasonable performance filter, with quality similiar to
	//FilterBilinear.
	FilterGood filter = C.CAIRO_FILTER_GOOD

	//FilterBest is the highest quality filter, but may not be suitable
	//for interactive use.
	FilterBest filter = C.CAIRO_FILTER_BEST

	//FilterNearest is nearest-neighbor filtering.
	FilterNearest filter = C.CAIRO_FILTER_NEAREST

	//FilterBilinear uses linear interpolation in two dimensions.
	FilterBilinear filter = C.CAIRO_FILTER_BILINEAR
)

func (f filter) String() string {
	var s string
	switch f {
	case FilterFast:
		s = "Fast"
	case FilterGood:
		s = "Good"
	case FilterBest:
		s = "Best"
	case FilterNearest:
		s = "Nearest"
	case FilterBilinear:
		s = "Bilinear"
	default:
		s = "unknown"
	}
	return s + " filter"
}

//cairo_font_slant_t
type fontSlant int

//Specifies variants of a font face based on their slant.
//
//Originally cairo_font_slant_t.
const (
	//FontSlantNormal is standard upright font style.
	FontSlantNormal fontSlant = C.CAIRO_FONT_SLANT_NORMAL
	//FontSlantItalic is italic font style.
	FontSlantItalic fontSlant = C.CAIRO_FONT_SLANT_ITALIC
	//FontSlantOblique is oblique font style.
	FontSlantOblique fontSlant = C.CAIRO_FONT_SLANT_OBLIQUE
)

func (s fontSlant) String() string {
	switch s {
	case FontSlantNormal:
		return "normal font slant"
	case FontSlantItalic:
		return "italic font slant"
	case FontSlantOblique:
		return "oblique font slant"
	}
	return "unknown font slant"
}

//cairo_font_type_t
type fontType int

//A fontType describes the type of a given font face or scaled font.
//The font types are also known as "font backends" within cairo.
//
//Originally cairo_font_type_t.
const (
	//FontTypeToy fonts are created using cairo's toy font api.
	FontTypeToy fontType = C.CAIRO_FONT_TYPE_TOY
	//FontTypeWin32 is a native Windows font.
	FontTypeWin32 fontType = C.CAIRO_FONT_TYPE_WIN32
	//FontTypeQuartz is a native Macintosh font.
	FontTypeQuartz fontType = C.CAIRO_FONT_TYPE_QUARTZ //previously knonw as CAIRO_FONT_TYPE_ATSUI
	//FontTypeUser was created using cairo's user font api.
	FontTypeUser fontType = C.CAIRO_FONT_TYPE_USER
)

func (f fontType) String() string {
	s := ""
	switch f {
	case FontTypeToy:
		s = "toy"
	case FontTypeWin32:
		s = "Win32"
	case FontTypeQuartz:
		s = "Quartz"
	case FontTypeUser:
		s = "user"
	default:
		s = "unknown"
	}
	return "Font type " + s
}

//cairo_font_weight_t
type fontWeight int

//Specifies variants of a font face based on their weight.
//
//Orginally cairo_font_weight_t.
const (
	//FontWeightNormal is normal font weight.
	FontWeightNormal fontWeight = C.CAIRO_FONT_WEIGHT_NORMAL
	//FontWeightBold is bold font weight.
	FontWeightBold fontWeight = C.CAIRO_FONT_WEIGHT_BOLD
)

func (w fontWeight) String() string {
	switch w {
	case FontWeightNormal:
		return "normal font weight"
	case FontWeightBold:
		return "bold font weight"
	}
	return "unknown font weight"
}

//cairo_format_t
type format int

//A format identifies the memory format of image data.
//
//Originally cairo_format_t.
const (
	//FormatInvalid specifies an unsupported or nonexistent format.
	FormatInvalid format = C.CAIRO_FORMAT_INVALID

	//FormatARGB32 specifies that each pixel is a native-endian 32 bit quanity
	//listed as transparency, red, green, and then blue.
	FormatARGB32 format = C.CAIRO_FORMAT_ARGB32 //zero value

	//FormatRGB24 is the same as FormatARGB32 but the 8-bits of transparency
	//are unused.
	FormatRGB24 format = C.CAIRO_FORMAT_RGB24

	//FormatA8 stores each pixel in an 8-bit quantity holding an alpha value.
	FormatA8 format = C.CAIRO_FORMAT_A8

	//FormatA1 stores each pixel in a 1-bit quantity holding an alpha value.
	FormatA1 format = C.CAIRO_FORMAT_A1

	//FormatRGB16_565 stores each pixel as a 16-bit quantity with 5 bits for
	//red, 6 bits for green, and 5 bits for blue.
	FormatRGB16_565 format = C.CAIRO_FORMAT_RGB16_565

	//FormatRGB30 is like FormatRGB24 but with 10 bits per pixel instead
	//of 8.
	FormatRGB30 format = C.CAIRO_FORMAT_RGB30
)

func (f format) String() string {
	var s string
	switch f {
	case FormatARGB32:
		s = "32bit ARGB"
	case FormatRGB24:
		s = "24bit RGB"
	case FormatA8:
		s = "A8"
	case FormatA1:
		s = "A1"
	case FormatRGB16_565:
		s = "5-6-5 RGB16"
	case FormatRGB30:
		s = "RGB30"
	default: //grabs format invalid too
		s = "unknown"
	}
	return s + " format"
}

//cairo_hint_metrics_t
type hintMetrics int

//Specifies whether to hint font metrics; hinting font metrics means quantizing
//them so that they are integer values in device space. Doing this improves
//the consistency of letter and line spacing, however it also means that text
//will be laid out differently at different zoom factors.
//
//Oringally cairo_hint_metrics_t.
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

//Originally cairo_hint_style_t.
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
//
//Originally cairo_line_cap_t.
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

//Specifies how to render the junction of two lines when stroking.
//
//Originally cairo_line_join_t.
const (
	//LineJoinMiter uses a sharp (angled) corner.
	LineJoinMiter lineJoin = C.CAIRO_LINE_JOIN_MITER //default
	//LineJoinRound uses a rounded join, the center of the circle
	//is the join point.
	LineJoinRound lineJoin = C.CAIRO_LINE_JOIN_ROUND
	//LineJoinBevel uses a cut-off join, the join is cut off at half
	//the line width from the joint point.
	LineJoinBevel lineJoin = C.CAIRO_LINE_JOIN_BEVEL
)

//cairo_operator_t
type op int

//Originally cairo_operator_t.
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

//Originally cairo_path_data_type_t.
const (
	PathMoveTo    pathDataType = C.CAIRO_PATH_MOVE_TO
	PathLineTo    pathDataType = C.CAIRO_PATH_LINE_TO
	PathCurveTo   pathDataType = C.CAIRO_PATH_CURVE_TO
	PathClosePath pathDataType = C.CAIRO_PATH_CLOSE_PATH
)

//cairo_pdf_version_t
type pdfVersion int

//The pdfVersion type describes the version number of the PDF specification
//that a generated PDF file will conform to.
//
//Originally cairo_pdf_version_t.
const (
	//PDFVersion1_4 is the version 1.4 of the PDF specification.
	PDFVersion1_4 pdfVersion = C.CAIRO_PDF_VERSION_1_4
	//PDFVersion1_5 is the version 1.5 of the PDF specification.
	PDFVersion1_5 pdfVersion = C.CAIRO_PDF_VERSION_1_5
)

func (p pdfVersion) String() string {
	return C.GoString(C.cairo_pdf_version_to_string(C.cairo_pdf_version_t(p)))
}

//cairo_ps_level_t
type psLevel int

//The psLevel type is used to describe the version number of the PDF
//specification that a generated PDF file will conform to.
//
//Since libcairo 1.6. Originally cairo_ps_level_t.
const (
	//PSLevel2 is the language level 2 of the PostScript specification.
	PSLevel2 psLevel = C.CAIRO_PS_LEVEL_2
	//PSLevel3 is the language level 3 of the PostScript specification.
	PSLevel3 psLevel = C.CAIRO_PS_LEVEL_3
)

func (p psLevel) String() string {
	return C.GoString(C.cairo_ps_level_to_string(C.cairo_ps_level_t(p)))
}

//cairo_status_t is handled in error.go
//BUG(jmf): make this not a lie ^

//cairo_subpixel_order_t
type subpixelOrder int

//Originally cairo_subpixel_order_t.
const (
	SubpixelOrderDefault subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_DEFAULT
	SubpixelOrderRGB     subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_RGB
	SubpixelOrderBGR     subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_BGR
	SubpixelOrderVRGB    subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_VRGB
	SubpixelOrderVBGR    subpixelOrder = C.CAIRO_SUBPIXEL_ORDER_VBGR
)

//cairo_surface_type_t
type surfaceType int

//Originally cairo_surface_type_t.
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
