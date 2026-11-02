package tsc

import (
	"fmt"
	"math"
)

type Fit struct {
	Model      string    `json:"model"`
	Energies   []float64 `json:"energies_eV"`
	Amplitudes []float64 `json:"peak_amplitudes"`
	Offset     float64   `json:"offset"`
	Gamma      float64   `json:"transport_exponent"`
	RSS        float64   `json:"rss"`
	BIC        float64   `json:"bic"`
	Curve      []float64 `json:"-"`
}
type Solver struct {
	T      []float64
	S      float64
	Beta   float64
	Egrid  []float64
	Shapes [][]float64
	Dot    [][]float64
}

func NewSolver(t []float64, s, beta float64) (*Solver, error) {
	if len(t) < 25 || s <= 0 || beta <= 0 {
		return nil, fmt.Errorf("invalid solver")
	}
	o := &Solver{T: t, S: s, Beta: beta}
	for e := 0.215; e <= 0.3210001; e += 0.002 {
		v, err := Shape(t, e, s, beta)
		if err != nil {
			return nil, err
		}
		o.Egrid = append(o.Egrid, e)
		o.Shapes = append(o.Shapes, v)
	}
	n := len(o.Egrid)
	o.Dot = make([][]float64, n)
	for i := range o.Dot {
		o.Dot[i] = make([]float64, n)
		for j := 0; j <= i; j++ {
			v := dot(o.Shapes[i], o.Shapes[j])
			o.Dot[i][j] = v
			if j != i {
				o.Dot[j][i] = v
			}
		}
	}
	return o, nil
}
func dot(a, b []float64) float64 {
	v := 0.
	for i := range a {
		v += a[i] * b[i]
	}
	return v
}
func bic(rss float64, n, k int) float64 {
	if rss < 1e-25 {
		rss = 1e-25
	}
	return float64(n)*math.Log(rss/float64(n)) + float64(k)*math.Log(float64(n))
}
func fit1ByShape(y, shape []float64) (float64, float64) {
	v := dot(shape, y) / dot(shape, shape)
	if v < 0 {
		v = 0
	}
	return v, MeanSquare(y) - v*dot(shape, y)
}
func fit2ByShape(y, a, b []float64) (float64, float64, float64) {
	aa, bb, ab := dot(a, a), dot(b, b), dot(a, b)
	ya, yb := dot(y, a), dot(y, b)
	d := aa*bb - ab*ab
	if d < 1e-11 {
		x, r := fit1ByShape(y, a)
		return x, 0, r
	}
	x, z := (bb*ya-ab*yb)/d, (aa*yb-ab*ya)/d
	if x < 0 {
		x = 0
		z = math.Max(0, yb/bb)
	}
	if z < 0 {
		z = 0
		x = math.Max(0, ya/aa)
	}
	rss := MeanSquare(y) - 2*x*ya - 2*z*yb + x*x*aa + 2*x*z*ab + z*z*bb
	return x, z, math.Max(0, rss)
}
func makeCurve(a []float64, scale float64) (c []float64) {
	c = make([]float64, len(a))
	for i := range a {
		c[i] = scale * a[i]
	}
	return
}
func (o *Solver) FitOne(y []float64) (Fit, error) {
	if err := o.validY(y); err != nil {
		return Fit{}, err
	}
	best := Fit{RSS: math.Inf(1)}
	for j, e := range o.Egrid {
		amp, rss := fit1ByShape(y, o.Shapes[j])
		if rss < best.RSS {
			best = Fit{Model: "one", Energies: []float64{e}, Amplitudes: []float64{amp}, RSS: rss}
		}
	}
	best = o.refineOne(y, best)
	best.BIC = bic(best.RSS, len(y), 2)
	return best, nil
}
func (o *Solver) validY(y []float64) error {
	if len(y) != len(o.T) {
		return fmt.Errorf("data length mismatch")
	}
	for _, v := range y {
		if !good(v) {
			return fmt.Errorf("nonfinite observation")
		}
	}
	return nil
}
func (o *Solver) refineOne(y []float64, f Fit) Fit {
	e := f.Energies[0]
	for _, step := range []float64{0.0005, 0.0001, 0.00002} {
		for k := 0; k < 12; k++ {
			moved := false
			for _, trial := range []float64{e - step, e + step} {
				if trial < 0.205 || trial > 0.33 {
					continue
				}
				sh, _ := Shape(o.T, trial, o.S, o.Beta)
				a, r := fit1ByShape(y, sh)
				if r < f.RSS-1e-12 {
					e = trial
					f.Energies[0] = e
					f.Amplitudes[0] = a
					f.RSS = r
					moved = true
				}
			}
			if !moved {
				break
			}
		}
	}
	sh, _ := Shape(o.T, e, o.S, o.Beta)
	f.Curve = makeCurve(sh, f.Amplitudes[0])
	return f
}
func (o *Solver) FitTwo(y []float64) (Fit, error) {
	if err := o.validY(y); err != nil {
		return Fit{}, err
	}
	best := Fit{RSS: math.Inf(1)}
	n := len(o.Egrid)
	ys := make([]float64, n)
	for i := range ys {
		ys[i] = dot(y, o.Shapes[i])
	}
	yy := MeanSquare(y)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			aa, ab, bb := o.Dot[i][i], o.Dot[i][j], o.Dot[j][j]
			d := aa*bb - ab*ab
			if d <= 1e-10 {
				continue
			}
			a, b := (bb*ys[i]-ab*ys[j])/d, (aa*ys[j]-ab*ys[i])/d
			if a < 0 || b < 0 {
				continue
			}
			r := math.Max(0, yy-a*ys[i]-b*ys[j])
			if r < best.RSS {
				best = Fit{Model: "two", Energies: []float64{o.Egrid[i], o.Egrid[j]}, Amplitudes: []float64{a, b}, RSS: r}
			}
		}
	}
	if len(best.Energies) == 0 {
		return Fit{}, fmt.Errorf("two trap optimizer could not find positive components")
	}
	best = o.refineTwo(y, best)
	best.BIC = bic(best.RSS, len(y), 4)
	return best, nil
}
func (o *Solver) refineTwo(y []float64, f Fit) Fit {
	e1, e2 := f.Energies[0], f.Energies[1]
	for _, step := range []float64{0.0005, 0.0001, 0.00002} {
		for k := 0; k < 16; k++ {
			move := false
			for _, p := range [][2]float64{{e1 - step, e2}, {e1 + step, e2}, {e1, e2 - step}, {e1, e2 + step}} {
				a, b := p[0], p[1]
				if a < 0.205 || b > 0.33 || b-a < 0.0001 {
					continue
				}
				s1, _ := Shape(o.T, a, o.S, o.Beta)
				s2, _ := Shape(o.T, b, o.S, o.Beta)
				x, z, r := fit2ByShape(y, s1, s2)
				if r < f.RSS-1e-11 {
					f.RSS = r
					f.Amplitudes = []float64{x, z}
					e1, e2 = a, b
					move = true
				}
			}
			if !move {
				break
			}
		}
	}
	f.Energies = []float64{e1, e2}
	s1, _ := Shape(o.T, e1, o.S, o.Beta)
	s2, _ := Shape(o.T, e2, o.S, o.Beta)
	f.Curve = make([]float64, len(y))
	for i := range y {
		f.Curve[i] = f.Amplitudes[0]*s1[i] + f.Amplitudes[1]*s2[i]
	}
	return f
}

