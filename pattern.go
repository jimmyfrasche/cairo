package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"errors"
	"image/color"
	"runtime"
)

func getPatternType(p *C.cairo_pattern_t) patternType {
	return patternType(C.cairo_pattern_get_type(p))
}

func cPattern(p *C.cairo_pattern_t) (Pattern, error) {
	switch getPatternType(p) {
	case PatternTypeSolid:
		return cNewSolidPattern(p), nil
	case PatternTypeSurface:
		return cNewSurfacePattern(p)
	case PatternTypeLinear:
		return cNewLinearGradient(p), nil
	case PatternTypeRadial:
		return cNewRadialGradient(p), nil
	case PatternTypeMesh:
		return cNewMesh(p), nil
	case PatternTypeRasterSource:
		return nil, errors.New("unimplemented pattern type")
	}
	return nil, errors.New("unimplemented pattern type")
}

//BUG(jmf): does not support user patterns

//Pattern is a pattern used for drawing.
//
//Originally cairo_pattern_t.
type Pattern interface {
	Type() patternType
	Err() error
	Close() error

	SetExtend(extend)
	Extend() extend
	SetFilter(filter)
	Filter() filter
	SetMatrix(Matrix)
	Matrix() Matrix

	c() *C.cairo_pattern_t
}

type pattern struct {
	p *C.cairo_pattern_t
}

func newPattern(p *C.cairo_pattern_t) *pattern {
	P := &pattern{p}
	runtime.SetFinalizer(P, (*pattern).Close)
	return P
}

func (p *pattern) c() *C.cairo_pattern_t {
	return p.p
}

//Type returns the type of the pattern.
//
//Originally cairo_pattern_get_type.
func (p *pattern) Type() patternType {
	return getPatternType(p.p)
}

//Err reports any error on this pattern.
//
//Originally cairo_pattern_status.
func (p *pattern) Err() error {
	if p.p == nil {
		return ErrInvalidLibcairoHandle
	}
	return toerr(C.cairo_pattern_status(p.p))
}

//Close releases the pattern's resources.
//
//Originally cairo_pattern_destroy.
func (p *pattern) Close() error {
	if p.p == nil {
		return nil
	}
	err := p.Err()
	runtime.SetFinalizer(p, nil)
	C.cairo_pattern_destroy(p.p)
	p.p = nil
	return err
}

//SetExtend sets the mode used for drawing outside the area of this pattern.
//
//Originally cairo_pattern_set_extend.
func (p *pattern) SetExtend(e extend) {
	C.cairo_pattern_set_extend(p.p, e.c())
}

//Extend reports the mode used for drawing outside the area of this pattern.
//
//Originally cairo_pattern_get_extend.
func (p *pattern) Extend() extend {
	return extend(C.cairo_pattern_get_extend(p.p))
}

//SetFilter sets the filter used when resizing this pattern.
//
//Originally cairo_pattern_set_filter.
func (p *pattern) SetFilter(f filter) {
	C.cairo_pattern_set_filter(p.p, f.c())
}

//Filter returns the filter used when resizing this pattern.
//
//Originally cairo_pattern_get_filter.
func (p *pattern) Filter() filter {
	return filter(C.cairo_pattern_get_filter(p.p))
}

//SetMatrix sets the pattern's transformation matrix.
//This matrix is a transformation from user space to pattern space.
//
//When a pattern is first created it always has the identity matrix for its
//transformation matrix, which means that pattern space is initially identical
//to user space.
//
//Important
//
//Please note that the direction of this transformation matrix is from user
//space to pattern space.
//This means that if you imagine the flow from a pattern to user space
//(and on to device space), then coordinates in that flow will be transformed
//by the inverse of the pattern matrix.
//
//Originally cairo_pattern_set_matrix.
func (p pattern) SetMatrix(m Matrix) {
	C.cairo_pattern_set_matrix(p.p, &m.m)
}

//Matrix returns this patterns transformation matrix.
//
//Originally cairo_pattern_get_matrix.
func (p pattern) Matrix() Matrix {
	var m C.cairo_matrix_t
	C.cairo_pattern_get_matrix(p.p, &m)
	return Matrix{m}
}

//SolidPattern is a Pattern corresponding to a single translucent color.
type SolidPattern struct {
	*pattern
	col AlphaColor
}

