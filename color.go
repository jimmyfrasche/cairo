package cairo

import "C"

import (
	"image/color"
)

func clamp01(f float64) float64 {
	switch {
	case f < 0:
		return 0
	case f > 1:
		return 1
	}
	return f
}

//convert f from [0,1] to uint32.
func ctoi(f float64) uint32 {
	f = clamp01(f)
	return uint32(f*float64(1<<32-1) + .5)
}

//convert i to [0,1].
func cto01(i uint32) float64 {
	return float64(i) / float64(1<<32-1)
}

//Color represents an RGB color where each component is in [0, 1].
type Color struct {
	R, G, B float64
}

//Canon returns a new color with all values clamped to [0,1].
func (c Color) Canon() Color {
	return Color{clamp01(c.R), clamp01(c.G), clamp01(c.B)}
}

func (co Color) c() (r, g, b C.double) {
	return C.double(co.R), C.double(co.G), C.double(co.G)
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return ctoi(c.R), ctoi(c.G), ctoi(c.B), 0xffff
}

//Alpha returns c with an alpha channel set to alpha.
func (c Color) Alpha(alpha float64) AlphaColor {
	return AlphaColor{c.R, c.G, c.B, alpha}
}

//AlphaColor represents an RGBA color where each component is in [0,1].
type AlphaColor struct {
	R, G, B, A float64
}

func cColor(r, g, b, a C.double) AlphaColor {
	return AlphaColor{
		float64(r),
		float64(g),
		float64(b),
		float64(a),
	}
}

func colorToAlpha(c color.Color) AlphaColor {
	return AlphaColorModel.Convert(c).(AlphaColor).Canon()
}

//Canon returns a new color with all values clamped to [0,1].
func (a AlphaColor) Canon() AlphaColor {
	return AlphaColor{clamp01(a.R), clamp01(a.G), clamp01(a.B), clamp01(a.A)}
}

func (a AlphaColor) c() (r, g, b, alpha C.double) {
	return C.double(a.R), C.double(a.G), C.double(a.B), C.double(a.A)
}

func (a AlphaColor) RGBA() (r, g, b, alpha uint32) {
	return ctoi(a.R), ctoi(a.G), ctoi(a.B), ctoi(a.A)
}

//Color returns a without the alpha channel.
func (a AlphaColor) Color() Color {
	return Color{a.R, a.G, a.B}
}

//These models can convert any color.Color to themselves.
//
//The conversion may be lossy.
var (
	ColorModel = color.ModelFunc(func(c color.Color) color.Color {
		if c, ok := c.(Color); ok {
			return c
		}
		if c, ok := c.(AlphaColor); ok {
			return Color{c.R, c.G, c.B}
		}
		r, g, b, _ := c.RGBA()
		return Color{cto01(r), cto01(g), cto01(b)}
	})
	AlphaColorModel = color.ModelFunc(func(c color.Color) color.Color {
		if c, ok := c.(AlphaColor); ok {
			return c
		}
		if c, ok := c.(Color); ok {
			return AlphaColor{c.R, c.G, c.B, 1}
		}
		r, g, b, a := c.RGBA()
		return AlphaColor{cto01(r), cto01(g), cto01(b), cto01(a)}
	})
)

//Basic colors.
var (
	Black       = Color{}
	White       = Color{1, 1, 1}
	Red         = Color{R: 1}
	Green       = Color{G: 1}
	Blue        = Color{B: 1}
	Transparent = AlphaColor{}
)
