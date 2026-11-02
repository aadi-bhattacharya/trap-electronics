package tsc

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Row struct {
	SeparationMeV float64 `json:"separation_meV"`
	Trials        int     `json:"trials"`
	TwoPreferred  int     `json:"two_preferred"`
	BothAccurate  int     `json:"both_accurate"`
	OnePeak       int     `json:"one_peak"`
}
type StressRow struct {
	Scenario        string `json:"scenario"`
	Trials          int    `json:"trials"`
	NaiveFalseTwo   int    `json:"naive_false_two"`
	AlternativeWins int    `json:"alternative_wins"`
}
type Summary struct {
	Protocol           Protocol  `json:"protocol"`
	EffectivePrefactor float64   `json:"effective_prefactor_s_inv"`
	TrueEnergies       []float64 `json:"true_energies_eV"`
	CombinedPeaks      int       `json:"combined_peak_count"`
	PeakK              float64   `json:"peak_K"`
	First              Fit       `json:"one_fit"`
	Second             Fit       `json:"two_fit"`
	DeltaBIC           float64   `json:"delta_bic_one_minus_two"`
	Caveat             string    `json:"caveat"`
}
type Point struct{ T, First, Second, Clean, Observed, FitOne, FitTwo float64 }
type Results struct {
	Summary Summary
	Points  []Point
	Scan    []Row
	Stress  []StressRow
}

