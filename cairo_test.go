package cairo

func ExampleContext_Circle_drawEllipse() {
	//To achieve an elliptical arc, you can scale the current transformation
	//matrix by different amounts in the X and Y directions.
	//This draws an ellipse in the box given by Rectangle r.
	c.SaveRestore(func(c *Context) error {
		mid := Pt(r.Dx()/2, r.Dy()/2)
		c.Translate(r.Min.Add(mid))
		c.Scale(mid)
		c.Circle(UC)
		return nil
	})
}
