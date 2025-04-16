package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

const (
	bannerHt = 95.0
	xIndent  = 40.0
	taxRate  = 0.09
)

type LineItem struct {
	UnitName       string
	PricePerUnit   int
	UnitsPurchased int
}

func main() {
	lineItems := []LineItem{
		{
			UnitName:       "2x6 Lumber - 8",
			PricePerUnit:   375,
			UnitsPurchased: 220,
		},
		{
			UnitName:       "Drywall Sheet",
			PricePerUnit:   822,
			UnitsPurchased: 50,
		},
		{
			UnitName:       "Paint",
			PricePerUnit:   1455,
			UnitsPurchased: 3,
		},
		{
			UnitName:       "Drill",
			PricePerUnit:   30000,
			UnitsPurchased: 2,
		},
		{
			UnitName:       "Hammer",
			PricePerUnit:   1000,
			UnitsPurchased: 5,
		},
		{
			UnitName:       "Ladder",
			PricePerUnit:   20000,
			UnitsPurchased: 4,
		},
	}

	subTotal := 0
	for _, li := range lineItems {
		subTotal += li.PricePerUnit * li.UnitsPurchased
	}
	tax := int(float64(subTotal) * taxRate)
	total := subTotal + tax

	pdf := gofpdf.New(gofpdf.OrientationPortrait, gofpdf.UnitPoint, gofpdf.PageSizeA4, "")
	w, h := pdf.GetPageSize()

	pdf.AddPage()

	//TOP BANNER
	pdf.SetFillColor(103, 60, 79)
	pdf.Polygon([]gofpdf.PointType{
		{0, 0},
		{w, 0},
		{w, bannerHt},
		{0, bannerHt * 0.9},
	}, "F")

	pdf.Polygon([]gofpdf.PointType{
		{0, h},
		{0, h - (bannerHt * 0.2)},
		{w, h - (bannerHt * 0.1)},
		{w, h},
	}, "F")

	//INVOICE
	pdf.SetFont("Arial", "B", 40)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt := pdf.GetFontSize()
	pdf.Text(xIndent, bannerHt-(bannerHt/2.0)+lineHt/3.0, "INVOICE")

	//BANNER - PHONE, EMAIL, DOMAIN
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-2.0*155.0, (bannerHt-(lineHt*1.5*3.0))/2)
	pdf.MultiCell(185.0, lineHt*1.5, "1234 456-789\nkacperbrozyna@hotmail.com\nhttps://github.com/Kacperbrozyna", gofpdf.BorderNone, gofpdf.AlignRight, false)

	//BANNER - ADDRESS
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(255, 255, 255)
	_, lineHt = pdf.GetFontSize()
	pdf.MoveTo(w-xIndent-125.0, (bannerHt-(lineHt*1.5*3.0))/2)
	pdf.MultiCell(125.0, lineHt*1.5, "123 Fake Street\nSome Town\nABC 123", gofpdf.BorderNone, gofpdf.AlignRight, false)

	//BILLED TO, INVOICE, DATE OF ISSUE
	_, sy := summaryBlock(pdf, xIndent, bannerHt+lineHt*2, "Billed To", "Client Name", "Client Address", "City", "Postal")
	summaryBlock(pdf, xIndent*2.0+lineHt*11.0, bannerHt+lineHt*2, "Invoice Number", "00000000123")
	summaryBlock(pdf, xIndent*2.0+lineHt*11.0, bannerHt+lineHt*6.25, "Date of Issue", "16/04/2025")

	//INVOICE TOTAL
	x, y := w-xIndent-124.0, bannerHt+lineHt*2.25
	pdf.MoveTo(x, y)
	pdf.SetFont("times", "", 14)
	_, lineHt = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.CellFormat(124.0, lineHt, "Invoice Total", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x, y = x+2.0, y+lineHt*1.5
	pdf.MoveTo(x, y)
	pdf.SetFont("times", "", 48)
	_, lineHt = pdf.GetFontSize()
	alpha := 58
	pdf.SetTextColor(72+alpha, 42+alpha, 55+alpha)
	pdf.CellFormat(124.0, lineHt, toUSDString(total), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x, y = x-2.0, y+lineHt*1.5

	if sy > y {
		y = sy
	}
	x, y = xIndent-20.0, y+30.0
	pdf.Rect(x, y, w-(xIndent*2.0)+40.0, 3.0, "F")

	//HEADER FOR DATA
	pdf.SetFont("times", "", 14)
	_, lineHt = pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	x, y = xIndent-2.0, y+lineHt
	pdf.MoveTo(x, y)
	pdf.CellFormat(w/2.65+1.5, lineHt, "Description", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100, lineHt, "Price Per Unit", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80, lineHt, "Quantity", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = w - xIndent - 2.0 - 109
	pdf.MoveTo(x, y)
	pdf.CellFormat(109, lineHt, "Amount", gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	y = y + lineHt*1.25

	//DATA
	for _, li := range lineItems {
		x, y = lineItem(pdf, x, y, li)
	}

	//FOOTER
	x, y = w/1.75, y+lineHt*2.25
	x, y = trailerLine(pdf, x, y, "Subtotal", subTotal)
	x, y = trailerLine(pdf, x, y, "Tax", tax)
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(x+20, y, x+220, y)
	y = y + lineHt*0.5
	x, y = trailerLine(pdf, x, y, "Total", total)

	err := pdf.OutputFileAndClose("invoice.pdf")
	if err != nil {
		panic(err)
	}
}

func toUSDString(cents int) string {
	centsStr := fmt.Sprintf("%d", cents%100)
	if len(centsStr) < 2 {
		centsStr = "0" + centsStr
	}

	return fmt.Sprintf("$%d.%s", cents/100, centsStr)
}

func trailerLine(pdf *gofpdf.Fpdf, x, y float64, label string, amount int) (float64, float64) {
	origX := x

	w, _ := pdf.GetPageSize()
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(180, 180, 180)
	pdf.MoveTo(x, y)
	pdf.CellFormat(80, lineHt, label, gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")
	x = w - xIndent - 2.0 - 109
	pdf.MoveTo(x, y)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(109, lineHt, toUSDString(amount), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignRight, false, 0, "")

	y = y + lineHt*1.5
	return origX, y
}

func lineItem(pdf *gofpdf.Fpdf, x, y float64, item LineItem) (float64, float64) {
	origX := x
	w, _ := pdf.GetPageSize()

	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	pdf.SetTextColor(50, 50, 50)
	pdf.MoveTo(x, y)
	x, y = xIndent-2.0, y+lineHt*0.5
	pdf.MoveTo(x, y)
	pdf.MultiCell(w/2.65+1.5, lineHt, item.UnitName, gofpdf.BorderNone, gofpdf.AlignLeft, false)
	tmp := pdf.SplitLines([]byte(item.UnitName), w/2.65+1.5)
	maxY := y + float64(len(tmp))*lineHt - lineHt
	x = x + w/2.65 + 1.5
	pdf.MoveTo(x, y)
	pdf.CellFormat(100, lineHt, toUSDString(item.PricePerUnit), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = x + 100.0
	pdf.MoveTo(x, y)
	pdf.CellFormat(80, lineHt, fmt.Sprintf("%d", item.UnitsPurchased), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")
	x = w - xIndent - 2 - 109
	pdf.MoveTo(x, y)
	pdf.CellFormat(109, lineHt, toUSDString(item.UnitsPurchased*item.PricePerUnit), gofpdf.BorderNone, gofpdf.LineBreakNone, gofpdf.AlignLeft, false, 0, "")

	if maxY > y {
		y = maxY
	}
	y = y + lineHt*1.75
	pdf.SetDrawColor(180, 180, 180)
	pdf.Line(xIndent-10, y, w-xIndent+10, y)

	return origX, y
}

func summaryBlock(pdf *gofpdf.Fpdf, x, y float64, title string, data ...string) (float64, float64) {
	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("times", "", 14)
	_, lineHt := pdf.GetFontSize()
	y = y + lineHt
	pdf.Text(x, y, title)
	pdf.SetTextColor(180, 180, 180)
	y = y + lineHt*.25
	for _, str := range data {
		y = y + lineHt*1.25
		pdf.Text(x, y, str)
	}

	return x, y
}
