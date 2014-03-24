package cairo

func ExampleMesh_withTwoPatches() {
	red, green, blue := Color{R: 1}, Color{G: 1}, Color{B: 1}

	m := NewMesh()

	//Add a Coons patch
	err := m.BeginPatch().
		MoveTo(ZP).
		CurveTo(Pt(30, -30), Pt(60, 30), Pt(100, 0)).
		CurveTo(Pt(60, 30), Pt(130, 60), Pt(100, 100)).
		CurveTo(Pt(60, 70), Pt(30, 130), Pt(0, 100)).
		CurveTo(Pt(30, 70), Pt(-30, 30), ZP).
		SetCornerColor(0, red).
		SetCornerColor(1, green).
		SetCornerColor(2, blue).
		SetCornerColor(3, Color{R: 1, G: 1}).
		EndPatch()

	if err != nil {
		panic(err) //this is not how to handle errors, outside of examples
	}

	//Add a Gouraud-shaded triangle
	err = m.BeginPatch().
		MoveTo(Pt(100, 100)).
		LineTo(Pt(130, 130)).
		LineTo(Pt(130, 70)).
		SetCornerColor(0, red).
		SetCornerColor(1, green).
		SetCornerColor(2, blue).
		EndPatch()

	if err != nil {
		panic(err)
	}
}
