package cairo

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
import "C"

import (
	"image/color"
	"math"
	"runtime"
	"unsafe"
)

//Version returns the version of libcairo.
func Version() string {
	return C.GoString(C.cairo_version_string())
}

//The Context is the main object used when drawing with cairo.
//
//Defaults
//
//The default compositing operator is OpOver.
//
//The default source pattern is equivalent to
//	cairo.NewSolidPattern(color.Black)
//
//The default fill rule is FillRuleWinding.
//
//The default line cap rule is LineCapButt.
//
//The default line join is LineJoinMiter.
//
//The default line width is 2.
//
//The default miter limit is 10.
//
//The default operator is OpOver.
//
//The default tolerance is 0.1.
//
//The default font size is 10.
//
//The default font slant is SlantNormal.
//
//The default font weight is WeightNormal.
//
//The default font family is platform-specific but typically "sans-serif",
//using the toy font api.
//
//Originally cairo_t.
type Context struct {
	c *C.cairo_t
	s Surface
}

//New creates a new drawing context that draws on the target surface.
//
//Originally cairo_create.
func New(target Surface) (*Context, error) {
	s := target.XtensionRaw()
	c := &Context{
		c: C.cairo_create(s),
		s: target,
	}
	runtime.SetFinalizer(c, (*Context).Close)
	return c, c.Err()
}

//Close destroys c.
//
//Originally cairo_destroy.
func (c *Context) Close() error {
	if c == nil || c.c == nil {
		return nil
	}
	runtime.SetFinalizer(c, nil)
	err := c.Err()
	C.cairo_destroy(c.c)
	c.c = nil
	c.s = nil
	return err
}

//Err reports the current error state of c.
//
//Originally cairo_status.
func (c *Context) Err() error {
	if c.c == nil {
		return ErrInvalidLibcairoHandle
	}
	return toerr(C.cairo_status(c.c))
}

//Save makes a copy of the current drawing state on an internal stack
//of drawing states.
//
//Further calls on c will affect the current state but no saved states.
//A call to Restore returns c to the state it was in at the last invocation
//of Save.
//
//Originally cairo_save.
func (c *Context) Save() *Context {
	C.cairo_save(c.c)
	return c
}

//Restore the last drawing state saved with Save and
//remove that state from the stack of saved states.
//
//It is malformed to call Restore without a previous call to Save.
//
//Originally cairo_restore.
func (c *Context) Restore() error {
	C.cairo_restore(c.c)
	return c.Err()
}

//SaveRestore saves the current drawing state, runs f on c, and then restores
//the current drawing state.
func (c *Context) SaveRestore(f func(*Context) error) (err error) {
	c.Save()
	defer func() {
		e := c.Restore()
		if err != nil {
			err = e
		}
	}()
	err = f(c)
	return
}

//Target returns the surface passed to New.
//
//Originally cairo_get_target
func (c *Context) Target() Surface {
	return c.s
}

//PushGroup temporarily redirects drawing to an intermediate surface known
//as a group.
//The redirection lasts until the group is completed by a call to PopGroup
//or PopGroupToSource.
//These calls provide the result of any drawing to the group as a pattern,
//either as an explicit object or set as the source pattern.
//
//This group functionality can be convenient for performing intermediate
//compositing.
//One common use of a group is to render objects as opaque within the group,
//so that they occlude each other, and then blend the result with translucence
//onto the destination.
//
//Groups can be nested arbitrarily deep by making balanced calls to PushGroup
//and PopGroup/PopGroupToSource.
//Each call pushes/pops the new target group onto/from a stack.
//
//Like Save, any changes to the drawing state following PushGroup will not
//be visible after a call to PopGroup/PopGroupToSource.
//
//By default, this intermediate group will have a content type of
//ContentColorAlpha.
//Other content types may be specified by calling PushGroupWithContent instead.
//
//An example of a translucent filled and stroked path without any portion of the
//visible under the stroke:
//
//	c.PushGroup().
//		SetSource(fillPattern).
//		FillPreserve().
//		SetSource(strokePattern).
//		Stroke().
//		PopGroupToSource()
//	c.PaintWithAlpha(alpha)
//
//Originally cairo_push_group.
func (c *Context) PushGroup() *Context {
	C.cairo_push_group(c.c)
	return c
}

//PushGroupWithContent is similar to PushGroup except for additionally setting
//the content of the group.
//
//See the documentation for PushGroup for more information.
//
//Originally cairo_push_group_with_content.
func (c *Context) PushGroupWithContent(content Content) *Context {
	C.cairo_push_group_with_content(c.c, content.c())
	return c
}

//PopGroup terminates the redirection begun by a call to PushGroup
//or PushGroupWithContent and returns a new Pattern.
//
//The drawing state of c is reset to what it was before the call to
//PushGroup/PushGroupWithContent.
//
//See the documentation for PushGroup for more information.
//
//Originally cairo_pop_group.
func (c *Context) PopGroup() (Pattern, error) {
	p := C.cairo_pop_group(c.c)
	if err := c.Err(); err != nil {
		return nil, err
	}
	return cPattern(p)
}

//PopGroupToSource terminates the redirection begun by a call to PushGroup
//or PushGroupWithContent and installs the resultant pattern as the source
//pattern in c.
//
//The drawing state of c is reset to what it was before the call to
//PushGroup/PushGroupWithContent.
//
//It is a convenience method equivalent to
//	func PopGroupToSource(c *Context) error {
//		p, err := c.PopGroup()
//		if err != nil {
//			return err
//		}
//		c.SetSource(p)
//		return p.Close()
//	}
//
//See the documentation for PushGroup for more information.
//
//Originally cairo_pop_group_to_source.
func (c *Context) PopGroupToSource() error {
	C.cairo_pop_group_to_source(c.c)
	return c.Err()
}

//GroupTarget returns the current destination surface for c.
//If there have been no calls to PushGroup/PushGroupWithContent,
//this is equivalent to Target.
//Otherwise,
//Originally cairo_get_group_target.
func (c *Context) GroupTarget() (Surface, error) {
	s := C.cairo_get_group_target(c.c)
	s = C.cairo_surface_reference(s)
	return XtensionRevivifySurface(s)
}

