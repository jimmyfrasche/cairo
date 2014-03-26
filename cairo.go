package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

import (
	"image/color"
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
//Originally cairo_t.
type Context struct {
	c *C.cairo_t
	s Surface
}

//New creates a new drawing context that draws on the target surface.
//
//Originally cairo_create.
func New(target Surface) (*Context, error) {
	s := target.ExtensionRaw()
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
func (c *Context) PushGroupWithContent(content content) *Context {
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
	//XXX should we track group depth and return c.Target(), nil if 0?
	s := C.cairo_get_group_target(c.c)
	s = C.cairo_surface_reference(s)
	return cSurface(s)
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
	sr := s.ExtensionRaw()
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
func (c *Context) SetOperator(op op) *Context {
	C.cairo_set_operator(c.c, op.c())
	return c
}

//Operator reports the current compositing operator.
//
//Originally cairo_get_operator.
func (c *Context) Operator() op {
	return op(C.cairo_get_operator(c.c))
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
	C.cairo_mask_surface(c.c, s.ExtensionRaw(), x, y)
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
