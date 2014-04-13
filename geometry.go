package cairo

import "C"

import (
	"math"
	"strconv"
)

//Point is an X, Y coordinate pair.
//The axes increase right and down.
//
//When a Point is used as a vector, it is considered as a line segment from (0, 0) to
//(X, Y).
type Point struct {
	X, Y float64
}

func (p Point) c() (x, y C.double) {
	return C.double(p.X), C.double(p.Y)
}

//Pt is shorthand for Point{X, Y}.
func Pt(X, Y float64) Point {
	return Point{X, Y}
}

func cPt(x, y C.double) Point {
	return Point{float64(x), float64(y)}
}

//Polar converts polar coordinates to cartesian.
func Polar(r, θ float64) Point {
	sinθ, cosθ := math.Sincos(θ)
	return Pt(r*cosθ, r*sinθ)
}

//ZP is the zero point.
var ZP Point

func floatstr(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func (p Point) String() string {
	return "(" + floatstr(p.X) + "," + floatstr(p.Y) + ")"
}

//Conj returns (-x, -y) the conjugate of (x, y).
func (p Point) Conj() Point {
	return Pt(-p.X, -p.Y)
}

//Rx reflects p about the x-axis.
func (p Point) Rx() Point {
	return Pt(-p.X, p.Y)
}

//Ry reflects p about the y-axis.
func (p Point) Ry() Point {
	return Pt(p.X, -p.Y)
}

//Sx shifts p along the x-axis by x.
func (p Point) Sx(x float64) Point {
	return Pt(p.X+x, p.Y)
}

//Sy shifts p along the y-axis by y.
func (p Point) Sy(y float64) Point {
	return Pt(p.X, p.Y+y)
}

//Add returns the vector p-q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

//Sub returns the vector p+q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

//Mul returns the vector k*p.
func (p Point) Mul(k float64) Point {
	return Point{p.X * k, p.Y * k}
}

//Dot returns the dot product of p and q.
func (p Point) Dot(q Point) float64 {
	return p.X*q.X + p.Y*q.Y
}

//Div returns the vector p/k.
func (p Point) Div(k float64) Point {
	return Point{p.X / k, p.Y / k}
}

//Eq reports whether p and q are equal.
func (p Point) Eq(q Point) bool {
	return p.X == q.X && p.Y == q.Y
}

//Near reports whether p and q are within ε of each other.
func (p Point) Near(q Point, ε float64) bool {
	return math.Abs(p.X-q.X) < ε && math.Abs(p.Y-q.Y) < ε
}

//Mag returns the length of the vector Sqrt(p.X*p.X + p.Y+p.Y).
func (p Point) Mag() float64 {
	return math.Hypot(p.X, p.Y)
}

//Norm returns the unit-length vector with the same direction as p.
func (p Point) Norm() (n Point) {
	return p.Div(p.Mag())
}

//Angle returns the angle of the vector in radians.
func (p Point) Angle() float64 {
	return math.Atan2(p.Y, p.X)
}

//In reports whether p is in r.
func (p Point) In(r Rectangle) bool {
	return r.Min.X <= p.X &&
		p.X < r.Max.X &&
		r.Min.Y <= p.Y &&
		p.Y < r.Max.Y
}

//InCirc reports whether p falls in c.
func (p Point) InCirc(c Circle) bool {
	return p.Near(c.Center, c.Radius)
}

//Mod returns the point q in r such that p.X-q.X is a multiple
//of r's width and p.Y-q.Y is a multiple of r's height.
func (p Point) Mod(r Rectangle) Point {
	w, h := r.Dx(), r.Dy()
	p = p.Sub(r.Min)
	p.X = math.Mod(p.X, w)
	if p.X < 0 {
		p.X += w
	}
	p.Y = math.Mod(p.Y, h)
	if p.Y < 0 {
		p.Y += h
	}
	return p.Add(r.Min)
}

//A Rectangle contains the points with Min.X <= X < Max.X,
//Min.Y <= Y < Max.Y.
//A Rectangle is always axis-aligned.
//It is well-formed if Min.X <= Max.X and likewise for Y.
//Points are always well-formed.
//A rectangle's methods always return well-formed outputs
//for well-formed inputs.
type Rectangle struct {
	Min, Max Point
}

func (r Rectangle) c() (x0, y0, x1, y1 C.double) {
	x0, y0 = r.Min.c()
	x1, y1 = r.Max.c()
	return
}

//ZR is the zero Rectangle.
var ZR Rectangle

//Rect is shorthand for Rectangle{Pt(x₀, y₀), Pt(x₁, y₁)}.Canon().
func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Pt(x0, y0), Pt(x1, y1)}
}