//SetSourceColor sets the source pattern to col.
//This color will be used for any subsequent drawing operations, until a new
//source pattern is set.
//
//Originally cairo_set_source_rgba.
func (c *Context) SetSourceColor(col color.Color) *Context {
	r, g, b, a := colorToAlpha(col).c()
	C.cairo_set_source_rgba(c.c, r, g, b, a)
	return c
}

//SetSource sets the source pattern of c to source.
//This pattern will be used for any subsequent drawing operations, until a new
//source pattern is set.
//
//Note
//
//The pattern's transformation matrix will be locked to the user space in effect
//at the time of SetSource.
//This means that further modifications of the current transformation matrix will
//not affect the source pattern.
//See Pattern.SetMatrix.
//
//Originally cairo_set_source.
func (c *Context) SetSource(source Pattern) *Context {
	C.cairo_set_source(c.c, source.c())
	return c
}

//SetSourceSurface is a convenience function for creating a pattern
//from a surface and setting it as the source pattern in c
//with SetSource.
//
//The originDisplacement vector gives the user-space coordinates
//at which the surface origin should appear.
//The surface origin is its upper-left corner before any transformation
//has been applied.
//The x and y components of the vector are negated and then set
//as the translation values in the pattern matrix.
//
//Other than the initial pattern matrix, described above, all other
//pattern attributes are set to their default values.
//The resulting pattern can be retrieved by calling Source.
//
//Originally cairo_set_source_surface.
func (c *Context) SetSourceSurface(s Surface, originDisplacement Point) error {
	sr := s.XtensionRaw()
	x, y := originDisplacement.c()
	C.cairo_set_source_surface(c.c, sr, x, y)
	return c.Err()
}

//Source returns the current source pattern for c.
//
//Originally cairo_get_source.
func (c *Context) Source() (Pattern, error) {
	p := C.cairo_get_source(c.c)
	p = C.cairo_pattern_reference(p)
	return cPattern(p)
}

//SetAntialiasMode sets the antialiasing mode of the rasterizer used
//for drawing shapes.
//This value is a hint, and a particular backend may or may not support
//a particular value.
//
//At the current time, no backend supports AntialiasSubpixel when drawing shapes.
//
//SetAntialiasMode does not affect text rendering, instead use
//FontOptions.SetAntialiasMode.
//
//Originally cairo_set_antialias.
func (c *Context) SetAntialiasMode(a antialias) *Context {
	C.cairo_set_antialias(c.c, a.c())
	return c
}

//AntialiasMode reports the current shape antialiasing mode.
//
//Originally cairo_get_antialias.
func (c *Context) AntialiasMode() antialias {
	return antialias(C.cairo_get_antialias(c.c))
}

//SetDash sets the dash pattern to be used by Stroke.
//
//A dash pattern is specified by a sequence of positive float64s.
//
//Each float64 represents the length of alternating "on" and "off"
//portions of the stroke.
//
//The offset specifies an offset into the pattern at which the stroke begins.
//
//Each "on" segment will have caps applied as if the segment were a separate
//sub-path.
//It is valid to use an "on" length of 0 with LineCapRound or LineCapSquare
//in order to distribute dots or squares along a path.
//
//Note
//
//The length values are in user-space units as evaluated
//at the time of stroking, which is not necessarily the same as the user space
//at the time SetDash is called.
//
//Special Cases
//
//If the length of dashes is 0, dashing is disabled.
//
//If the length of dashes is 1, a symmetric pattern is assumed,
//where the alternating off and on portions are of the single length provided.
//That is
//	SetDash(0, .5)
//and
//	SetDash(0, .5, .5)
//are equivalent.
//
//Errors
//
//If any of the elements of dashes is negative or all are zero,
//ErrInvalidDash is returned and the dash is not set.
//This differs from libcairo, which puts c into an error mode.
//
//Orginally cairo_set_dash.
func (c *Context) SetDash(offset float64, dashes ...float64) error {
	off := C.double(offset)
	nd := len(dashes)

	arr := make([]C.double, nd)
	allZero := true
	for i, d := range dashes {
		if d < 0 {
			return ErrInvalidDash
		}
		if d > 0 {
			allZero = false
		}
		arr[i] = C.double(d)
	}
	if allZero {
		return ErrInvalidDash
	}

	C.cairo_set_dash(c.c, &arr[0], C.int(nd), off)
	return nil
}

//DashCount reports the length of the dash sequence or 0 if dashing is not
//currently in effect.
//
//Originally cairo_get_dash_count.
func (c *Context) DashCount() int {
	return int(C.cairo_get_dash_count(c.c))
}

//Dashes reports the current dash sequence.
//If dashing is not currently in effect the length of dashes
//is 0.
//
//Originally cairo_get_dash
func (c *Context) Dashes() (offset float64, dashes []float64) {
	ln := c.DashCount()
	if ln == 0 {
		return
	}

	var off C.double
	arr := make([]C.double, ln)
	C.cairo_get_dash(c.c, &arr[0], &off)
	offset = float64(off)

	dashes = make([]float64, ln)
	for i, v := range arr {
		dashes[i] = float64(v)
	}

	return
}

//SetFillRule sets the fill rule on c.
//The fill rule is used to determine which regions are inside or outside
//a complex, potentially self-intersecting, path.
//
//The fill rule affects Fill and Clip.
//
//Originally cairo_set_fill_rule.
func (c *Context) SetFillRule(f fillRule) *Context {
	C.cairo_set_fill_rule(c.c, f.c())
	return c
}

//FillRule reports the current fill rule.
//
//Originally cairo_get_fill_rule.
func (c *Context) FillRule() fillRule {
	return fillRule(C.cairo_get_fill_rule(c.c))
}

//SetLineCap sets the line cap style.
//
//Originally cairo_set_line_cap
func (c *Context) SetLineCap(lc lineCap) *Context {
	C.cairo_set_line_cap(c.c, lc.c())
	return c
}