func Synthetic(t []float64, energies, amps []float64, s, beta, noise, gamma, baseline float64, rng *rand.Rand) ([]float64, []float64, error) {
	clean, err := Compose(t, energies, amps, s, beta)
	if err != nil {
		return nil, nil, err
	}
	mx := Max(clean)
	if mx <= 0 {
		return nil, nil, fmt.Errorf("empty signal")
	}
	y := make([]float64, len(t))
	for i, T := range t {
		clean[i] = clean[i]/mx*math.Pow(T/230, gamma) + baseline
		y[i] = clean[i] + noise*rng.NormFloat64()
	}
	return clean, y, nil
}
func Demo(p Protocol) (Summary, []Point, error) {
	t, err := p.Grid()
	if err != nil {
		return Summary{}, nil, err
	}
	s, err := CalibratedPrefactor(p.ReferenceE, p.ReferencePeak, p.Beta)
	if err != nil {
		return Summary{}, nil, err
	}
	o, err := NewSolver(t, s, p.Beta)
	if err != nil {
		return Summary{}, nil, err
	}
	rng := rand.New(rand.NewSource(p.Seed))
	clean, y, err := Synthetic(t, []float64{.26, .275}, []float64{1, 1}, s, p.Beta, .01, 0, 0, rng)
	if err != nil {
		return Summary{}, nil, err
	}
	f1, err := o.FitOne(y)
	if err != nil {
		return Summary{}, nil, err
	}
	f2, err := o.FitTwo(y)
	if err != nil {
		return Summary{}, nil, err
	}
	ca, _ := Shape(t, .26, s, p.Beta)
	cb, _ := Shape(t, .275, s, p.Beta)
	combinedUnscaledMax := 0.0
	for i := range ca {
		if ca[i]+cb[i] > combinedUnscaledMax {
			combinedUnscaledMax = ca[i] + cb[i]
		}
	}
	points := make([]Point, len(t))
	for i, T := range t {
		points[i] = Point{T, ca[i] / combinedUnscaledMax, cb[i] / combinedUnscaledMax, clean[i], y[i], f1.Curve[i], f2.Curve[i]}
	}
	return Summary{p, s, []float64{.26, .275}, PeakCount(clean), PeakTemp(t, clean), f1, f2, f1.BIC - f2.BIC, "Synthetic, fixed effective emission prefactor, independent Gaussian errors; components scaled to unit peak."}, points, nil
}
func Scan(p Protocol, trials int) ([]Row, error) {
	if trials < 1 || trials > 10000 {
		return nil, fmt.Errorf("invalid trials")
	}
	t, err := p.Grid()
	if err != nil {
		return nil, err
	}
	s, err := CalibratedPrefactor(p.ReferenceE, p.ReferencePeak, p.Beta)
	if err != nil {
		return nil, err
	}
	o, err := NewSolver(t, s, p.Beta)
	if err != nil {
		return nil, err
	}
	out := []Row{}
	for _, dm := range []float64{2, 5, 10, 15, 20} {
		rr := Row{SeparationMeV: dm, Trials: trials}
		for k := 0; k < trials; k++ {
			rng := rand.New(rand.NewSource(p.Seed + int64(10000*dm) + int64(k)))
			clean, y, err := Synthetic(t, []float64{.26, .26 + dm/1000}, []float64{1, 1}, s, p.Beta, .01, 0, 0, rng)
			if err != nil {
				return nil, err
			}
			if PeakCount(clean) == 1 {
				rr.OnePeak++
			}
			f1, err := o.FitOne(y)
			if err != nil {
				return nil, err
			}
			f2, err := o.FitTwo(y)
			if err != nil {
				return nil, err
			}
			if f1.BIC-f2.BIC > 10 {
				rr.TwoPreferred++
				if math.Abs(f2.Energies[0]-.26) <= .002 && math.Abs(f2.Energies[1]-(.26+dm/1000)) <= .002 {
					rr.BothAccurate++
				}
			}
		}
		out = append(out, rr)
	}
	return out, nil
}
func Stress(p Protocol, trials int) ([]StressRow, error) {
	if trials < 1 || trials > 10000 {
		return nil, fmt.Errorf("invalid trials")
	}
	t, err := p.Grid()
	if err != nil {
		return nil, err
	}
	s, err := CalibratedPrefactor(p.ReferenceE, p.ReferencePeak, p.Beta)
	if err != nil {
		return nil, err
	}
	o, err := NewSolver(t, s, p.Beta)
	if err != nil {
		return nil, err
	}
	out := []StressRow{}
	scenarios := []struct {
		name          string
		gamma, offset float64
		alt           string
	}{{"matched", 0, 0, "one"}, {"constant-background", 0, .015, "one+offset"}, {"temperature-weighting", -1.5, 0, "one+transport"}}
	for j, c := range scenarios {
		r := StressRow{Scenario: c.name, Trials: trials}
		for k := 0; k < trials; k++ {
			rng := rand.New(rand.NewSource(p.Seed + int64(j)*100000 + int64(k) + 900000))
			_, y, err := Synthetic(t, []float64{.26}, []float64{1}, s, p.Beta, .01, c.gamma, c.offset, rng)
			if err != nil {
				return nil, err
			}
			one, _ := o.FitOne(y)
			two, _ := o.FitTwo(y)
			if one.BIC-two.BIC > 10 {
				r.NaiveFalseTwo++
			}
			best := one.BIC
			if two.BIC < best {
				best = two.BIC
			}
			alt := one
			if c.alt == "one+offset" {
				alt, _ = o.FitOneOffset(y)
			}
			if c.alt == "one+transport" {
				alt, _ = o.FitOneTransport(y)
			}
			if c.alt != "one" && alt.BIC+10 < two.BIC {
				r.AlternativeWins++
			}
		}
		out = append(out, r)
	}
	return out, nil
}
func WriteJSON(path string, x any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	d, err := json.MarshalIndent(x, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(d, '\n'), 0644)
}
func WriteCSV(path string, head []string, rows [][]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	var sb strings.Builder
	sb.WriteString(strings.Join(head, ","))
	sb.WriteString("\n")
	for _, r := range rows {
		sb.WriteString(strings.Join(r, ","))
		sb.WriteString("\n")
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
func F(v float64) string { return strconv.FormatFloat(v, 'g', 12, 64) }
