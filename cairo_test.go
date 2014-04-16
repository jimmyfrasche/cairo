package cairo

import "log"

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