//LineCap reports the current line cap.
//
//Originally cairo_get_line_cap.
func (c *Context) LineCap() lineCap {
	return lineCap(C.cairo_get_line_cap(c.c))
}

//SetLineJoin sets the line join style.
//
//Originally cairo_set_line_join
func (c *Context) SetLineJoin(l lineJoin) *Context {
	C.cairo_set_line_join(c.c, l.c())
	return c
}

//LineJoin reports the current line join style.
//
//Originally cairo_get_line_join.
func (c *Context) LineJoin() lineJoin {
	return lineJoin(C.cairo_get_line_join(c.c))
}

//SetLineWidth sets the current line width.
//The line width specifies the diameter of a pen that is circular
//in user space, though the device space pen may be an ellipse due
//to shearing in the coordinate transform matrix.
//The user space and coordinate transform matrix referred to above are computed
//at stroke time, not at the time SetLineWidth is called.
//
//Originally cairo_set_line_width.
func (c *Context) SetLineWidth(width float64) *Context {
	C.cairo_set_line_width(c.c, C.double(width))
	return c
}

//LineWidth reports the line width as set by SetLineWidth and does not
//take any intervening changes to the coordinate transform matrix into account.
//
//Originally cairo_get_line_width
func (c *Context) LineWidth() float64 {
	return float64(C.cairo_get_line_width(c.c))
}

//SetMiterLimit sets the miter limit.
//
//When the current line join style is LineJoinMiter, the miter limit is used to
//determine whether the lines should be joined with a bevel instead of a miter.
//
//Cairo divides the length of the miter by the line width.
//If the result is greater than the miter limit, the style is converted to a
//bevel.
//
//A miter limit for a given angle can be computed by:
//	miter limit = 1/sin(angle/2)
//
//Examples
//
//For the default of 10, joins with interior angles less than 11 degrees are
//converted from miters to bevels.
//
//For reference, a mite limit of 2 makes the miter cutoff at 60 degrees,
//and a miter limit of 1.414 makes the cutoff at 90 degrees.
//
//Originally cairo_set_miter_limit.
func (c *Context) SetMiterLimit(ml float64) *Context {
	C.cairo_set_miter_limit(c.c, C.double(ml))
	return c
}

//MiterLimit returns the current miter limit as set by SetMiterLimit.
//
//Originally cairo_get_miter_limit.
func (c *Context) MiterLimit() float64 {
	return float64(C.cairo_get_miter_limit(c.c))
}

//SetOperator sets the compositing operator used for all drawing operations.
//
//Originally cairo_set_operator.
func (c *Context) SetOperator(op operator) *Context {
	C.cairo_set_operator(c.c, op.c())
	return c
}

//Operator reports the current compositing operator.
//
//Originally cairo_get_operator.
func (c *Context) Operator() operator {
	return operator(C.cairo_get_operator(c.c))
}

//SetTolerance sets the tolerance, in device units, when converting paths into
//trapezoids.
//Curved segments of the path will be subdivided until the maximum deviation
//between the original path and the polygonal approximation is less than
//tolerance.
//
//A larger value than the default of 0.1 will give better performance.
//While in general a lower value improves appearance, it is unlikely a value
//lower than .1 will improve appearance significantly.
//
//The accuracy of paths within libcairo is limited by the precision of its
//internal arithmetic and tolerance is restricted by the smallest representable
//internal value.
//
//Originally cairo_set_tolerance.
func (c *Context) SetTolerance(tolerance float64) *Context {
	C.cairo_set_tolerance(c.c, C.double(tolerance))
	return c
}

//Tolerance reports the tolerance in device units.
//
//Originally cairo_get_tolerance.
func (c *Context) Tolerance() float64 {
	return float64(C.cairo_get_tolerance(c.c))
}

//Clip establishes a new clip region by intersecting the current clip region
//with the current path as it would be filled by Fill and according to the
//current fill rule.
//
//After Clip, the current path will be cleared from the cairo context.
//
//The current clip region affects all drawing operations by effectively masking
//out any changes to the surface that are outside the current clip region.
//
//Clip can only make the clip region smaller, never larger.
//But the current clip is part of the graphics state, so a temporary restriction
//of the clip region can be achieved by calling Clip within a Save/Restore pair.
//The only other means of increasing the size of the clip region is ResetClip.
//
//Originally cairo_clip.
func (c *Context) Clip() *Context {
	C.cairo_clip(c.c)
	return c
}

//ClipPreserve is identical Clip but preserves the path in c.
//
//Originally cairo_clip_preserve.
func (c *Context) ClipPreserve() *Context {
	C.cairo_clip_preserve(c.c)
	return c
}

//ClipExtents computes a bounding box in user coordinates covering the area
//inside the current clip.
//
//Originally cairo_clip_extents.
func (c *Context) ClipExtents() Rectangle {
	var x0, y0, x1, y1 C.double
	C.cairo_clip_extents(c.c, &x0, &y0, &x1, &y1)
	return cRect(x0, y0, x1, y1)
}

//InClip reports whether pt is in the currently visible area defined by the
//clipping region.
//
//Originally cairo_in_clip.
func (c *Context) InClip(pt Point) bool {
	x, y := pt.c()
	return C.cairo_in_clip(c.c, x, y) == 1
}

//ResetClip resets the current clip region to its original, unrestricted state.
//That is, set the clip region to an infinitely large shape containing
//the target surface.
//Equivalently, one can imagine the clip region being reset to the exact bounds
//of the target surface.
//
//Note that code meant to be reusable should not call ResetClip as it will
//cause results unexpected by higher-level code which calls ResetClip.
//Consider using Save and Restore around Clip as a more robust means of
//temporarily restricting the clip region.
//
//Originally cairo_reset_clip.
func (c *Context) ResetClip() *Context {
	C.cairo_reset_clip(c.c)
	return c
}