//NewSolidPattern creates a solid pattern of color c.
//
//Originally cairo_pattern_create_rgba.
func NewSolidPattern(c color.Color) SolidPattern {
	col := colorToAlpha(c)
	r, g, b, a := col.c()
	p := C.cairo_pattern_create_rgba(r, g, b, a)
	return SolidPattern{
		pattern: newPattern(p),
		col:     col,
	}
}

func cNewSolidPattern(p *C.cairo_pattern_t) Pattern {
	var r, g, b, a C.double
	C.cairo_pattern_get_rgba(p, &r, &g, &b, &a)
	return SolidPattern{
		pattern: newPattern(p),
		col:     cColor(r, g, b, a),
	}
}

//Color returns the color this pattern was created with.
//
//Originally cairo_pattern_get_rgba.
func (s SolidPattern) Color() AlphaColor {
	return s.col
}

//SurfacePattern is a Pattern backed by a Surface.
type SurfacePattern struct {
	*pattern
	s Surface
}

//NewSurfacePattern creates a Pattern from a Surface.
//
//Originally cairo_pattern_create_for_surface.
func NewSurfacePattern(s Surface) (sp SurfacePattern, err error) {
	if err = s.Err(); err != nil {
		return
	}
	r := s.XtensionRaw()
	p := C.cairo_pattern_create_for_surface(r)
	sp = SurfacePattern{
		pattern: newPattern(p),
		s:       s,
	}
	return sp, sp.Err()
}

//BUG(jmf): (potentially) assuming pattern returned by cairo_pattern_create_for_surface
//is the same as the pattern put into it. If this is not true, things could get messy.

func cNewSurfacePattern(p *C.cairo_pattern_t) (Pattern, error) {
	var s *C.cairo_surface_t
	C.cairo_pattern_get_surface(p, &s)
	C.cairo_surface_reference(s) //returned surface does not up libcairo refcount
	S, err := XtensionRevivifySurface(p)
	if err != nil {
		return nil, err
	}
	P := SurfacePattern{
		pattern: newPattern(p),
		s:       S,
	}
	return P, nil
}

//Surface returns the Surface of this Pattern.
//
//Originally cairo_pattern_get_surface.
func (s SurfacePattern) Surface() Surface {
	return s.s
}

//A ColorStop is the color of a single gradient stop.
//
//Note that when defining gradients it two, or more, stops are specified
//with identical offset values, they will be sorted according to the order
//in which the stops are added.
//Stops added earlier will compare less than stops added later.
//This can be useful for reliably making sharp color transitions
//instead of the typical blend.
type ColorStop struct {
	//Offset specifies the location of this color stop along the gradient's
	//control vector.
	Offset float64
	Color  color.Color
}

func (c ColorStop) c() (o, r, g, b, a C.double) {
	o = C.double(clamp01(c.Offset))
	r, g, b, a = colorToAlpha(c.Color).c()
	return
}

//Gradient is a linear or radial gradient.
type Gradient interface {
	Pattern
	ColorStops() []ColorStop
	addColorStops(cs []ColorStop)
}

type patternGradient struct {
	*pattern
}

func (p patternGradient) addColorStops(cs []ColorStop) {
	for _, c := range cs {
		o, r, g, b, a := c.c()
		C.cairo_pattern_add_color_stop_rgba(p.p, o, r, g, b, a)
	}
}

//ColorStops reports the color stops defined on this gradient.
//
//Originally cairo_pattern_get_color_stop_count and
//cairo_pattern_get_color_stop_rgba.
func (p patternGradient) ColorStops() (cs []ColorStop) {
	var i, n C.int
	//only returns error if not a gradient, but disallowed by construction.
	_ = C.cairo_pattern_get_color_stop_count(p.p, &n)
	cs = make([]ColorStop, 0, int(n))
	for ; i < n; i++ {
		var o, r, g, b, a C.double
		//only returns error if not a gradient or invalid index, but disallowed by construction.
		_ = C.cairo_pattern_get_color_stop_rgba(p.p, i, &o, &r, &g, &b, &a)
		cs = append(cs, ColorStop{
			Offset: float64(o),
			Color:  cColor(r, g, b, a),
		})
	}

	return
}

//BUG(jmf): get_color_stop says "Values of index are 0 to 1 less than the number returned by" count, but what does that imply?

//LinearGradient is a linear gradient pattern.
type LinearGradient struct {
	patternGradient
	start, end Point
}

