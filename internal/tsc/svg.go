package tsc

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// Minimal vector plots intentionally avoid an external plotting dependency.
func PlotSVG(path, title, xlabel, ylabel string, x []float64, ys [][]float64, labels []string, ymax float64) error {
	if len(x) < 2 || len(ys) == 0 || len(ys) != len(labels) {
		return fmt.Errorf("invalid plot")
	}
	for _, y := range ys {
		if len(y) != len(x) {
			return fmt.Errorf("mismatched plot vectors")
		}
	}
	if ymax <= 0 {
		for _, y := range ys {
			if Max(y) > ymax {
				ymax = Max(y)
			}
		}
		ymax *= 1.1
	}
	if ymax <= 0 {
		ymax = 1
	}
	xlo, xhi := x[0], x[len(x)-1]
	const W, H = 820., 470.
	left, right, top, bottom := 82., 25., 59., 66.
	pw, ph := W-left-right, H-top-bottom
	colors := []string{"#345a8a", "#c47b42", "#2e927e", "#888d9c"}
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="820" height="470" viewBox="0 0 820 470"><rect x="0" y="0" width="820" height="470" fill="white"/><style>text{font-family:Arial,Helvetica,sans-serif;fill:#263245}</style><text x="82" y="30" font-size="18" font-weight="bold">%s</text>`, title))
	for j := 0; j <= 4; j++ {
		xx := left + pw*float64(j)/4
		yy := top + ph*float64(j)/4
		b.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#e7e9ef" stroke-width="1"/>`, left, yy, W-right, yy))
		b.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-size="12" text-anchor="end">%.2g</text>`, left-8, yy+4, ymax*(1-float64(j)/4)))
		b.WriteString(fmt.Sprintf(`<text x="%.2f" y="%.2f" font-size="12" text-anchor="middle">%.3g</text>`, xx, H-bottom+21, xlo+(xhi-xlo)*float64(j)/4))
	}
	b.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#283648" stroke-width="1.4"/><line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="#283648" stroke-width="1.4"/>`, left, top, left, H-bottom, left, H-bottom, W-right, H-bottom))
	for k, y := range ys {
		co := colors[k%len(colors)]
		b.WriteString(fmt.Sprintf(`<polyline fill="none" stroke="%s" stroke-width="2.4" points="`, co))
		stride := int(math.Max(1, float64(len(x))/1000))
		for i := 0; i < len(x); i += stride {
			xp := left + (x[i]-xlo)/(xhi-xlo)*pw
			yp := top + ph*(1-y[i]/ymax)
			b.WriteString(fmt.Sprintf("%.2f,%.2f ", xp, yp))
		}
		b.WriteString(`"/>`)
		b.WriteString(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="3"/><text x="%.2f" y="%.2f" font-size="12">%s</text>`, left+float64(k)*176, 44.0, left+float64(k)*176+19, 44.0, co, left+float64(k)*176+24, 48.0, labels[k]))
	}
	b.WriteString(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="13" text-anchor="middle">%s</text><text transform="translate(18,%.1f) rotate(-90)" font-size="13" text-anchor="middle">%s</text></svg>`, left+pw/2, H-13, xlabel, top+ph/2, ylabel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(b.String()), 0644)
}
