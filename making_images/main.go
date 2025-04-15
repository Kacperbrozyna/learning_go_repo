package main

import (
	"os"

	svg "github.com/ajstarks/svgo"
)

func main() {

	f, err := os.OpenFile("demo.svg", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	canvas := svg.New(f)
	data := []struct {
		Month string
		Usage int
	}{
		{"Jan", 150},
		{"Feb", 10},
		{"Mar", 17},
		{"Apr", 250},
		{"May", 86},
		{"Jun", 90},
		{"Aug", 50},
		{"Sept", 67},
		{"Oct", 90},
		{"Nov", 110},
		{"Dec", 200},
	}

	width := len(data)*60 + 10
	height := 300
	max := 0
	threshhold := 120
	style := "fill:rgb(77,200,232)"

	for _, item := range data {
		if item.Usage > max {
			max = item.Usage
		}
	}

	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")
	for i, val := range data {
		percent := val.Usage * (height - 50) / max
		canvas.Rect(i*60+10, (height-50)-percent, 50, percent, style)
		canvas.Text(i*60+35, height-25, val.Month, "stroke:white; font-size:20pt;text-anchor:middle")
	}

	threshholdPercent := threshhold * (height - 50) / max
	canvas.Line(0, height-threshholdPercent, width, height-threshholdPercent, "stroke:red; opacity:0.8; stroke-width:2")
	canvas.Rect(0, 0, width, height-threshholdPercent, "fill:rgb(255,100,100); opacity:0.3")
	canvas.Line(0, height-50, width, height-50, "stroke: black; stroke-width:2")

	canvas.End()
}