//ClipRectangles reports the current clip region as a list of rectangles
//in user coordinates or an error if the clip region cannot be so represented.
//
//Originally cairo_copy_clip_rectangle_list.
func (c *Context) ClipRectangles() (list []Rectangle, err error) {
	rects := C.cairo_copy_clip_rectangle_list(c.c)
	if err := toerr(rects.status); err != nil {
		return nil, err
	}

	n := int(rects.num_rectangles)
	if n == 0 {
		return nil, nil
	}
	rs := (*[1 << 30]C.cairo_rectangle_t)(unsafe.Pointer(rects.rectangles))[:n:n]
	list = make([]Rectangle, n)
	for i, v := range rs {
		list[i] = cRect(v.x, v.y, v.x+v.width, v.y+v.height)
	}

	C.cairo_rectangle_list_destroy(rects)

	return
}

//Fill fills the current path according to the current fill rule,
//(each sub-path is implicitly closed before being filled).
//After fill, the current path will be cleared from the cairo context.
//
//Originally cairo_fill.
func (c *Context) Fill() *Context {
	C.cairo_fill(c.c)
	return c
}

//FillPreserve is identical to Fill except it does not clear the current path.
//
//Originally cairo_fill_preserve.
func (c *Context) FillPreserve() *Context {
	C.cairo_fill_preserve(c.c)
	return c
}

//FillExtents computes a bounding box in user coordinates covering the area that
//would be affected, (the "inked" area), by a Fill operation given the current
//path and fill parameters.
//If the current path is empty, it returns ZR.
//Surface dimensions and clipping are not taken into account.
//
//Contrast with PathExtents, which is similar, but returns non-zero extents
//for some paths with no inked area, (such as a simple line segment).
//
//FillExtents must necessarily do more work to compute the precise inked areas
//in light of the fill rule, so PathExtents may be more desirable for sake of
//performance if the non-inked path extents are desired.
//
//Originally cairo_fill_extents.
func (c *Context) FillExtents() Rectangle {
	var x1, y1, x2, y2 C.double
	C.cairo_fill_extents(c.c, &x1, &y1, &x2, &y2)
	return cRect(x1, y1, x2, y2)
}

//InFill reports whether the given point is inside the area that would be
//affected by a Fill, given the current path and filling parameters.
//Surface dimensions and clipping are not taken into account.
//
//Originally cairo_in_fill.
func (c *Context) InFill(pt Point) bool {
	x, y := pt.c()
	return C.cairo_in_fill(c.c, x, y) == 1
}

//Mask paints the current source using the alpha channel of pattern as a mask.
//Opaque areas of pattern are painted with the source, transparent areas are
//not painted.
//
//Originally cairo_mask.
func (c *Context) Mask(p Pattern) *Context {
	C.cairo_mask(c.c, p.c())
	return c
}

//MaskSurface paints the current source using the alpha channel of surface
//as a mask.
//Opaque areas of surface are painted with the source, transparent areas are
//not painted.
//
//Originally cairo_mask_surface.
func (c *Context) MaskSurface(s Surface, offsetVector Point) *Context {
	x, y := offsetVector.c()
	C.cairo_mask_surface(c.c, s.XtensionRaw(), x, y)
	return c
}

//Paint paints the current source everywhere within the current clip region.
//
//originally cairo_paint.
func (c *Context) Paint() *Context {
	C.cairo_paint(c.c)
	return c
}

//PaintAlpha paints the current source everywhere within the current clip
//region, using a mask of constant alpha value alpha.
//
//The effect is similar to Paint, but the drawing is faded out using
//the alpha value.
//
//originally cairo_paint.
func (c *Context) PaintAlpha(alpha float64) *Context {
	C.cairo_paint_with_alpha(c.c, C.double(alpha))
	return c
}

//Stroke strokes the current path according to the current line width,
//line join, line cap, and dash settings.
//
//After Stroke, the current path will be cleared from the cairo context.
//
//Degenerate segments and sub-paths are treated specially and provide a useful
//result.
//These can result in two different situations:
//
//1. Zero-length "on" segments set in SetDash.
//If the cap style is LineCapRound or LineCapSquare then these segments will be
//drawn as circular dots or squares respectively.
//In the case of LineCapSquare, the orientation of the squares is determined by
//the direction of the underlying path.
//
//2. A sub-path created by MoveTo followed by either a ClosePath or one or more
//calls to LineTo to the same coordinate as the MoveTo.
//If the cap style is LineCapRound then these sub-paths will be drawn as
//circular dots.
//Note that in the case of LineCapSquare a degenerate sub-path will not be drawn
//at all, as the correct orientation is indeterminate.
//
//In no case will a cap style of LineCapButt cause anything to be drawn in the
//case of either degenerate segments or sub-paths.
//
//Originally cairo_stroke.
func (c *Context) Stroke() *Context {
	C.cairo_stroke(c.c)
	return c
}

//StrokePreserve is identical to Stroke except the path is not cleared.
//
//Originally cairo_stroke_preserve.
func (c *Context) StrokePreserve() *Context {
	C.cairo_stroke_preserve(c.c)
	return c
}

//StrokeExtents computes a bounding box in user coordinates covering the area
//that would be affected, (the "inked" area), by a Stroke operation given the
//current path and stroke parameters.
//If the current path is empty, returns the empty rectangle ZR.
//Surface dimensions and clipping are not taken into account.
//
//Note that StrokeExtents must necessarily do more work to compute the precise
//inked areas in light of the stroke parameters, so PathExtents may be more
//desirable for sake of performance if non-inked path extents are desired.
//
//Originally cairo_stroke_extents.
func (c *Context) StrokeExtents() Rectangle {
	var x1, y1, x2, y2 C.double
	C.cairo_stroke_extents(c.c, &x1, &y1, &x2, &y2)
	return cRect(x1, y1, x2, y2)
}

//InStroke reports whether the given point is inside the area that would be
//affected by a Stroke operation given the current path and stroking
//parameters.
//Surface dimensions and clipping are not taken into account.
//
//Originally cairo_in_stroke.
func (c *Context) InStroke(pt Point) bool {
	x, y := pt.c()
	return C.cairo_in_stroke(c.c, x, y) == 1
}