// FitOneOffset treats the constant dark-current contribution as a nuisance
// parameter, rather than inferring that every residual must be another trap.
func (o *Solver) FitOneOffset(y []float64) (Fit, error) {
	if err := o.validY(y); err != nil {
		return Fit{}, err
	}
	best := Fit{RSS: math.Inf(1)}
	for _, e := range o.Egrid {
		f := o.offsetAt(y, e)
		if f.RSS < best.RSS {
			best = f
		}
	}
	e := best.Energies[0]
	for _, step := range []float64{.0005, .0001, .00002} {
		for k := 0; k < 16; k++ {
			changed := false
			for _, ee := range []float64{e - step, e + step} {
				if ee < .205 || ee > .33 {
					continue
				}
				f := o.offsetAt(y, ee)
				if f.RSS < best.RSS-1e-12 {
					best = f
					e = ee
					changed = true
				}
			}
			if !changed {
				break
			}
		}
	}
	best.BIC = bic(best.RSS, len(y), 3)
	return best, nil
}
func (o *Solver) offsetAt(y []float64, e float64) Fit {
	sh, _ := Shape(o.T, e, o.S, o.Beta)
	n := float64(len(y))
	sumy, sumx, sumxx, sumxy := 0., 0., 0., 0.
	for i, v := range sh {
		sumy += y[i]
		sumx += v
		sumxx += v * v
		sumxy += v * y[i]
	}
	den := n*sumxx - sumx*sumx
	if den < 1e-20 {
		return Fit{RSS: math.Inf(1)}
	}
	a := (n*sumxy - sumx*sumy) / den
	if a < 0 {
		return Fit{RSS: math.Inf(1)}
	}
	b := (sumy - a*sumx) / n
	rss := 0.
	for i, v := range sh {
		d := y[i] - a*v - b
		rss += d * d
	}
	return Fit{Model: "one+offset", Energies: []float64{e}, Amplitudes: []float64{a}, Offset: b, RSS: rss}
}

// One-trap transport weighting is an explicitly hypothetical test, not a
// material-specific mobility law. Both E and gamma are fitted.
func (o *Solver) FitOneTransport(y []float64) (Fit, error) {
	if err := o.validY(y); err != nil {
		return Fit{}, err
	}
	best := Fit{RSS: math.Inf(1)}
	for gamma := -2.; gamma <= 1.0001; gamma += .25 {
		for _, e := range o.Egrid {
			f := o.transportAt(y, e, gamma)
			if f.RSS < best.RSS {
				best = f
			}
		}
	}
	e, g := best.Energies[0], best.Gamma
	for _, pair := range [][2]float64{{.0005, .1}, {.0001, .025}, {.00002, .005}} {
		for k := 0; k < 20; k++ {
			changed := false
			for _, v := range [][2]float64{{e - pair[0], g}, {e + pair[0], g}, {e, g - pair[1]}, {e, g + pair[1]}} {
				ee, gg := v[0], v[1]
				if ee < .205 || ee > .33 || gg < -3 || gg > 2 {
					continue
				}
				f := o.transportAt(y, ee, gg)
				if f.RSS < best.RSS-1e-12 {
					best = f
					e = ee
					g = gg
					changed = true
				}
			}
			if !changed {
				break
			}
		}
	}
	best.BIC = bic(best.RSS, len(y), 3)
	return best, nil
}
func (o *Solver) transportAt(y []float64, e, gamma float64) Fit {
	sh, _ := Shape(o.T, e, o.S, o.Beta)
	for i, T := range o.T {
		sh[i] *= math.Pow(T/230, gamma)
	}
	a, r := fit1ByShape(y, sh)
	return Fit{Model: "one+transport", Energies: []float64{e}, Amplitudes: []float64{a}, Gamma: gamma, RSS: r}
}
