package cairo

func ExamplePatch_coons() {
	coons := &Patch{}
	coons.MoveTo(ZP)
	coons.CurveTo(Pt(30, -30), Pt(60, 30), Pt(100, 0))
	coons.CurveTo(Pt(60, 30), Pt(130, 60), Pt(100, 100))
	coons.CurveTo(Pt(30, 70), Pt(-30, 30), Pt(0, 100))
	coons.SetCornerColors(Red, Green, Blue, Color{R: 1, G: 1})
}

func ExamplePatch_gouraudShadedTriangle() {
	gst := &Patch{}
	gst.MoveTo(Pt(100, 100))
	gst.LineTo(Pt(130, 130))
	gst.LineTo(Pt(130, 70))
	gst.SetCornerColors(Red, Green, Blue)
}