//CopyPage emits the current page for backends that support multiple pages,
//but doesn't clear it, so, the contents of the current page will be retained
//for the next page too.
//
//Use ShowPage if you want to get an empty page after the emission.
//
//This is a convenience function that simply calls CopyPage on c's target.
//
//Originally cairo_copy_page.
func (c *Context) CopyPage() *Context {
	C.cairo_copy_page(c.c)
	return c
}

//ShowPage emits and clears the current page for backends that support multiple
//pages.
//
//Use CopyPage if you don't want to clear the page.
//
//This is a convenience function that simply calls ShowPage on c's target.
//
//Originally cairo_show_page.
func (c *Context) ShowPage() *Context {
	C.cairo_show_page(c.c)
	return c
}

//CopyPath returns a copy of the current path.
//
//Originally cairo_copy_path.
func (c *Context) CopyPath() (Path, error) {
	p := C.cairo_copy_path(c.c)
	defer C.cairo_path_destroy(p)
	return cPath(p)
}

//CopyPathFlat returns a linearized copy of the current path.
//
//CopyPathFlat behaves like CopyPath except that any curves in the path will be
//approximated with piecewise-linear approximations, accurate to within the
//current tolerance value.
//That is, the result is guaranteed to not have any elements of type PathCurveTo
//which will instead be replaced by a series of PathLineTo elements.
//
//Originally cairo_copy_path_flat.
func (c *Context) CopyPathFlat() (Path, error) {
	p := C.cairo_copy_path_flat(c.c)
	defer C.cairo_path_destroy(p)
	return cPath(p)
}

//AppendPath appends path onto the current path of c.
//
//Originally cairo_append_path.
func (c *Context) AppendPath(path Path) error {
	p, err := path.c()
	if err != nil {
		return err
	}
	C.cairo_append_path(c.c, p)
	C.cairo_path_destroy(p) //BUG(jmf): does cairo take control of path after append path?
	return c.Err()
}

//CurrentPoint reports the current point of the current path.
//The current point is, conceptually, the final point reached by the path
//so far.
//
//The current point is returned in the user-space coordinate system.
//
//If there is no defined current point, or if c is in an error state,
//(ZP, false) will be returned.
//Otherwise the (cp, true) will be returned where cp is the current point.
//
//Most path constructions alter the current point.
//
//Some functions use and alter the current point, but do not otherwise change
//the current path, see ShowText.
//
//Some functions unset the current path, and, as a result, the current point,
//such as Fill.
//
//Originally cairo_has_current_point and cairo_get_current_point.
func (c *Context) CurrentPoint() (cp Point, defined bool) {
	has := C.cairo_has_current_point(c.c) == 1
	if !has {
		return
	}
	var x, y C.double
	C.cairo_get_current_point(c.c, &x, &y)
	return cPt(x, y), true
}

//NewPath clears the current path and, by extension, the current point.
//
//Originally cairo_new_path.
func (c *Context) NewPath() *Context {
	C.cairo_new_path(c.c)
	return c
}

//NewSubPath begins a new sub-path.
//The existing path is not affected, but the current point is cleared.
//
//In many cases, this is not needed since new sub-paths are frequently started
//with MoveTo.
//
//NewSubPath is particularly useful when beginning a new sub-path with one of the
//Arc calls, as, in this case, it is no longer necessary to manually computer the
//arc's inital coordinates for use with MoveTo.
//
//Originally cairo_new_sub_path.
func (c *Context) NewSubPath() *Context {
	C.cairo_new_sub_path(c.c)
	return c
}

//ClosePath adds a line segment from the current point to the beginning
//of the current sub-path and closes the sub-path.
//After this call the current point will be at the joined endpoint of the
//sub-path.
//
//The behavior of ClosePath is distinct from simply calling LineTo with the
//equivalent coordinate in the case of stroking.
//When a closed sub-path is stroked, there are no caps on the ends of the
//sub-path.
//Instead, there is a line join connecting the final and initial segments
//of the sub-path.
//
//If there is no current point, this method will have no effect.
//
//ClosePath will place an explicit PathMoveTo following the PathClosePath into
//the current path.
func (c *Context) ClosePath() *Context {
	C.cairo_close_path(c.c)
	return c
}

//Arc adds a circular arc along the surface of circle from fromAngle
//increasing to toAngle.
//
//If fromAngle < toAngle, then toAngle will be increased by 2π until
//fromAngle > toAngle.
//
//If there is a current point, an initial line segment will be added
//to the path to connect the current point to the beginning of the arc.
//If this initial line is undesired, call ClosePath before Arc.
//
//Angles are measured in radians.
//An angle of 0 is in the direction of the positive X axis in user space.
//An angle of π/2 radians (90°) is in the direction of the positive Y axis
//in user space.
//With the default transformation matrix, angles increase clockwise.
//
//To convert from degrees to radians use
//	degrees * π/180
//
//Arc gives the arc in the direction of increasing angles.
//Use ArcNegative to get the arc in the direction of decreasing
//angles.
//
//The arc is circular in user space.
//
//Originally cairo_arc.
func (c *Context) Arc(circle Circle, fromAngle, toAngle float64) *Context {
	x, y, r := circle.c()
	a1, a2 := C.double(fromAngle), C.double(toAngle)
	C.cairo_arc(c.c, x, y, r, a1, a2)
	return c
}

//Circle is shorthand for calling Arc from 0 to 2π.
func (c *Context) Circle(circle Circle) *Context {
	return c.Arc(circle, 0, 2*math.Pi)
}

//ArcNegative adds a circular arc along the surface of circle from fromAngle
//decreasing to toAngle.
//
//If fromAngle > toAngle, then toAngle will be dereased by 2π until
//fromAngle < toAngle.
//
//ArcNegative gives the arc in the direction of decreasing angles.
//Use Arc to get the arc in the direction of increasing angles.
//
//Originally cairo_arc_negative.
func (c *Context) ArcNegative(circle Circle, fromAngle, toAngle float64) *Context {
	x, y, r := circle.c()
	a1, a2 := C.double(fromAngle), C.double(toAngle)
	C.cairo_arc_negative(c.c, x, y, r, a1, a2)
	return c
}

