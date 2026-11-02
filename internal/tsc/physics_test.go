package tsc

import (
	"math"
	"testing"
)

func prot() Protocol { return Protocol{100, 340, .5, .05, 230, .26, 20261102} }
func TestPeakIdentity(t *testing.T) {
	p := prot()
	s, e := CalibratedPrefactor(.26, 230, p.Beta)
	if e != nil {
		t.Fatal(e)
	}
	ts, _ := p.Grid()
	flux, e := Flux(ts, .26, s, p.Beta)
	if e != nil {
		t.Fatal(e)
	}
	if math.Abs(PeakTemp(ts, flux)-230) > .6 {
		t.Fatalf("peak=%v", PeakTemp(ts, flux))
	}
}
func TestRejectMalformed(t *testing.T) {
	p := prot()
	p.Step = 0
	if _, e := p.Grid(); e == nil {
		t.Fatal("bad grid accepted")
	}
	if _, e := Flux([]float64{100, 100, 101}, .26, 1, .05); e == nil {
		t.Fatal("unordered grid accepted")
	}
}
func TestConservation(t *testing.T) {
	p := prot()
	s, _ := CalibratedPrefactor(.26, 230, p.Beta)
	ts, _ := p.Grid()
	f, _ := Flux(ts, .26, s, p.Beta)
	area := 0.
	for i := 1; i < len(ts); i++ {
		area += (f[i] + f[i-1]) * .5 * (ts[i] - ts[i-1]) / p.Beta
	}
	if math.Abs(area-1) > .002 {
		t.Fatalf("normalized release area %g", area)
	}
}
func TestModelFit(t *testing.T) {
	p := prot()
	sum, _, err := Demo(p)
	if err != nil {
		t.Fatal(err)
	}
	if sum.CombinedPeaks != 1 {
		t.Fatalf("expected one maximum, got %d", sum.CombinedPeaks)
	}
	if math.Abs(sum.Second.Energies[0]-.26) > .003 || math.Abs(sum.Second.Energies[1]-.275) > .003 {
		t.Fatalf("recovery %v", sum.Second.Energies)
	}
	if sum.DeltaBIC <= 10 {
		t.Fatalf("model score %g", sum.DeltaBIC)
	}
}
func TestAlternativeModels(t *testing.T) {
	p := prot()
	ts, _ := p.Grid()
	s, _ := CalibratedPrefactor(.26, 230, p.Beta)
	sol, _ := NewSolver(ts, s, p.Beta)
	base, _ := Compose(ts, []float64{.26}, []float64{1}, s, p.Beta)
	maximum := Max(base)
	for i := range base {
		base[i] = base[i]/maximum + .01
	}
	fit, e := sol.FitOneOffset(base)
	if e != nil {
		t.Fatal(e)
	}
	if math.Abs(fit.Energies[0]-.26) > .004 {
		t.Fatalf("offset e=%v", fit.Energies)
	}
}
func MaxCopy(a []float64) float64 { return Max(a) }

func TestComponentPlotAdditivity(t *testing.T) {
	p := prot()
	_, points, err := Demo(p)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range points {
		if math.Abs(v.First+v.Second-v.Clean) > 1e-10 {
			t.Fatalf("nonadditive figure at %.2f K", v.T)
		}
	}
}
func TestIndependentGridRefinement(t *testing.T) {
	p := prot()
	s, _ := CalibratedPrefactor(.26, 230, p.Beta)
	coarse, _ := p.Grid()
	fc, _ := Flux(coarse, .26, s, p.Beta)
	p.Step = .125
	fine, _ := p.Grid()
	ff, _ := Flux(fine, .26, s, p.Beta)
	maxErr := 0.
	for i := range coarse {
		delta := math.Abs(fc[i] - ff[4*i])
		if delta > maxErr {
			maxErr = delta
		}
	}
	if maxErr/Max(ff) > 3e-4 {
		t.Fatalf("grid error %.8g", maxErr/Max(ff))
	}
}
