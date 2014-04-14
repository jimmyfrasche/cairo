//Examples come from http://cairographics.org/samples and as such these
//examples are in the public domain and were originally contributed
//by Øyvind Kolås.
//
//In this directory, run
//	go run examples.go
//and each item in the examples slice will output a pdf of that name.
package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/jimmyfrasche/cairo"
	"github.com/jimmyfrasche/cairo/pdf"
)

var img image.Image

func getImage() (image.Image, error) {
	f, err := os.Open("romedalen.png")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	im, _, err := image.Decode(f)
	return im, err
}

var pt = cairo.Pt

func deg2rad(rad float64) (deg float64) {
	const f = math.Pi / 180.
	return rad * f
}

var examples = []struct {
	Name string
	Run  func(*cairo.Context) error
}{
	{"arc", func(c *cairo.Context) error {
		p := cairo.Circ(128, 128, 100)

		ang1, ang2 := deg2rad(45), deg2rad(180)

		c.
			SetLineWidth(10).
			Arc(p, ang1, ang2).
			Stroke()

		c.NewPath()

		//draw helping lines
		c.
			SetSourceColor(cairo.AlphaColor{1, .2, .2, .6}).
			SetLineWidth(6)

		c.
			Circle(p.Mul(1 / 10.)).
			Fill()

		c.
			Arc(p, ang1, ang1).
			LineTo(p.Center)

		c.
			Arc(p, ang2, ang2).
			LineTo(p.Center)

		c.Stroke()

		return nil
	}},
	{"arc-negative", func(c *cairo.Context) error {
		p := cairo.Circ(128, 128, 100)

		ang1, ang2 := deg2rad(45), deg2rad(180)

		c.
			SetLineWidth(10).
			ArcNegative(p, ang1, ang2).
			Stroke()

		c.NewPath()

		//draw helping lines
		c.
			SetSourceColor(cairo.AlphaColor{1, .2, .2, .6}).
			SetLineWidth(6)

		c.
			Arc(cairo.Circle{p.Center, 10}, 0, 2*math.Pi).
			Fill()

		c.
			Arc(p, ang1, ang1).
			LineTo(p.Center)

		c.
			Arc(p, ang2, ang2).
			LineTo(p.Center)

		c.Stroke()

		return nil
	}},
	{"clip", func(c *cairo.Context) error {
		c.
			Circle(cairo.Circ(128, 128, 76.8)).
			Clip()

		//Current path is not consumed by Clip.
		c.NewPath()

		c.
			Rectangle(cairo.Rect(0, 0, 256, 256)).
			Fill()

		c.
			SetSourceColor(cairo.Color{G: 1}).
			MoveTo(cairo.ZP).
			LineTo(pt(256, 256)).
			MoveTo(pt(256, 0)).
			LineTo(pt(0, 256)).
			SetLineWidth(10).
			Stroke()

		return nil
	}},
	{"clip-image", func(c *cairo.Context) error {
		i, err := cairo.FromImage(img) //img declared globally and set in main
		if err != nil {
			return err
		}
		defer i.Close()

		c.
			Circle(cairo.Circ(128, 128, 76.8)).
			Clip().
			NewPath()

		sz := i.Size()
		c.Scale(pt(256/sz.X, 256/sz.Y))
		if err = c.SetSourceSurface(i, cairo.ZP); err != nil {
			return err
		}
		c.Paint()

		return nil
	}},
	{"curve-rectangle", func(c *cairo.Context) error {
		curveRect := func(c *cairo.Context, r cairo.Rectangle, radius float64) {
			if r.Dx() == 0 || r.Dy() == 0 {
				return
			}

			lt := r.Min
			lb := pt(r.Min.X, r.Max.Y)
			rt := pt(r.Max.X, r.Min.Y)
			rb := r.Max

			ya := (lt.Y + rb.Y) / 2
			xa := (lt.X + rb.Y) / 2

			xab := pt(lt.X, ya)
			xat := pt(rb.X, ya)
			yab := pt(xa, lt.Y)
			yat := pt(xa, rb.Y)

			xr := pt(radius, 0)
			yr := pt(0, radius)

			ltayr := lt.Add(yr)
			lbsyr := lb.Sub(yr)
			ltaxr := lt.Add(xr)
			lbaxr := lb.Add(xr)
			rtayr := rt.Add(yr)
			rbsyr := rb.Sub(yr)
			rtsxr := rt.Sub(xr)
			rbsxr := rb.Sub(xr)

			if r.Dx()/2 < radius {
				if r.Dy()/2 < radius {
					c.
						MoveTo(xab).
						CurveTo(lt, lt, yab).
						CurveTo(rt, rt, xat).
						CurveTo(rb, rb, yat).
						CurveTo(lb, lb, xab)
				} else {
					c.
						MoveTo(ltayr).
						CurveTo(lt, lt, yab).
						CurveTo(rt, rt, rtayr).
						LineTo(rbsyr).
						CurveTo(rb, rb, yat).
						CurveTo(lb, lb, lbsyr)
				}
			} else {
				if r.Dy()/2 < radius {
					c.
						MoveTo(xab).
						CurveTo(lt, lt, ltaxr).
						LineTo(rtsxr).
						CurveTo(rt, rt, xat).
						CurveTo(rb, rb, rbsxr).
						LineTo(lbaxr).
						CurveTo(lb, lb, xab)
				} else {
					c.
						MoveTo(ltayr).
						CurveTo(lt, lt, ltaxr).
						LineTo(rtsxr).
						CurveTo(rt, rt, rtayr).
						LineTo(rbsyr).
						CurveTo(rb, rb, rbsxr).
						LineTo(lbaxr).
						CurveTo(lb, lb, ltayr)
				}
			}

			c.ClosePath()
			c.SetSourceColor(cairo.Color{.5, .5, 1})
			c.FillPreserve()
			c.SetSourceColor(cairo.AlphaColor{R: .5, A: .5})
			c.SetLineWidth(10)
			c.Stroke()
		}

		curveRect(c, cairo.RectWH(25.6, 25.6, 204.8, 204.8), 102.4)
		return nil
	}},
	{"curve-to", func(c *cairo.Context) error {
		ps := pt(25.6, 128)
		p1 := pt(102.4, 230.4)
		p2 := pt(153.6, 25.6)
		p3 := pt(230.4, 128)

		c.
			MoveTo(ps).
			CurveTo(p1, p2, p3).
			SetLineWidth(10).
			Stroke()

		c.
			SetSourceColor(cairo.AlphaColor{1, .2, .2, .6}).
			SetLineWidth(6).
			MoveTo(ps).LineTo(p1).
			MoveTo(p2).LineTo(p3).
			Stroke()

		return nil
	}},
	{"dash", func(c *cairo.Context) error {
		//NB: the offset comes first since we're using ...args
		if err := c.SetDash(-50, 50, 10, 10, 10); err != nil {
			return err
		}

		c.SetLineWidth(10)

		c.
			MoveTo(pt(128, 25.6)).
			LineTo(pt(230.4, 230.4))

		//This can only err if c is err'd or no current point
		if err := c.RelLineTo(pt(-102.4, 0)); err != nil {
			return err
		}

		c.CurveTo(pt(51.2, 230.4), pt(51.2, 128), pt(128, 128))

		c.Stroke()

		return nil
	}},
	{"fill-and-stroke2", func(c *cairo.Context) error {
		//pregnant triangle
		c.
			MoveTo(pt(128, 25.6)).
			LineTo(pt(230.4, 230.4))
		if err := c.RelLineTo(pt(-102.4, 0)); err != nil {
			return err
		}
		c.
			CurveTo(pt(51.2, 230.4), pt(51.2, 128), pt(128, 128)).
			ClosePath()

		//diamond
		c.MoveTo(pt(64, 25.6))
		p := pt(51.2, 51.2)
		if err := c.RelLineTo(p); err != nil {
			return err
		}
		if err := c.RelLineTo(p.Rx()); err != nil {
			return err
		}
		if err := c.RelLineTo(p.Conj()); err != nil {
			return err
		}
		c.ClosePath()

		//fill in blue
		c.
			SetLineWidth(10).
			SetSourceColor(cairo.Blue).
			FillPreserve()

		//stroke black
		c.
			SetSourceColor(cairo.Black).
			Stroke()

		return nil
	}},
	{"fill-style", func(c *cairo.Context) error {
		c.SetLineWidth(6)

		c1, c2 := cairo.Circ(64, 64, 40), cairo.Circ(192, 64, 40)
		r := cairo.RectWH(12, 12, 232, 70)

		c.
			Rectangle(r).
			NewSubPath().Arc(c1, 0, 2*math.Pi).
			NewSubPath().ArcNegative(c2, 0, -2*math.Pi)

		c.
			SetFillRule(cairo.FillRuleEvenOdd).
			SetSourceColor(cairo.Color{G: .7}).
			FillPreserve().
			SetSourceColor(cairo.Black).
			Stroke()

		c.
			Translate(pt(0, 128)).
			Rectangle(r).
			NewSubPath().Arc(c1, 0, 2*math.Pi).
			NewSubPath().ArcNegative(c2, 0, -2*math.Pi)

		c.
			SetFillRule(cairo.FillRuleWinding).
			SetSourceColor(cairo.Color{B: .9}).
			FillPreserve().
			SetSourceColor(cairo.Black).
			Stroke()

		return nil
	}},
	{"gradient", func(c *cairo.Context) error {
		lg := cairo.NewLinearGradient(cairo.ZP, pt(0, 256), []cairo.ColorStop{
			{1, cairo.Black},
			{0, cairo.White},
		}...)
		defer lg.Close()

		c.
			Rectangle(cairo.Rectangle{Max: pt(256, 256)}).
			SetSource(lg).
			Fill()

		rg := cairo.NewRadialGradient(
			cairo.Circ(115.2, 102.4, 25.6),
			cairo.Circ(102.4, 102.4, 128),
			[]cairo.ColorStop{
				{0, cairo.White},
				{1, cairo.Black},
			}...)
		defer rg.Close()
		c.
			SetSource(rg).
			Circle(cairo.Circ(128, 128, 76)).
			Fill()
		return nil
	}},
	{"image", func(c *cairo.Context) error {
		i, err := cairo.FromImage(img) //img declared globally and set in main
		if err != nil {
			return err
		}
		defer i.Close()

		sz := i.Size()
		c.
			Translate(pt(128, 128)).
			Rotate(deg2rad(45)).
			Scale(pt(256/sz.X, 256/sz.Y)).
			Translate(sz.Mul(-.5))

		if err = c.SetSourceSurface(i, cairo.ZP); err != nil {
			return err
		}
		c.Paint()

		return nil
	}},
	{"image-pattern", func(c *cairo.Context) error {
		i, err := cairo.FromImage(img) //img declared globally and set in main
		if err != nil {
			return err
		}
		defer i.Close()
		sz := i.Size()

		pat, err := cairo.NewSurfacePattern(i)
		if err != nil {
			return err
		}
		defer pat.Close()
		pat.SetExtend(cairo.ExtendRepeat)

		off := pt(128, 128)
		c.
			Translate(off).
			Rotate(math.Pi / 4).
			Scale(pt(1/math.Sqrt2, 1/math.Sqrt2)).
			Translate(off.Conj())

		pat.SetMatrix(cairo.NewScaleMatrix(sz.Div(256).Mul(5)))

		c.
			SetSource(pat).
			Rectangle(cairo.Rect(0, 0, 256, 256)).
			Fill()
		return nil
	}},
	{"multi-segment-caps", func(c *cairo.Context) error {
		line := func(c *cairo.Context, y float64) {
			c.MoveTo(pt(50, y)).LineTo(pt(200, y))
		}
		line(c, 75)
		line(c, 125)
		line(c, 175)
		c.
			SetLineWidth(30).
			SetLineCap(cairo.LineCapRound).
			Stroke()
		return nil
	}},
	{"rounded-rectangle", func(c *cairo.Context) error {
		roundedRect := func(c *cairo.Context, r cairo.Rectangle, aspect, radius float64) {
			radius = radius / aspect

			circ := func(x, y float64) cairo.Circle {
				return cairo.Circle{r.Min.Add(pt(x, y)), radius}
			}
			w, h := r.Dx(), r.Dy()

			c.
				NewSubPath().
				Arc(circ(w-radius, radius), deg2rad(-90), 0).
				Arc(circ(w-radius, h-radius), 0, deg2rad(90)).
				Arc(circ(radius, h-radius), deg2rad(90), deg2rad(180)).
				Arc(circ(radius, radius), deg2rad(180), deg2rad(270)).
				ClosePath().
				SetSourceColor(cairo.Color{.5, .5, 1}).
				FillPreserve().
				SetSourceColor(cairo.AlphaColor{R: .5, A: .5}).
				SetLineWidth(10).
				Stroke()
		}

		roundedRect(c, cairo.RectWH(25.6, 25.6, 204.8, 204.8), 1, 204.8/10)
		return nil
	}},
	{"set-line-cap", func(c *cairo.Context) error {
		c.SetLineWidth(30)

		line := func(x float64, c *cairo.Context) {
			c.
				MoveTo(pt(x, 50)).
				LineTo(pt(x, 200)).
				Stroke()
		}

		line(64, c.SetLineCap(cairo.LineCapButt)) //Default line cap
		line(128, c.SetLineCap(cairo.LineCapRound))
		line(192, c.SetLineCap(cairo.LineCapSquare))

		//draw helper lines
		c.
			SetSourceColor(cairo.Color{1, .2, .2}).
			SetLineWidth(2.56)
		line(64, c)
		line(128, c)
		line(192, c)
		return nil
	}},
	{"set-line-join", func(c *cairo.Context) error {
		angle := func(y float64, c *cairo.Context) (err error) {
			c.MoveTo(pt(76.8, y))
			if err = c.RelLineTo(pt(51.2, -51.2)); err != nil {
				return
			}
			if err = c.RelLineTo(pt(51.2, 51.2)); err != nil {
				return
			}
			c.Stroke()
			return nil
		}

		c.SetLineWidth(40.96)

		if err := angle(84.48, c.SetLineJoin(cairo.LineJoinMiter)); err != nil { //default
			return err
		}
		if err := angle(161.28, c.SetLineJoin(cairo.LineJoinBevel)); err != nil {
			return err
		}
		return angle(238.08, c.SetLineJoin(cairo.LineJoinRound))
	}},
	{"text", func(c *cairo.Context) error {
		c.
			SelectFont("Sans", 0, cairo.WeightBold).
			SetFontSize(90)

		c.
			MoveTo(pt(10, 135)).
			ShowText("Hello")

		c.
			MoveTo(pt(70, 165)).
			TextPath("void").
			SetSourceColor(cairo.Color{.5, .5, 1}).
			FillPreserve().
			SetSourceColor(cairo.Black).
			SetLineWidth(2.56).
			Stroke()

		//draw helper dots
		c.
			SetSourceColor(cairo.AlphaColor{1, .2, .2, .6}).
			Circle(cairo.Circ(10, 135, 5.12)).
			ClosePath().
			Circle(cairo.Circ(70, 165, 5.12)).
			Fill()
		return nil
	}},
	{"text-align-center", func(c *cairo.Context) error {
		const str = "cairo"
		c.
			SelectFont("Sans", 0, 0).
			SetFontSize(52)

		te := c.TextExtents(str)
		cp := pt(128-(te.Width/2+te.BearingX), 128-(te.Height/2+te.BearingY))

		c.
			MoveTo(cp).
			ShowText(str)

		//draw helping lines
		c.
			SetSourceColor(cairo.AlphaColor{1, .2, .2, .6}).
			SetLineWidth(6).
			Circle(cairo.Circle{cp, 10}).
			Fill()
		if err := c.MoveTo(pt(128, 0)).RelLineTo(pt(0, 256)); err != nil {
			return err
		}
		if err := c.MoveTo(pt(0, 128)).RelLineTo(pt(256, 0)); err != nil {
			return err
		}
		c.Stroke()

		return nil
	}},
	{"text-extents", func(c *cairo.Context) error {
		const s = "cairo"

		c.
			SelectFont("Sans", 0, 0).
			SetFontSize(100)

		te := c.TextExtents(s)

		start := pt(25, 150)
		c.
			MoveTo(start).
			ShowText(s)

		//draw helping lines
		c.
			SetSourceColor(cairo.AlphaColor{1, .2, .2, .6}).
			SetLineWidth(6).
			Circle(cairo.Circle{start, 10}).
			Fill().
			MoveTo(start)
		if err := c.RelLineTo(pt(0, -te.Height)); err != nil {
			return err
		}
		if err := c.RelLineTo(pt(te.Width, 0)); err != nil {
			return err
		}
		if err := c.RelLineTo(pt(te.BearingX, -te.BearingY)); err != nil {
			return err
		}
		c.Stroke()
		return nil
	}},
}

func main() {
	log.SetFlags(0)

	var err error
	img, err = getImage()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Using libcairo:", cairo.Version())

	count := 0
	logerr := func(name string, err error) {
		count++
		log.Printf("%s: %s", name, err)
	}

	for _, example := range examples {
		nm := example.Name

		outname := nm + ".pdf"

		outfile, err := os.Create(outname)
		if err != nil {
			logerr(nm, err)
			continue
		}
		defer outfile.Close()

		surface, err := pdf.New(outfile, 595, 842) //A4
		if err != nil {
			logerr(nm, err)
			continue
		}
		defer surface.Close()

		context, err := cairo.New(surface)
		if err != nil {
			logerr(nm, err)
			continue
		}
		defer context.Close()

		if err = example.Run(context); err != nil {
			logerr(nm, err)
			continue
		}

		if err = context.Err(); err != nil {
			logerr(nm, err)
		}
	}

	fmt.Printf("%d of %d examples failed\n", count, len(examples))
}