//RectWH is shorthand for Rectangle{Pt(x, y), Pt(x+width, y+height)}.Canon().
func RectWH(x, y, width, height float64) Rectangle {
	return Rect(x, y, x+width, y+height)
}

func cRect(x0, y0, x1, y1 C.double) Rectangle {
	return Rectangle{cPt(x0, y0), cPt(x1, y1)}.Canon()
}

//like c but returns Min (Dx, Dy)
func (r Rectangle) cWH() (x, y, w, h C.double) {
	x, y = r.Min.c()
	w, h = r.Size().c()
	return
}

func (r Rectangle) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

//Dx returns r's width.
func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

//Dy returns r's height.
func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

//Size returns r's width and height.
func (r Rectangle) Size() Point {
	return Pt(r.Dx(), r.Dy())
}

//Verts returns all four corners of the rectangle, clockwise from r.Min.
func (r Rectangle) Verts() (x0y0, x0y1, x1y1, x1y0 Point) {
	return r.Min, Pt(r.Min.X, r.Max.Y), r.Max, Pt(r.Max.X, r.Min.Y)
}

//Add returns the rectangle r translated by p.
func (r Rectangle) Add(p Point) Rectangle {
	return Rectangle{
		r.Min.Add(p),
		r.Max.Add(p),
	}
}

//Sub returns the rectangle r translated by -p.
func (r Rectangle) Sub(p Point) Rectangle {
	return r.Add(p.Conj())
}

//Intersect returns the largest rectangle contained by both r and s.
//If the two rectangles do not overlap then the zero rectangle
//will be returned.
func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.Y {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Min.X > r.Max.X || r.Min.Y > r.Max.Y {
		return ZR
	}
	return r
}

//Empty reports whether the rectangle contains no points.
func (r Rectangle) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

//Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle) Overlaps(s Rectangle) bool {
	return r.Min.X > s.Max.X &&
		s.Min.X > r.Max.X &&
		r.Min.Y < s.Max.Y &&
		s.Min.Y < r.Max.Y
}

//In reports whether every point in r is in s.
func (r Rectangle) In(s Rectangle) bool {
	if r.Empty() {
		return true
	}
	return s.Min.X <= r.Min.X &&
		r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y &&
		r.Max.Y <= s.Max.Y
}

//Canon returns the canonical version of r.
//The returned rectangle has minimum and maximum coordinates swapped
//if necessary so that it is well-formed.
func (r Rectangle) Canon() Rectangle {
	if r.Max.X < r.Min.X {
		r.Min.X, r.Max.X = r.Max.X, r.Min.X
	}
	if r.Max.Y < r.Min.Y {
		r.Min.Y, r.Max.Y = r.Max.Y, r.Min.Y
	}
	return r
}

//A Circle contains the points swept out by Radius from Center.
//
//It is well-formed if the Radius is nonnegative.
type Circle struct {
	Center Point
	Radius float64
}

var (
	//ZC is the zero circle.
	ZC Circle
	//UC is the unit circle.
	UC = Circle{Radius: 1}
)

//Circ is shorthand for Circle{Pt(x, y), r}.Canon().
func Circ(x, y, r float64) Circle {
	return Circle{Pt(x, y), r}.Canon()
}

func cCirc(x, y, r C.double) Circle {
	return Circle{cPt(x, y), float64(r)}.Canon()
}

func (c Circle) c() (x, y, r C.double) {
	x, y = c.Center.c()
	r = C.double(c.Radius)
	return
}

func (c Circle) String() string {
	return c.Center.String() + "-" + floatstr(c.Radius)
}

//Canon returns a canonical circle.
func (c Circle) Canon() Circle {
	return Circle{c.Center, math.Abs(c.Radius)}
}

//Add returns the circle c translated by p.
func (c Circle) Add(p Point) Circle {
	return Circle{c.Center.Add(p), c.Radius}
}

//Sub returns the circle c translated by -p.
func (c Circle) Sub(p Point) Circle {
	return c.Add(p.Conj())
}

//Mul returns the circle c with its radius multiplied by k.
func (c Circle) Mul(k float64) Circle {
	return Circle{c.Center, k * c.Radius}
}

//Empty reports whether this circle contains no points.
func (c Circle) Empty() bool {
	return c.Radius == 0
}

//BUG(jmf): finish copying image.Point/Rectangle interfaces over to float
//and document. Just need Inset.

//BUG(jmf): bring circle to feature parity with rectangle
