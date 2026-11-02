package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"trapresolution/internal/tsc"
)

func run() error {
	mode := "all"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	cfg := "configs/experiment.json"
	if len(os.Args) > 2 {
		cfg = os.Args[2]
	}
	d, err := os.ReadFile(cfg)
	if err != nil {
		return err
	}
	var p tsc.Protocol
	if err = json.Unmarshal(d, &p); err != nil {
		return err
	}
	n := 60
	if len(os.Args) > 3 {
		n, err = strconv.Atoi(os.Args[3])
		if err != nil {
			return err
		}
	}
	root := "results"
	if err = os.MkdirAll(root, 0755); err != nil {
		return err
	}
	if mode == "all" || mode == "demo" {
		s, pts, err := tsc.Demo(p)
		if err != nil {
			return err
		}
		if err = tsc.WriteJSON(filepath.Join(root, "demo.json"), s); err != nil {
			return err
		}
		xx := []float64{}
		a, b, c, d := []float64{}, []float64{}, []float64{}, []float64{}
		rows := [][]string{}
		for _, v := range pts {
			xx = append(xx, v.T)
			a = append(a, v.First)
			b = append(b, v.Second)
			c = append(c, v.Clean)
			d = append(d, v.FitTwo)
			rows = append(rows, []string{tsc.F(v.T), tsc.F(v.First), tsc.F(v.Second), tsc.F(v.Clean), tsc.F(v.Observed), tsc.F(v.FitOne), tsc.F(v.FitTwo)})
		}
		if err = tsc.WriteCSV(filepath.Join(root, "demo.csv"), []string{"temperature_K", "component1_relative", "component2_relative", "combined_clean", "combined_noisy", "one_fit", "two_fit"}, rows); err != nil {
			return err
		}
		if err = tsc.PlotSVG(filepath.Join(root, "demo.svg"), "Two components, one visible maximum", "Temperature (K)", "Relative intensity (arbitrary units)", xx, [][]float64{a, b, c, d}, []string{"Trap A", "Trap B", "Sum", "Two-trap fit"}, 1.16); err != nil {
			return err
		}
		fmt.Printf("demo: E=%v, fitted=%v, deltaBIC=%.2f, peaks=%d\n", s.TrueEnergies, s.Second.Energies, s.DeltaBIC, s.CombinedPeaks)
	}
	if mode == "all" || mode == "scan" {
		r, err := tsc.Scan(p, n)
		if err != nil {
			return err
		}
		if err = tsc.WriteJSON(filepath.Join(root, "resolution.json"), r); err != nil {
			return err
		}
		rows := [][]string{}
		x, model, rec := []float64{}, []float64{}, []float64{}
		for _, v := range r {
			rows = append(rows, []string{tsc.F(v.SeparationMeV), fmt.Sprint(v.Trials), fmt.Sprint(v.TwoPreferred), fmt.Sprint(v.BothAccurate), fmt.Sprint(v.OnePeak)})
			x = append(x, v.SeparationMeV)
			model = append(model, 100*float64(v.TwoPreferred)/float64(v.Trials))
			rec = append(rec, 100*float64(v.BothAccurate)/float64(v.Trials))
		}
		if err = tsc.WriteCSV(filepath.Join(root, "resolution.csv"), []string{"separation_meV", "trials", "two_strong_BIC", "both_energies_accurate", "single_visible_peak"}, rows); err != nil {
			return err
		}
		if err = tsc.PlotSVG(filepath.Join(root, "resolution.svg"), "Conditional resolution: fixed emission prefactor", "Trap separation (meV)", "Trials (%)", x, [][]float64{model, rec}, []string{"Model preferred", "Both energies accurate"}, 105); err != nil {
			return err
		}
		fmt.Println("resolution scan:", r)
	}
	if mode == "all" || mode == "stress" {
		r, err := tsc.Stress(p, n)
		if err != nil {
			return err
		}
		if err = tsc.WriteJSON(filepath.Join(root, "stress.json"), r); err != nil {
			return err
		}
		fmt.Println("model-mismatch control:", r)
	}
	// Produce a readable model-mismatch example, separate from aggregate trials.
	if mode == "all" || mode == "stress" {
		tt, err := p.Grid()
		if err != nil {
			return err
		}
		ss, err := tsc.CalibratedPrefactor(p.ReferenceE, p.ReferencePeak, p.Beta)
		if err != nil {
			return err
		}
		solver, err := tsc.NewSolver(tt, ss, p.Beta)
		if err != nil {
			return err
		}
		clean, observed, err := tsc.Synthetic(tt, []float64{0.26}, []float64{1}, ss, p.Beta, 0.01, 0, 0.015, rand.New(rand.NewSource(p.Seed+234)))
		if err != nil {
			return err
		}
		one, err := solver.FitOne(observed)
		if err != nil {
			return err
		}
		two, err := solver.FitTwo(observed)
		if err != nil {
			return err
		}
		offset, err := solver.FitOneOffset(observed)
		if err != nil {
			return err
		}
		sh, _ := tsc.Shape(tt, offset.Energies[0], ss, p.Beta)
		alt := make([]float64, len(tt))
		for j := range alt {
			alt[j] = offset.Amplitudes[0]*sh[j] + offset.Offset
		}
		errorOne, errorTwo, errorAlt := make([]float64, len(tt)), make([]float64, len(tt)), make([]float64, len(tt))
		for j := range tt {
			errorOne[j] = math.Abs(one.Curve[j] - clean[j])
			errorTwo[j] = math.Abs(two.Curve[j] - clean[j])
			errorAlt[j] = math.Abs(alt[j] - clean[j])
		}
		if err = tsc.PlotSVG(filepath.Join(root, "mismatch_example.svg"), "Model error when the one-trap physics is misspecified", "Temperature (K)", "Absolute error (normalized units)", tt, [][]float64{errorOne, errorTwo, errorAlt}, []string{"One-trap naive", "Two-trap naive", "One + background"}, 0); err != nil {
			return err
		}
	}
	if mode != "all" && mode != "demo" && mode != "scan" && mode != "stress" {
		return fmt.Errorf("usage: traplab [all|demo|scan|stress] [config] [trials]")
	}
	return nil
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