//NewLinearGradient creates a new linear gradient, from start to end,
//with specified color stops.
//
//Originally cairo_pattern_create_linear and cairo_pattern_add_color_stop_rgba.
func NewLinearGradient(start, end Point, colorStops ...ColorStop) LinearGradient {
	x0, y0 := start.c()
	x1, y1 := end.c()
	p := C.cairo_pattern_create_linear(x0, y0, x1, y1)
	P := patternGradient{
		pattern: newPattern(p),
	}
	P.addColorStops(colorStops)
	return LinearGradient{
		patternGradient: P,
		start:           start,
		end:             end,
	}
}

func cNewLinearGradient(p *C.cairo_pattern_t) Pattern {
	var x0, y0, x1, y1 C.double
	C.cairo_pattern_get_linear_points(p, &x0, &y0, &x1, &y1)
	return LinearGradient{
		patternGradient: patternGradient{
			pattern: newPattern(p),
		},
		start: cPt(x0, y0),
		end:   cPt(x1, y1),
	}
}

//Line returns the start and end points of this linear gradient.
//
//Originally cairo_pattern_get_linear_points.
func (l LinearGradient) Line() (start, end Point) {
	return l.start, l.end
}

//RadialGradient is a radial gradient pattern.
type RadialGradient struct {
	patternGradient
	start, end Circle
}

//NewRadialGradient creates a new radial gradient between the circles
//start and end, with specified color stops.
//
//Originally cairo_pattern_create_radial and cairo_pattern_add_color_stop_rgba.
func NewRadialGradient(start, end Circle, colorStops ...ColorStop) RadialGradient {
	x0, y0, r0 := start.c()
	x1, y1, r1 := end.c()
	p := C.cairo_pattern_create_radial(x0, y0, r0, x1, y1, r1)
	P := patternGradient{
		pattern: newPattern(p),
	}
	P.addColorStops(colorStops)
	return RadialGradient{
		patternGradient: P,
		start:           start,
		end:             end,
	}
}

func cNewRadialGradient(p *C.cairo_pattern_t) Pattern {
	var x0, y0, r0, x1, y1, r1 C.double
	C.cairo_pattern_get_radial_circles(p, &x0, &y0, &r0, &x1, &y1, &r1)
	return RadialGradient{
		patternGradient: patternGradient{
			pattern: newPattern(p),
		},
		start: cCirc(x0, y0, r0),
		end:   cCirc(x0, y0, r0),
	}
}

//RadialCircles reports the gradient endpoints.
//
//Originally cairo_pattern_get_radial_circles.
func (r RadialGradient) RadialCircles() (start, end Circle) {
	return r.start, r.end
}

