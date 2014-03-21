package cairo

//#cgo pkg-config: cairo
//#include <cairo/cairo.h>
import "C"

//Matrix is used throughout cairo to convert between different coordinate
//spaces.
//
//A Matrix holds an affine transformation, such as a scale, rotation, shear,
//or a combination of these.
//
//The transformation of a point (x,y) is given by:
//	x' = xx*x + xy*y + x0
//	y' = yx*y + yy*y + y0
//
//Originally cairo_matrix_t.
type Matrix struct {
	m C.cairo_matrix_t
}

//NewMatrix creates the affine transformation matrix given by
//xx, yx, xy, yy, x0, y0.
//
//Originally cairo_matrix_init
func NewMatrix(xx, yx, xy, yy, x0, y0 float64) Matrix {
	var m C.cairo_matrix_t
	C.cairo_matrix_init(&m, C.double(xx), C.double(yx), C.double(xy), C.double(yy), C.double(x0), C.double(y0))
	return Matrix{m}
}

//NewIdentityMatrix creates the identity I.
//
//Originally cairo_init_identity.
func NewIdentityMatrix() (I Matrix) {
	var m C.cairo_matrix_t
	C.cairo_matrix_init_identity(&m)
	return Matrix{m}
}

//NewTranslateMatrix matrix creates a matrix that translates by vector.
//
//Originally cairo_init_translate.
func NewTranslateMatrix(vector Point) Matrix {
	var m C.cairo_matrix_t
	C.cairo_matrix_init_translate(&m, C.double(vector.X), C.double(vector.Y))
	return Matrix{m}
}

//NewScaleMatrix creates a matrix that scales by vector.
//
//Originally cairo_matrix_init_scale.
func NewScaleMatrix(vector Point) Matrix {
	var m C.cairo_matrix_t
	C.cairo_matrix_init_scale(&m, C.double(vector.X), C.double(vector.Y))
	return Matrix{m}
}

//NewRotateMatrix creates a matrix that rotates by radians.
//
//Originally cairo_matrix_init_rotate.
func NewRotateMatrix(radians float64) Matrix {
	var m C.cairo_matrix_t
	C.cairo_matrix_init_rotate(&m, C.double(radians))
	return Matrix{m}
}

//Clone returns a new matrix that is the same as m.
func (m Matrix) Clone() Matrix {
	var n C.cairo_matrix_t
	C.cairo_matrix_init(&n, m.m.xx, m.m.yx, m.m.xy, m.m.yy, m.m.x0, m.m.y0)
	return Matrix{n}
}

//XX returns the XX component of the matrix.
func (m Matrix) XX() float64 {
	return float64(m.m.xx)
}

//XY returns the XY component of the matrix.
func (m Matrix) XY() float64 {
	return float64(m.m.xy)
}

//YX returns the YX component of the matrix.
func (m Matrix) YX() float64 {
	return float64(m.m.yx)
}

//YY returns the YY component of the matrix.
func (m Matrix) YY() float64 {
	return float64(m.m.yy)
}

//X0 returns the X0 component of the matrix.
func (m Matrix) X0() float64 {
	return float64(m.m.x0)
}

//Y0 returns the Y0 component of the matrix.
func (m Matrix) Y0() float64 {
	return float64(m.m.y0)
}

//Translate translates m by vector and returns itself.
//
//Originally cairo_matrix_translate.
func (m Matrix) Translate(vector Point) Matrix {
	C.cairo_matrix_translate(&m.m, C.double(vector.X), C.double(vector.Y))
	return m
}

//Scale scales m by vector and returns itself.
//
//Originally cairo_matrix_scale.
func (m Matrix) Scale(vector Point) Matrix {
	C.cairo_matrix_scale(&m.m, C.double(vector.X), C.double(vector.Y))
	return m
}

//Rotate rotates m by radians and returns itself.
//
//Originally cairo_matrix_rotate.
func (m Matrix) Rotate(radians float64) Matrix {
	C.cairo_matrix_rotate(&m.m, C.double(radians))
	return m
}

//Invert inverts A and returns itself.
//
//Originally cairo_matrix_invert.
func (m Matrix) Invert() Matrix {
	C.cairo_matrix_invert(&m.m)
	return m
}

//Mul multiples m by n and returns the new result r, such that r = m*n.
//
//Originally cairo_matrix_multiply.
func (m Matrix) Mul(n Matrix) Matrix {
	var r Matrix
	C.cairo_matrix_multiply(&r.m, &m.m, &n.m)
	return r
}

//TransformDistance transforms the distance vector p by m.
//This is similar to Transform except that the translation component of m
//are ignored.
//
//Originally cairo_matrix_transform_distance.
func (p Point) TransformDistance(m Matrix) Point {
	x := C.double(p.X)
	y := C.double(p.Y)
	C.cairo_matrix_transform_distance(&m.m, &x, &y)
	return Point{float64(x), float64(y)}
}

//Transform transforms p by m, returning a new point.
//
//Originally cairo_matrix_transform.
func (p Point) Transform(m Matrix) Point {
	x := C.double(p.X)
	y := C.double(p.Y)
	C.cairo_matrix_transform_point(&m.m, &x, &y)
	return Point{float64(x), float64(y)}
}
