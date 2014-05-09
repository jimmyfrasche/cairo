package cairo

import (
	"image/png"
	"log"
	"os"
)

//A simple hello world that writes a blue "Hello, world" to hello.png.
//
//Adapted from http://cairographics.org/FAQ/#minimal_C_program .
func Example_helloWorld() {
	surface, err := NewImageSurface(FormatARGB32, 240, 80)
	if err != nil {
		log.Fatalln(err)
	}
	defer surface.Close()

	cr, err := New(surface)
	if err != nil {
		log.Fatalln(err)
	}
	defer cr.Close()

	cr.
		SelectFont("serif", 0, WeightBold).
		SetFontSize(32).
		SetSourceColor(Blue).
		MoveTo(Pt(10, 50)).
		ShowText("Hello, world")

	img, err := surface.ToImage()
	if err != nil {
		log.Fatalln(err)
	}

	f, err := os.Create("hello.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Fatalln(err)
	}
}

func ExampleContext_Circle_drawEllipse() {
	is, err := NewImageSurface(FormatARGB32, 500, 500)
	if err != nil {
		log.Fatalln(err)
	}
	c, err := New(is)
	if err != nil {
		log.Fatalln(err)
	}
	//To achieve an elliptical arc, you can scale the current transformation
	//matrix by different amounts in the X and Y directions.
	//This draws an ellipse in the box given by Rectangle r.
	r := Rect(0, 0, 50, 60)
	c.SaveRestore(func(c *Context) error {
		mid := r.Size().Div(2)
		c.
			Translate(r.Min.Add(mid)).
			Scale(mid).
			Circle(UC)
		return nil
	})
}