//Mesh is a mesh pattern.
//
//Mesh patterns are tensor-product patch meshes (type 7 shadings in PDF).
//Mesh patterns may also be used to create other types of shadings that are
//special cases of tensor-product patch meshes such as Coons patch meshes
//(type 6 shading in PDF) and Gouraud-shaded triangle meshes
//(type 4 and 5 shadings in PDF).
//
//Mesh patterns consist of one or more tensor-product patches,
//which should be defined before using the mesh pattern.
//Using a mesh pattern with a partially defined patch as source or mask
//will put the context in an error
//
//Definition
//
//A tensor-product patch is defined by 4 Bézier curves (side 0, 1, 2, 3)
//and by 4 additional control points (P₀, P₁, P₂, P₃) that provide further
//control over the patch and complete the definition of the tensor-product
//patch.
//The corner C₀ is the first point of the patch.
//
//All methods that take a control point or corner point index are taken mod 4.
//
//	      C₁     Side 1       C₂
//	       +---------------+
//	       |               |
//	       |  P₁       P₂  |
//	       |               |
//	Side 0 |               | Side 2
//	       |               |
//	       |               |
//	       |  P₀       P₃  |
//	       |               |
//	       +---------------+
// 	    C₀     Side 3        C₃
//
//Degenerate sides are permitted so straight lines may be used.
//A zero length line on one side may be used to create 3 sided patches.
//
//Each patch is constructed by calling BeginPatch, followed by MoveTo
//to specify the first point in the patch C₀.
//The sides are then specified by calls to CurveTo and LineTo.
//
//The four additional control points (P₀, P₁, P₂, P₃) in a patch can be
//specified with SetControlPoint.
//
//At each corner of the patch (C₀, C₁, C₂, C₃) a color may be specified
//with SetCornerColor.
//
//Note: The coordinates are always in pattern space. For a new pattern,
//pattern space is identical to user space, but the relationship between
//the spaces can be changed with SetMatrix.
//
//Special cases
//
//A Coons patch is a special case of the tensor-product patch
//where the control points are implicitly defined by the sides of the patch.
//The default value for any control point not specified is the implicit value
//for a Coons patch, i.e. if no control points are specified the patch is a
//Coons patch.
//
//A triangle is a special case of the tensor-product patch where the control
//points are implicitly defined by the sides of the patch, all the sides are
//lines and one of them has length 0.
//That is, if the patch is specified using just 3 lines, it is a triangle.
//
//If the corners connected by the 0-length side have the same color, the patch
//is a Gouraud-shaded triangle.
//
//Orientation
//
//Patches may be oriented differently to the above diagram.
//For example, the first point could be at the top left.
//The diagram only shows the relationship between the sides, corners and control
//points.
//
//Regardless of where the first point is located, when specifying colors,
//corner 0 will always be the first point, corner 1 the point between side 0
//and side 1, and so on.
//
//Defaults
//
//Calling EndPatch completes the current patch.
//If less than 4 sides have been defined, the first missing side is defined
//as a line from the current point to the first point of the patch (C₀)
//and the other sides are degenerate lines from C₀ to C₀.
//The corners between the added sides will all be coincident with C₀
//of the patch and their color will be set to be the same as the color of C₀.
//
//Any corner color whose color is not explicitly specified defaults to
//transparent black.
//
//Multiple Patches
//
//Additional patches may be added with additional calls to BeginPatch/EndPatch.
//
//	red, green, blue := Color{R: 1}, Color{G: 1}, Color{B: 1}
//
//	m := NewMesh()
//
//	//Add a Coons patch
//	err := m.BeginPatch().
//		MoveTo(ZP).
//		CurveTo(Pt(30, -30), Pt(60, 30), Pt(100, 0)).
//		CurveTo(Pt(60, 30), Pt(130, 60), Pt(100, 100)).
//		CurveTo(Pt(60, 70), Pt(30, 130), Pt(0, 100)).
//		CurveTo(Pt(30, 70), Pt(-30, 30), ZP).
//		SetCornerColor(0, red).
//		SetCornerColor(1, green).
//		SetCornerColor(2, blue).
//		SetCornerColor(3, Color{R: 1, G: 1}).
//		EndPatch()
//
//	if err != nil {
//		panic(err) //this is not how to handle errors, outside of examples
//	}
//
//	//Add a Gouraud-shaded triangle
//	err = m.BeginPatch().
//		MoveTo(Pt(100, 100)).
//		LineTo(Pt(130, 130)).
//		LineTo(Pt(130, 70)).
//		SetCornerColor(0, red).
//		SetCornerColor(1, green).
//		SetCornerColor(2, blue).
//		EndPatch()
//
//	if err != nil {
//		panic(err)
//	}
//
//When two patches overlap, the last one that has been added is drawn over
//the first one.
//
//When a patch folds over itself, points are sorted depending on their parameter
//coordinates inside the patch.
//The v coordinate ranges from 0 to 1 when moving from side 3 to side 1;
//the u coordinate ranges from 0 to 1 when going from side 0 to side 2.
//Points with higher v coordinate hide points with lower v coordinate.
//When two points have the same v coordinate, the one with higher u coordinate
//is above.
//This means that points nearer to side 1 are above points nearer to side 3;
//when this is not sufficient to decide which point is above
//(for example when both points belong to side 1 or side 3)
//points nearer to side 2 are above points nearer to side 0.
//
//More information
//
//For a complete definition of tensor-product patches,
//see the PDF specification (ISO32000), which describes
//the parametrization in detail.
type Mesh struct {
	*pattern
}

//BUG(jmf): v, u coordinates in mesh docs are undefined. Presumably vectors, but not clear in original docs.

//BUG(jmf): add line and curve types to geometry?

//NewMesh creates a new mesh pattern.
//
//Originally cairo_pattern_create_mesh.
func NewMesh() Mesh {
	p := C.cairo_pattern_create_mesh()
	return cNewMesh(p)
}

func cNewMesh(p *C.cairo_pattern_t) Mesh {
	return Mesh{
		pattern: newPattern(p),
	}
}

//BeginPatch creates a new patch.
//
//Originally cairo_mesh_pattern_begin_patch.
func (m Mesh) BeginPatch() Mesh {
	C.cairo_mesh_pattern_begin_patch(m.p)
	return m
}

