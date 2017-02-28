// hellopdf project main.go
package main

import (
	"strings"

	"github.com/signintech/gopdf"
)

type FontStyle struct {
	Name  string
	Style string
	Size  int
}

type CellText struct {
	Left      float64
	Top       float64
	Rect      gopdf.Rect
	Option    gopdf.CellOption
	Text      string
	FontStyle FontStyle
}

type PdfTemplate struct {
	Titles  []CellText
	Cells   []CellText
	Details []CellText
}

func DrawCellText(pdf *gopdf.GoPdf, cellText CellText) (err error) {

	//draw border
	pdf.SetX(cellText.Left)
	pdf.SetY(cellText.Top)
	pdf.CellWithOption(&cellText.Rect, "",
		gopdf.CellOption{Border: cellText.Option.Border})

	//draw text
	err = pdf.SetFont(cellText.FontStyle.Name, cellText.FontStyle.Style,
		cellText.FontStyle.Size)
	if err != nil {
		return
	}
	charWitdh, err := pdf.MeasureTextWidth("中")
	if err != nil {
		return
	}
	bw := charWitdh / 3
	lineHeight := float64(cellText.FontStyle.Size)
	newLeft := cellText.Left + bw

	//charCount := int((cellText.Rect.W - bw*2) / charWitdh)
	lines := strings.Split(cellText.Text, "\n")
	cellwidth := cellText.Rect.W - 2*bw
	strs := []string{}
	for _, line := range lines {
		rs := []rune(line)
		rsleft := rs
		rsright := []rune("")
		strwidth, _ := pdf.MeasureTextWidth(string(rsleft))
		for strwidth > cellwidth {

			rsright = []rune(string(rsleft[len(rsleft)-1:]) + string(rsright))
			rsleft = rsleft[:len(rsleft)-1]
			strwidth, _ = pdf.MeasureTextWidth(string(rsleft))
			if !(strwidth > cellwidth) {
				strs = append(strs, string(rsleft))
				rsleft = rsright
				rsright = []rune("")
				strwidth, _ = pdf.MeasureTextWidth(string(rsleft))
			}
		}
		strs = append(strs, string(rsleft))
	}

	linecount := len(strs)
	newTop := cellText.Top + bw
	if cellText.Option.Align&gopdf.Middle == gopdf.Middle {
		newTop = newTop + (cellText.Rect.H-2*bw-float64(linecount)*lineHeight)/2
	} else if cellText.Option.Align&gopdf.Bottom == gopdf.Bottom {
		newTop = newTop + cellText.Rect.H - 2*bw - float64(linecount)*lineHeight
	}
	for i, str := range strs {
		lineTop := newTop + float64(i)*lineHeight
		newRect := gopdf.Rect{W: cellText.Rect.W - 2*bw, H: lineHeight}
		pdf.SetX(newLeft)
		pdf.SetY(lineTop)
		pdf.CellWithOption(&newRect, str,
			gopdf.CellOption{Align: cellText.Option.Align})
	}
	return
}
