// paperdata derives manuscript macros solely from machine-readable experiment
// files so the written report cannot silently use old numerical outputs.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"trapresolution/internal/tsc"
)

func mustRead(path string, v any) {
	d, e := os.ReadFile(path)
	if e != nil {
		panic(e)
	}
	if e = json.Unmarshal(d, v); e != nil {
		panic(e)
	}
}
func main() {
	var d tsc.Summary
	var rows []tsc.Row
	var stress []tsc.StressRow
	mustRead("results/demo.json", &d)
	mustRead("results/resolution.json", &rows)
	mustRead("results/stress.json", &stress)
	var b strings.Builder
	f := func(k, v string) { fmt.Fprintf(&b, "\\newcommand{\\%s}{%s}\n", k, v) }
	f("FitOne", fmt.Sprintf("%.5f", d.First.Energies[0]))
	f("FitTwoA", fmt.Sprintf("%.5f", d.Second.Energies[0]))
	f("FitTwoB", fmt.Sprintf("%.5f", d.Second.Energies[1]))
	f("DeltaBIC", fmt.Sprintf("%.1f", d.DeltaBIC))
	f("PeakTemp", fmt.Sprintf("%.1f", d.PeakK))
	f("Prefactor", fmt.Sprintf("%.1f", d.EffectivePrefactor))
	f("PeakCount", fmt.Sprint(d.CombinedPeaks))
	for _, v := range rows {
		if v.SeparationMeV == 5 {
			f("FiveSelected", fmt.Sprint(v.TwoPreferred))
			f("FiveRecovered", fmt.Sprint(v.BothAccurate))
			f("FiveTrials", fmt.Sprint(v.Trials))
		}
	}
	for _, v := range stress {
		switch v.Scenario {
		case "matched":
			f("MatchedFalse", fmt.Sprint(v.NaiveFalseTwo))
		case "constant-background":
			f("BackgroundFalse", fmt.Sprint(v.NaiveFalseTwo))
			f("BackgroundAlternative", fmt.Sprint(v.AlternativeWins))
		case "temperature-weighting":
			f("TransportFalse", fmt.Sprint(v.NaiveFalseTwo))
			f("TransportAlternative", fmt.Sprint(v.AlternativeWins))
			f("StressTrials", fmt.Sprint(v.Trials))
		}
	}
	if e := os.MkdirAll("paper/data", 0755); e != nil {
		panic(e)
	}
	if e := os.WriteFile(filepath.Join("paper/data", "numbers.tex"), []byte(b.String()), 0644); e != nil {
		panic(e)
	}
}