//CurveTo adds a cubic Bézier spline to the path from the current point to p3
//in user space coordinates, using p1 and p2 as the control points.
//After calling CurveTo, the current point will be p3.
//
//If there is no current point, CurveTo will behave as if preceded by a call
//to MoveTo(p1)
//
//Originally cairo_curve_to.
func (c *Context) CurveTo(p1, p2, p3 Point) *Context {
	x0, y0 := p1.c()
	x1, y1 := p2.c()
	x2, y2 := p3.c()
	C.cairo_curve_to(c.c, x0, y0, x1, y1, x2, y2)
	return c
}

//LineTo adds a line to the path from the current point to p in user space
//coordinates.
//
//If there is no current point, LineTo will behave as if preceded by a call
//to MoveTo(p)
//
//Originally cairo_line_to.
func (c *Context) LineTo(p Point) *Context {
	x, y := p.c()
	C.cairo_line_to(c.c, x, y)
	return c
}

//MoveTo begins a new sub-path and sets the current point to p.
//
//Originally cairo_move_to.
func (c *Context) MoveTo(p Point) *Context {
	x, y := p.c()
	C.cairo_move_to(c.c, x, y)
	return c
}

//Rectangle adds a closed sub-path rectangle to the current path at position
//r.Min in user-space coordinates.
//
//This function is logically equivalent to:
//	c.MoveTo(r.Min)
//	c.LineTo(Pt(r.Dx(), 0))
//	c.LineTo(Pt(0, r.Dy()))
//	c.LineTo(Pt(-r.Dx(), 0))
//	c.ClosePath()
//
//Originally cairo_rectangle.
func (c *Context) Rectangle(r Rectangle) *Context {
	x, y, w, h := r.cWH()
	C.cairo_rectangle(c.c, x, y, w, h)
	return c
}

//RelCurveTo is a relative-coordinate version of CurveTo.
//All points are considered as vectors with an origin at the current point.
//
//It is equivalent to
//	if p, ok := c.CurrentPoint(); ok {
//		c.CurveTo(p.Add(v1), p.Add(v2), p.Add(v3))
//	} else {
//		// c is broken now
//	}
//
//Originally cairo_rel_curve_to.
func (c *Context) RelCurveTo(v1, v2, v3 Point) error {
	x0, y0 := v1.c()
	x1, y1 := v2.c()
	x2, y2 := v3.c()
	C.cairo_rel_curve_to(c.c, x0, y0, x1, y1, x2, y2)
	return c.Err()
}

//RelLineTo is a relative-coordinate version of LineTo.
//The point v is considered a vector with the origin at the current point.
//
//It is equivalent to
//	if p, ok := c.CurrentPoint(); ok {
//		c.LineTo(p.Add(v))
//	} else {
//		// c is broken now
//	}
//
//Originally cairo_rel_line_to.
func (c *Context) RelLineTo(v Point) error {
	x, y := v.c()
	C.cairo_rel_line_to(c.c, x, y)
	return c.Err()
}

//RelMoveTo begins a new sub-path.
//After this call the current point will be offset by v.
//
//It is equivalent to
//	if p, ok := c.CurrentPoint(); ok {
//		c.Move(p.Add(v))
//	} else {
//		// c is broken now
//	}
//
//Originally cairo_rel_move_to.
func (c *Context) RelMoveTo(v Point) error {
	x, y := v.c()
	C.cairo_rel_move_to(c.c, x, y)
	return c.Err()
}

//PathExtents computes a bounding box in user-space coordinates covering
//the points on the current path.
//If the current path is empty, returns ZR.
//Stroke parameters, fill rule, surface dimensions and clipping are not taken
//into account.
//
//PathExtents is in contrast to FillExtents and StrokeExtents which return
//the extents of only the area that would be "inked" by the corresponding
//drawing operations.
//
//The result of PathExtents is defined as equivalent to the limit
//of StrokeExtents with LineCapRound as the line width approaches 0, but never
//approaching the empty rectangle returned by StrokeExtents for a line width
//of 0.
//
//Specifically, this means that zero-area sub-paths such as MoveTo contribute
//to the extents.
//However, a lone MoveTo will not contribute to the results of PathExtents.
//
//Originally cairo_path_extents.
func (c *Context) PathExtents() Rectangle {
	var x, y, x1, y1 C.double
	C.cairo_path_extents(c.c, &x, &y, &x1, &y1)
	return cRect(x, y, x1, y1)
}

//GlyphPath adds a closed path for the glyphs to the current path.
//The generated path, if filled, achieves an effect similar to that
//of ShowGlyphs.
//
//Originally cairo_glyph_path.
func (c *Context) GlyphPath(glyphs []Glyph) *Context {
	gs, n := glyphsC(glyphs)
	C.cairo_glyph_path(c.c, gs, n)
	return c
}

//TextPath adds closed paths for text to the current path.
//The generated path if filled, achieves an effect similar to that of ShowText.
//
//Text conversion and positioning is done similar to ShowText.
//
//Like ShowText, After this call the current point is moved to the origin
//of where the next glyph would be placed in this same progression.
//hat is, the current point will be at the origin of the final glyph offset
//by its advance values.
//This allows for chaining multiple calls to to TextPath without having to set
//current point in between.
//
//Note: The TextPath method is part of what the libcairo designers call
//the "toy" text API.
//It is convenient for short demos and simple programs, but it is not expected
//to be adequate for serious text-using applications.
//See GlyphPath for the "real" text path in cairo.
//
//Originally cairo_text_path.
func (c *Context) TextPath(s string) *Context {
	cs := C.CString(s)
	C.cairo_text_path(c.c, cs)
	C.free(unsafe.Pointer(cs))
	return c
}