//EndPatch closes the current patch and reports any error.
//
//Originally cairo_mesh_pattern_end_patch.
func (m Mesh) EndPatch() error {
	C.cairo_mesh_pattern_end_patch(m.p)
	return m.Err()
}

//MoveTo defines the first point of the current patch in the mesh.
//
//After this call the current point is p.
//
//Originally cairo_mesh_pattern_move_to.
func (m Mesh) MoveTo(p Point) Mesh {
	x, y := p.c()
	C.cairo_mesh_pattern_move_to(m.p, x, y)
	return m
}

//LineTo adds a line to the current patch from the current point to p.
//
//If there is no current point, this call is equivalent to MoveTo.
//
//After this call the current point is p.
//
//Originally cairo_mesh_pattern_line_to.
func (m Mesh) LineTo(p Point) Mesh {
	x, y := p.c()
	C.cairo_mesh_pattern_line_to(m.p, x, y)
	return m
}

//CurveTo adds a cubic Bézier spline to the current patch,
//from the current point to p2, using p0 and p1 as the control points.
//
//If the current patch has no current point, this method will behave
//as if preceded by a call to MoveTo(p0).
//
//After this call the current point will be p2.
//
//Originally cairo_mesh_pattern_curve_to.
func (m Mesh) CurveTo(p0, p1, p2 Point) Mesh {
	x0, y0 := p0.c()
	x1, y1 := p1.c()
	x2, y2 := p2.c()
	C.cairo_mesh_pattern_curve_to(m.p, x0, y0, x1, y1, x2, y2)
	return m
}

func clampMeshPoint(n int) C.uint {
	return C.uint(n % 4)
}

//SetControlPoint sets the internal contorl point whichPoint
//of the current patch.
//
//Originally cairo_mesh_pattern_set_control_point.
func (m Mesh) SetControlPoint(whichPoint int, p Point) Mesh {
	x, y := p.c()
	wp := clampMeshPoint(whichPoint)
	C.cairo_mesh_pattern_set_control_point(m.p, wp, x, y)
	return m
}

//SetCornerColor sets the color of whichCorner to c in the current patch.
//
//Originally cairo_mesh_pattern_set_corner_color_rgba.
func (m Mesh) SetCornerColor(whichCorner int, c color.Color) Mesh {
	wc := clampMeshPoint(whichCorner)
	r, g, b, a := colorToAlpha(c).c()
	C.cairo_mesh_pattern_set_corner_color_rgba(m.p, wc, r, g, b, a)
	return m
}

//Patches reports the number of patches defined on this mesh.
//
//Originally cairo_mesh_pattern_get_patch_count.
func (m Mesh) Patches() int {
	var pc C.uint
	C.cairo_mesh_pattern_get_patch_count(m.p, &pc)
	return int(pc)
}

//Path gets path defining patch whichPatch.
//
//Originally cairo_mesh_pattern_get_path.
func (m Mesh) Path(whichPatch int) (Path, error) {
	p := C.cairo_mesh_pattern_get_path(m.p, C.uint(whichPatch))
	defer C.cairo_path_destroy(p)
	return cPath(p)
}

//BUG(jmf): cairo_mesh_patten_get_control_point. Do not understand doc:
//"patch_num can range 0 to 1 less than the number returned by cairo_mesh_pattern_get_patch_count"
//so whatever this caveat is, it is undocumented in this binding.

//ControlPoint reports control point whichPoint of patch whichPatch.
//
//Originally cairo_mesh_pattern_get_control_point.
func (m Mesh) ControlPoint(whichPatch, whichPoint int) Point {
	wp := clampMeshPoint(whichPoint)
	var x, y C.double
	C.cairo_mesh_pattern_get_control_point(m.p, C.uint(whichPatch), wp, &x, &y)
	return cPt(x, y)
}

//CornerColor reports the color of corner whichCorner of patch whichPatch.
//
//Originally cairo_mesh_pattern_get_corner_color_rgba.
func (m Mesh) CornerColor(whichPatch, whichCorner int) color.Color {
	pn := C.uint(whichPatch)
	wc := clampMeshPoint(whichCorner)
	var r, g, b, a C.double
	C.cairo_mesh_pattern_get_corner_color_rgba(m.p, pn, wc, &r, &g, &b, &a)
	return cColor(r, g, b, a)
}
