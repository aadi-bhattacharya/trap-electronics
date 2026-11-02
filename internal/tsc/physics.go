// Package tsc studies idealized thermal-release curves. It does not model
// measured electrical current unless fast carrier relaxation and approximately
// temperature-independent collection, lifetime and mobility are assumed.
package tsc

import (
	"fmt"
	"math"
)

const KB = 8.617333262145e-5 // Boltzmann constant in eV/K

type Protocol struct {
	Start         float64 `json:"start_K"`
	End           float64 `json:"end_K"`
	Step          float64 `json:"step_K"`
	Beta          float64 `json:"heating_K_per_s"`
	ReferencePeak float64 `json:"reference_peak_K"`
	ReferenceE    float64 `json:"reference_energy_eV"`
	Seed          int64   `json:"seed"`
}

func (p Protocol) Grid() ([]float64, error) {
	if !(p.Start >= 1 && p.End > p.Start && p.Step >= 0.1 && p.Beta > 0) || !good(p.End) || !good(p.Step) || !good(p.Beta) {
		return nil, fmt.Errorf("invalid protocol")
	}
	n := int(math.Floor((p.End-p.Start)/p.Step+1e-8)) + 1
	if n < 20 || n > 5000 {
		return nil, fmt.Errorf("sample count outside 20..5000")
	}
	t := make([]float64, n)
	for i := range t {
		t[i] = p.Start + float64(i)*p.Step
	}
	return t, nil
}
func good(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }
func CalibratedPrefactor(e, tm, beta float64) (float64, error) {
	if !good(e) || !good(tm) || !good(beta) || e <= 0 || tm <= 0 || beta <= 0 {
		return 0, fmt.Errorf("invalid prefactor arguments")
	}
	s := beta * e / (KB * tm * tm) * math.Exp(e/(KB*tm))
	if !good(s) || s <= 0 {
		return 0, fmt.Errorf("nonfinite prefactor")
	}
	return s, nil
}
func Rate(e, s, t float64) float64 { return s * math.Exp(-e/(KB*t)) }
func Flux(temps []float64, e, s, beta float64) ([]float64, error) {
	if e <= 0 || s <= 0 || beta <= 0 || !good(e) || !good(s) || !good(beta) || len(temps) < 3 {
		return nil, fmt.Errorf("invalid flux input")
	}
	a := make([]float64, len(temps))
	H := 0.
	old := 0.
	for i, t := range temps {
		if !(t > 0) || !good(t) || (i > 0 && t <= temps[i-1]) {
			return nil, fmt.Errorf("invalid temperature order")
		}
		r := Rate(e, s, t)
		if i > 0 {
			H += 0.5 * (r + old) * (t - temps[i-1]) / beta
		}
		a[i] = r * math.Exp(-H)
		old = r
	}
	return a, nil
}

// Shape sets each ideal-trap component to maximum 1. The fitted coefficient
// therefore denotes a peak-signal scale, NOT a measured physical trap density.
func Shape(t []float64, e, s, beta float64) ([]float64, error) {
	a, err := Flux(t, e, s, beta)
	if err != nil {
		return nil, err
	}
	mx := 0.
	for _, v := range a {
		if v > mx {
			mx = v
		}
	}
	if mx <= 0 {
		return nil, fmt.Errorf("vanishing signal")
	}
	for i := range a {
		a[i] /= mx
	}
	return a, nil
}
func PeakTemp(t, a []float64) float64 {
	ix := 0
	for i := 1; i < len(a); i++ {
		if a[i] > a[ix] {
			ix = i
		}
	}
	return t[ix]
}
func PeakCount(a []float64) int {
	n := 0
	for i := 1; i < len(a)-1; i++ {
		if a[i] > a[i-1] && a[i] >= a[i+1] {
			n++
		}
	}
	return n
}
func MeanSquare(a []float64) float64 {
	z := 0.
	for _, v := range a {
		z += v * v
	}
	return z
}
func Max(a []float64) float64 {
	m := math.Inf(-1)
	for _, v := range a {
		if v > m {
			m = v
		}
	}
	return m
}
func Compose(t, energies, amplitudes []float64, s, beta float64) ([]float64, error) {
	if len(energies) == 0 || len(energies) != len(amplitudes) {
		return nil, fmt.Errorf("mismatched components")
	}
	a := make([]float64, len(t))
	for j, e := range energies {
		if amplitudes[j] < 0 || !good(amplitudes[j]) {
			return nil, fmt.Errorf("negative amplitude")
		}
		c, err := Shape(t, e, s, beta)
		if err != nil {
			return nil, err
		}
		for i := range a {
			a[i] += amplitudes[j] * c[i]
		}
	}
	return a, nil
}