//SelectFont selects a family and style of font from a simplified
//description as a family name, slant and weight.
//
//Libcairo provides no operation to list available family names
//on the system, as this is part of the toy api.
//The standard CSS2 generic family names, ("serif", "sans-serif", "cursive",
//"fantasy", "monospace"), are likely to work as expected.
//
//If family starts with the string "cairo:", or if no native font backends are
//compiled in, libcairo will use an internal font family.
//
//If text is drawn without a call to SelectFont or SetFont
//or SetScaledFont, the platform-specific default family is used, typically
//"sans-serif", and the default slant is SlantNormal and the default weight
//is WeightNormal.
//
//For "real" see the font-backend-specific factor functions for the font
//backend you are using.
//The resulting font face could then be used with ScaledFontCreate and
//SetScaledFont.
//
//Note
//
//SelectFont is part of the "toy" text API.
//
//Originally cairo_select_font_face.
func (c *Context) SelectFont(family string, slant slant, weight weight) *Context {
	f := C.CString(family)
	C.cairo_select_font_face(c.c, f, slant.c(), weight.c())
	C.free(unsafe.Pointer(f))
	return c
}

//SetFontSize sets the current font matrix to a scale by a factor of size,
//replacing any font matrix previously set with SetFontSize or SetFontMatrix.
//
//Originally cairo_set_font_size.
func (c *Context) SetFontSize(size float64) *Context {
	C.cairo_set_font_size(c.c, C.double(size))
	return c
}

//SetFontMatrix sets the current font matrix to m.
//The font matrix gives a transformation from the design space of the font
//(in this space, the em-square is 1 unit by 1 unit) to user space.
//
//Originally cairo_set_font_matrix.
func (c *Context) SetFontMatrix(m Matrix) *Context {
	C.cairo_set_font_matrix(c.c, &m.m)
	return c
}

//FontMatrix reports the current font matrix.
//
//Originally cairo_get_font_matrix.
func (c *Context) FontMatrix() Matrix {
	var m C.cairo_matrix_t
	C.cairo_get_font_matrix(c.c, &m)
	return Matrix{m}
}

//SetFontOptions sets the custom font rendering options for c.
//Rendering operations are derived by merging opts with the
//FontOptions of the underlying surface:
//if any of the options in opts is the default, the value of
//the surface is used.
//
//Originally cairo_set_font_options.
func (c *Context) SetFontOptions(opts *FontOptions) *Context {
	C.cairo_set_font_options(c.c, opts.fo)
	return c
}

//FontOptions reports the current font options of c.
//
//Originally cairo_get_font_options.
func (c *Context) FontOptions() *FontOptions {
	f := NewFontOptions()
	C.cairo_get_font_options(c.c, f.fo)
	return f
}

//SetFont replaces the current font face of c with f.
//
//Originally cairo_set_font_face.
func (c *Context) SetFont(f Font) *Context {
	C.cairo_set_font_face(c.c, f.XtensionRaw())
	return c
}

//Font returns the current font face for c.
//
//Originally cairo_get_font_face.
func (c *Context) Font() (Font, error) {
	fr := C.cairo_get_font_face(c.c)
	fr = C.cairo_font_face_reference(fr)
	return cFont(fr)
}

//SetScaledFont replaces the current font face, font matrix, and font options
//with those of sf.
//Except for some translation, the current coordinate transforms matrix of c
//should be the same as that of sf.
//
//Originally cairo_set_scaled_font.
func (c *Context) SetScaledFont(sf *ScaledFont) *Context {
	C.cairo_set_scaled_font(c.c, sf.f)
	return c
}

//ScaledFont reports the scaled font of c.
//
//Originally cairo_get_scaled_font.
func (c *Context) ScaledFont() (*ScaledFont, error) {
	f := C.cairo_get_scaled_font(c.c)
	f = C.cairo_scaled_font_reference(f)
	sf := cNewScaledFont(f)
	return sf, c.Err()
}

//ShowText draws a shape generated from s, rendered according to the current
//Font, font matrix, and FontOptions.
//
//This method first computes a set of glyphs for the string of text.
//The first glyph is placed so that its origin is at the current point.
//The origin of each subsequent glyph is offset from that of the previous glyph
//by the advance values of the previous glyph.
//
//After this call the current point is moved to the origin of where the next
//glyph would be placed in this same progression.
//That is, the current point will be at the origin of the final glyph offset
//by its advance values.
//This allows for easy display of a single logical string with multiple calls
//to ShowText.
//
//Note
//
//ShowText is part of the toy api.
//See ShowGlyphs for the "real" text display api.
//
//Originally cairo_show_text.
func (c *Context) ShowText(s string) *Context { //BUG(jmf): what if there is no current point before calling ShowText?
	cs := C.CString(s)
	C.cairo_show_text(c.c, cs)
	C.free(unsafe.Pointer(cs))
	return c
}

//ShowGlyphs draws a shape generated from gylphs rendered according
//to the current Font, font matrix, and FontOptions.
//
//Originally cairo_show_glyphs.
func (c *Context) ShowGlyphs(glyphs []Glyph) *Context {
	gs, n := glyphsC(glyphs)
	C.cairo_show_glyphs(c.c, gs, n)
	return c
}

//ShowTextGlyphs renders similarly to ShowGlyphs but, if the target surface
//support it, uses the provided text and cluster mapping to embed the text
//for the glyphs shown in the output.
//If the target surface does not support the extended attributes, this method
//behaves exactly as ShowGlyphs(glyphs).
//
//The mapping between s and glyphs is provided by clusters.
//Each cluster covers a number of text bytes and glyphs, and neighboring
//clusters cover neighboring areas of s and glyphs.
//The clusters should collectively cover s and glyphs in entirety.
//
//The first cluster always covers bytes from the beginning of s.
//If flags do not have TextClusterBackward set, the first cluster also covers
//the beginning of glyphs, otherwise it covers the end of the glyphs array and
//following clusters move backward.
//
//Originally cairo_show_text_glyphs.
func (c *Context) ShowTextGlyphs(s string, glyphs []Glyph, clusters []TextCluster, flags textClusterFlags) *Context {
	gs, gn := glyphsC(glyphs)
	ts, tn := clustersC(clusters)
	cs, cn := C.CString(s), C.int(len(s))
	C.cairo_show_text_glyphs(c.c, cs, cn, gs, gn, ts, tn, flags.c())
	C.free(unsafe.Pointer(cs))
	return c
}

