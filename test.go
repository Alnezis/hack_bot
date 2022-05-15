package main

import (
	"fmt"
	"math/rand"

	"github.com/fogleman/gg"
)

func main() {
	q := "AgACAgIAAxkBAAIBp2KAv5LPDr2AR1RSrsKtJ-nKfINoAALWuTEb7QUJSEWuYgg7fHgHAQADAgADeQADJAQ"

	if q == "AgACAgIAAxkBAAIBp2KAv5LPDr2AR1RSrsKtJ-nKfINoAALWuTEb7QUJSEWuYgg7fHgHAQADAgADeQADJAQ" {
		fmt.Println(true)
	}
	const W = 1024
	const H = 1024
	dc := gg.NewContext(W, H)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	//for i := 0; i < 10; i++ {
	x1 := rand.Float64() * 200
	y1 := rand.Float64() * 200
	x2 := rand.Float64() * 400
	y2 := rand.Float64() * 400
	r := 3.0
	g := 3.0
	b := 3.0
	a := rand.Float64()*0.5 + 0.5
	w := 6.0
	dc.SetRGBA(r, g, b, a)
	dc.SetLineWidth(w)
	dc.DrawLine(x1, y1, x2, y2)
	dc.Stroke()
	///	}
	dc.SavePNG("out.png")
}
