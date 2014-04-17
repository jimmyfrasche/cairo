package ps

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/jimmyfrasche/cairo"
)

//When generating landscape PostScript output the surface should not be created
//with a width greater than the height. Instead create the surface
//with a height greater than the width and rotate the cairo drawing context.
//The "%PageOrientation" DSC comment is used by PostScript viewers to indicate
//the orientation of the page.
//
//The steps to create a landscape page are:
//Set the page size to a portrait size.
//Rotate user space 90 degrees counterclockwise and move the origin
//to the correct location.
//Insert the "%PageOrientation: Landscape" DSC comment.
func Example_landscape() {
	draw := func(c *cairo.Context, text string, width, height int) error {
		const Border = 50
		border := cairo.Pt(Border, Border)
		size := cairo.Pt(float64(width), float64(height))

		c.
			Rectangle(cairo.Rectangle{border, border.Sub(size)}).
			SetLineWidth(2).
			Stroke()

		c.
			SelectFont("Sans", 0, 0).
			SetFontSize(60).
			MoveTo(cairo.Pt(200, float64(height)/3)).
			ShowText(text)

		c.
			SetFontSize(18).
			MoveTo(cairo.Pt(120, float64(height)*2/3.)).
			ShowText(fmt.Sprintf("Width: %d points\t\tHeight: %d points"))

		return nil
	}

	//A4
	const (
		PageWidth  = 595
		PageHeight = 842
	)
	surface, err := New(os.Stdout, PageWidth, PageHeight, false, nil, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer surface.Close()

	c, err := cairo.New(surface)
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close()

	//Print portrait page
	surface.AddComment("PageOrientation", "Portrait")
	err = draw(c, "Portrait", PageWidth, PageHeight)
	if err != nil {
		log.Println(err)
		return
	}
	surface.ShowPage()

	//Print landscape page
	surface.AddComment("PageOrientation", "Landscape")

	//Move the origin to landscape origin and rotate counterclockwise
	c.
		Translate(cairo.Pt(0, PageHeight)).
		Rotate(-math.Pi / 2)

	err = draw(c, "Landscape", PageHeight, PageWidth)
	if err != nil {
		log.Println(err)
		return
	}
	surface.ShowPage()
}