//FontExtents returns the extents of the currently selected font.
//
//Originally cairo_font_extents.
func (c *Context) FontExtents() FontExtents {
	var f C.cairo_font_extents_t
	C.cairo_font_extents(c.c, &f)
	return newFontExtents(f)
}

//TextExtents reports the extents for s.
//The extents describe a user-space rectangle that encloses the "inked" portion
//of the text, as it would be drawn with ShownText.
//Additionally, the x_advance and y_advance values indicate the amount by which
//the current point would be advanced by ShowText.
//
//Note that whitespace characters do not directly contribute to the size
//of the rectangle (extents.Width and extents.Height), but they do contribute
//indirectly by changing the position of non-whitespace characters.
//In particular, trailing whitespace characters are likely to not affect
//the size of the rectangle, though they will affect the AdvanceX and AdvanceY
//values.
//
//Originally cairo_text_extents.
func (c *Context) TextExtents(s string) TextExtents {
	var t C.cairo_text_extents_t
	cs := C.CString(s)
	C.cairo_text_extents(c.c, cs, &t)
	C.free(unsafe.Pointer(cs))
	return newTextExtents(t)
}

//GlyphExtents reports the extents for glyphs.
//
//The extents describe a user-space rectangle that encloses the "inked" portion
//of the text, as it would be drawn with ShownText.
//Additionally, the x_advance and y_advance values indicate the amount by which
//the current point would be advanced by ShowText.
//
//Note that whitespace characters do not directly contribute to the size
//of the rectangle (extents.Width and extents.Height), but they do contribute
//indirectly by changing the position of non-whitespace characters.
//In particular, trailing whitespace characters are likely to not affect
//the size of the rectangle, though they will affect the AdvanceX and AdvanceY
//values.
//
//Originally cairo_glyph_extents.
func (c *Context) GlyphExtents(glyphs []Glyph) TextExtents {
	var t C.cairo_text_extents_t
	gs, n := glyphsC(glyphs)
	C.cairo_glyph_extents(c.c, gs, n, &t)
	return newTextExtents(t)
}

//Translate the current transformation matrix by vector v.
//This offset is interpreted as a user-space coordinate according to the CTM
//in place before the new call to Translate.
//In other words, the translation of the user-space origin takes place after
//any existing transformation.
//
//Originally cairo_translate.
func (c *Context) Translate(v Point) *Context {
	x, y := v.c()
	C.cairo_translate(c.c, x, y)
	return c
}

//Scale scales the current transformation matrix by v by scaling the user-space
//axes by v.X and v.Y.
//The scaling of the axes takes place after any existing transformation
//of user space.
//
//Originally cairo_scale.
func (c *Context) Scale(v Point) *Context {
	x, y := v.c()
	C.cairo_scale(c.c, x, y)
	return c
}

//Rotate the current transformation matrix by θ by rotating the user-space
//axes.
//The rotation of the axes takes places after any existing transformation
//of user space.
//The rotation direction for positive angles is from the positive X axis toward
//the positive Y axis.
//
//Originally cairo_rotate.
func (c *Context) Rotate(θ float64) *Context {
	C.cairo_rotate(c.c, C.double(θ))
	return c
}

//Transform applies m to the current transformation matrix as an additional
//transformation.
//The new transformation of user space takes place after any existing
//transformation.
//
//Originally cairo_transform.
func (c *Context) Transform(m Matrix) *Context {
	C.cairo_transform(c.c, &m.m)
	return c
}

//SetMatrix sets the current transformation matrix to m.
//
//Originally cairo_set_matrix.
func (c *Context) SetMatrix(m Matrix) *Context {
	C.cairo_set_matrix(c.c, &m.m)
	return c
}

//Matrix returns the current transformation matrix.
//
//Originally cairo_get_matrix.
func (c *Context) Matrix() Matrix {
	var m C.cairo_matrix_t
	C.cairo_get_matrix(c.c, &m)
	return Matrix{m}
}

//ResetMatrix resets the current transformation matrix to the identity
//matrix.
//
//Originally cairo_identity_matrix.
func (c *Context) ResetMatrix() *Context {
	C.cairo_identity_matrix(c.c)
	return c
}

//UserToDevice takes the point p from user space to the point q in device space
//by multiplication with the current transformation matrix.
//
//Originally cairo_user_to_device.
func (c *Context) UserToDevice(p Point) (q Point) {
	x, y := p.c()
	C.cairo_user_to_device(c.c, &x, &y)
	return cPt(x, y)
}

//UserToDeviceDistance transforms a distance vector v from user to device
//space.
//This method is similar to UserToDevice, except that the translation
//components of the current transformation matrix will be ignored.
//
//Originally cairo_user_to_device_distance.
func (c *Context) UserToDeviceDistance(v Point) Point {
	x, y := v.c()
	C.cairo_user_to_device_distance(c.c, &x, &y)
	return cPt(x, y)
}

//DeviceToUser takes the point p from device space to the point q in user space
//by multiplication with the inverse of the current transformation matrix.
//
//Originally cairo_device_to_user.
func (c *Context) DeviceToUser(p Point) (q Point) {
	x, y := p.c()
	C.cairo_device_to_user(c.c, &x, &y)
	return cPt(x, y)
}

//DeviceToUserDistance transforms a distance vector v from device to user
//space.
//This method is similar to DeviceToUser, except that the translation
//components of the current transformation matrix will be ignored.
//
//Originally cairo_device_to_user_distance.
func (c *Context) DeviceToUserDistance(v Point) Point {
	x, y := v.c()
	C.cairo_device_to_user_distance(c.c, &x, &y)
	return cPt(x, y)
}
