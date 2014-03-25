package cairo

//#cgo pkg-config: cairo
//#include <stdlib.h>
//#include <cairo/cairo.h>
//
//static void gocairo_get_header(cairo_path_data_t* head, int i, cairo_path_data_type_t* typ, int* length) {
//	*typ = head[i].header.type;
//	*length = head[i].header.length;
//}
//
//static void gocairo_get_point(cairo_path_data_t* head, int i, double* x, double* y) {
//	*x = head[i].point.x;
//	*y = head[i].point.y;
//}
//
//static void gocairo_set_header(cairo_path_data_t* head, int i, cairo_path_data_type_t typ, int length) {
//	head[i].header.type = typ;
//	head[i].header.length = length;
//}
//
//static void gocairo_set_point(cairo_path_data_t* head, int i, double x, double y) {
//	head[i].point.x = x;
//	head[i].point.y = y;
//}
//
//static cairo_path_t* gocairo_new_path(int numelms) {
//	cairo_path_t* path = malloc(sizeof(cairo_path_t));
//	path->status = CAIRO_STATUS_SUCCESS;
//	path->data = malloc(numelms * sizeof(cairo_path_data_t));
//	path->num_data = numelms;
//	return path;
//}
import "C"

func cPath(p *C.cairo_path_t) (path Path, err error) {
	if p == nil {
		err = ErrInvalidPathData
		return
	}
	if err = toerr(p.status); err != nil {
		return
	}

	//We assume the path is correctly formed here as it can only come directly from libcairo.

	N := int(p.num_data)
	for i := 0; i < N; {
		//grab header
		var t C.cairo_path_data_type_t
		var length C.int
		C.gocairo_get_header(p.data, C.int(i), &t, &length)
		typ := pathDataType(t)
		ln := int(length)

		//get points, if any
		var pts []Point
		for j := i; j < i+ln; j++ {
			var x, y C.double
			C.gocairo_get_point(p.data, C.int(j), &x, &y)

			pts = append(pts, cPt(x, y))
		}

		//handle case of extra info stuffed into path.
		if len(pts) > 3 {
			pts = pts[:3:3]
		}
		if typ == PathClosePath && len(pts) != 0 {
			pts = nil
		}
		if (typ == PathMoveTo || typ == PathLineTo) && len(pts) != 1 {
			pts = pts[:1:1]
		}

		pe := pathElement{
			dtype:  pathDataType(t),
			points: pts,
		}
		if !pe.valid() {
			return nil, ErrInvalidPathData
		}

		//add path element
		path = append(path, pe)

		//advance index
		i += ln
	}

	return
}

//A PathElement is a single item in a Path.
type PathElement interface {
	//Type reports the type of this path element.
	Type() pathDataType
	//Points returns a copy of this path element's points.
	//Its length will always be 0, 1, or 3.
	Points() []Point
	//String renders a human readable interpretation of this path element.
	String() string
	valid() bool
	len() int
	pts() []Point
}

type pathElement struct {
	dtype  pathDataType
	points []Point //len = 0, 1, or 3
}

func (p pathElement) valid() bool {
	switch ln := len(p.points); p.dtype {
	default:
		return false
	case PathClosePath:
		return ln == 0
	case PathMoveTo, PathLineTo:
		return ln == 1
	case PathCurveTo:
		return ln == 3
	}
}

func (p pathElement) String() string {
	s := func(i int) string {
		return p.points[i].String()
	}
	switch p.dtype {
	case PathClosePath:
		return "close path"
	case PathMoveTo:
		return "move to " + s(0)
	case PathLineTo:
		return "line to " + s(0)
	case PathCurveTo:
		return "curve to " + s(0) + "-" + s(1) + "-" + s(2)
	}
	return "invalid path element"
}

func (p pathElement) Type() pathDataType {
	return p.dtype
}

func (p pathElement) Points() (out []Point) {
	out = make([]Point, 0, len(p.points))
	copy(out, p.points)
	return
}
func (p pathElement) pts() []Point {
	return p.points
}

func (p pathElement) len() int {
	switch p.dtype {
	case PathClosePath:
		return 0
	case PathMoveTo, PathLineTo:
		return 1
	case PathCurveTo:
		return 3
	}
	return -1
}

type Path []PathElement

func (p *Path) append(t pathDataType, ps ...Point) {
	*p = append(*p, pathElement{
		dtype:  t,
		points: ps,
	})
}

func (p *Path) MoveTo(pt Point) {
	p.append(PathMoveTo, pt)
}

func (p *Path) LineTo(pt Point) {
	p.append(PathLineTo, pt)
}

func (p *Path) CurveTo(p0, p1, p2 Point) {
	p.append(PathCurveTo, p0, p1, p2)
}

func (p *Path) ClosePath() {
	p.append(PathClosePath)
}

func (p Path) sizeMult() (sz int) {
	for _, pe := range p {
		sz += 1 + pe.len()
	}
	return
}

func (p Path) valid() bool {
	for _, pe := range p {
		if !pe.valid() {
			return false
		}
	}
	return true
}

//Note doc comments for cairo_path_destroy says it can't be called on paths created
//outside libcairo but, at least as of 2014.03.24, this does not seem to be true.
//If it causes problems, the approach will have to be rethought.
func (p Path) c() (path *C.cairo_path_t, err error) {
	if !p.valid() {
		return nil, ErrInvalidPathData
	}

	sz := p.sizeMult()
	path = C.gocairo_new_path(C.int(sz))

	i := C.int(0)
	for _, v := range p {
		//add header
		typ := C.cairo_path_data_type_t(v.Type())
		C.gocairo_set_header(path.data, i, typ, C.int(v.len()))
		i++
		//add points
		for _, pt := range v.pts() {
			x, y := pt.c()
			C.gocairo_set_point(path.data, i, x, y)
			i++
		}
	}

	return
}
