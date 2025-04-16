package main

import (
	"flag"
	"image/color"

	"github.com/jung-kurt/gofpdf"
)

type PDFOption func(*gofpdf.Fpdf)

type PDF struct {
	fpdf *gofpdf.Fpdf
	x, y float64
}

func FillColor(c color.RGBA) PDFOption {
	return func(f *gofpdf.Fpdf) {
		r, g, b := rgb(c)
		f.SetFillColor(r, g, b)
	}
}

func rgb(c color.RGBA) (int, int, int) {
	alpha := float64(c.A) / 255.0
	alphaWhite := int(255 * (1 - alpha))
	r := int(float64(c.R)*alpha) + alphaWhite
	g := int(float64(c.G)*alpha) + alphaWhite
	b := int(float64(c.B)*alpha) + alphaWhite

	return r, g, b
}

func (p *PDF) move(xDelta, yDelta float64) {
	p.x, p.y = p.x+xDelta, p.y+yDelta
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) moveAbs(x, y float64) {
	p.x, p.y = x, y
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) Text(text string) {
	p.fpdf.Text(p.x, p.y, text)
}

func (p *PDF) Polygon(pts []gofpdf.PointType, opts ...PDFOption) {
	for _, opt := range opts {
		opt(p.fpdf)
	}
	p.fpdf.Polygon(pts, "F")
}

func main() {
	name := flag.String("name", "Kacper Brozyna", "The name of the person who completed this course")
	flag.Parse()

	fpdf := gofpdf.New(gofpdf.OrientationLandscape, gofpdf.UnitPoint, gofpdf.PageSizeA4, "")
	w, h := fpdf.GetPageSize()
	fpdf.AddPage()
	pdf := PDF{
		fpdf: fpdf,
		x:    0,
		y:    0,
	}

	primary := color.RGBA{103, 60, 79, 255}
	secondary := color.RGBA{103, 60, 79, 220}

	//TOP AND BOTTOM BANNER
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{0, h / 9.0},
		{w - (w / 6), 0},
	}, FillColor(secondary))

	pdf.Polygon([]gofpdf.PointType{
		{w / 6.0, 0},
		{w, 0},
		{w, h / 9.0},
	}, FillColor(primary))

	pdf.Polygon([]gofpdf.PointType{
		{w, h},
		{w, h - h/8.0},
		{(w / 6), h},
	}, FillColor(secondary))

	pdf.Polygon([]gofpdf.PointType{
		{0, h},
		{0, h - h/8.0},
		{w - (w / 6), h},
	}, FillColor(primary))

	fpdf.SetFont("times", "B", 50)
	pdf.fpdf.SetTextColor(25, 25, 25)
	pdf.moveAbs(0, 100)
	_, lineHt := fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, "Certificate of Completion", gofpdf.AlignCenter)
	pdf.move(0, lineHt*2.0)

	fpdf.SetFont("Arial", "", 28)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, "This certificate is awarded to", gofpdf.AlignCenter)
	pdf.move(0, lineHt*2.0)

	fpdf.SetFont("times", "B", 42)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt, *name, gofpdf.AlignCenter)
	pdf.move(0, lineHt*1.75)

	fpdf.SetFont("Arial", "", 22)
	_, lineHt = fpdf.GetFontSize()
	fpdf.WriteAligned(0, lineHt*1.5, "For successfully completing all twenty programming exercices in the Gophercises programming course for budding Gophers (Go developers)", gofpdf.AlignCenter)

	pdf.move(0, lineHt*3.75)
	fpdf.ImageOptions("jump.png", w/2.0-50, pdf.y, 100.0, 0.0, false, gofpdf.ImageOptions{
		ReadDpi: true,
	}, -1, "")

	pdf.move(0, lineHt*3)
	fpdf.SetFillColor(100, 100, 100)
	fpdf.Rect(60, pdf.y, 240.0, 1.0, "F")
	fpdf.Rect(550, pdf.y, 240.0, 1.0, "F")

	fpdf.SetFont("arial", "", 12)
	pdf.move(0, lineHt/1.5)
	fpdf.SetTextColor(100, 100, 100)
	pdf.moveAbs(60+105, pdf.y)
	pdf.Text("Date")
	pdf.moveAbs(550+105, pdf.y)
	pdf.Text("Instructor")
	pdf.moveAbs(60.0+70, pdf.y-lineHt/1.25)
	fpdf.SetFont("times", "", 22)
	fpdf.SetTextColor(50, 50, 50)
	pdf.Text("16/04/2025")
	pdf.moveAbs(550.0+50, pdf.y)
	pdf.Text("Jonathan Calhoun")

	err := fpdf.OutputFileAndClose("certificate.pdf")
	if err != nil {
		panic(err)
	}
}
