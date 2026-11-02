# Independent scientific / implementation audit

**Purpose:** record exactly what was implemented and what remains hypothetical. This document is not part of your authored course prose.

## 1. Source integrity and scope

- Primary β-In₂S₃ experimental anchor: Pai et al. (2005) user's PDF. It reports a 0.05 K/s ramp, 100 K start, approximately 2.6 eV optical gap for its specific films, and a proposed 0.26 eV Cl-related trap peaking near 230 K.
- General TSC kinetics: Dussel and Bube (1967), user-provided paper. Their treatment of retrapping and changing carrier lifetime warns against universal peak-width interpretation.
- Crystallography: Pistor et al. (2016) space group I4₁/amd at 309 K, a=7.6231 Å, c=32.358 Å. Their crystallographic ordering does not chemically identify TSC levels.
- Band structure: Zhao et al. (2012) indirect *calculated* fundamental gap, not a proof that Pai's optical measurement must exhibit the same numerical threshold.
- Contemporary motivation: Kaur (2026), sulfur-vacancy engineering in β-In₂S₃ photoelectrodes, not direct corroboration of a synthetic 0.275 eV level.

## 2. Physical forward model

The basic program integrates e(T)=s exp(-E/kBT), dn_t/dt=-e n_t, and the normalized release flux e(T)n_t(T). It neither solves for n_free(T) nor predicts a current in amperes. It neglects retrapping and assumes fast relaxation and temperature-independent conversion between release flux and measured current. It has a **common fixed effective prefactor** and separate nonnegative *peak-signal amplitudes*. The latter are not trap densities or capture cross sections.

The prefactor is inverted from the 0.26 eV and 230 K pair under the ideal peak condition; it is not independent experimental knowledge. Parameter recovery conditional on that prefactor is therefore optimistic.

## 3. Numerical solver and test gates

- Numerical integration: composite trapezoid rule on a monotonic T grid. Output `results/demo.csv` is 481 rows from 100 through 340 K at 0.5 K steps.
- `TestPeakIdentity`: peak lies within 0.6 K of 230 K at the calibrated E,s,β.
- `TestConservation`: ∫flux dt ≈ 1 for a completely released initially occupied model population.
- `TestIndependentGridRefinement`: compares a 0.5 K integration with a finer 0.125 K grid at shared points.
- `TestComponentPlotAdditivity`: independently normalized visual display was an earlier bug; the corrected plot uses the SAME normalization across components and total.
- `TestModelFit`: recovers the planted 0.260/0.275 eV levels to ±0.003 eV in one noisy dataset.
- `TestAlternativeModels`: checks the add-offset fit on synthetic observations.
- `go vet` and `go test ./...` form part of `make all` and GitHub Actions.

## 4. Inverse model details and fragility

One-trap has E and A (k=2); two-trap has ordered E1<E2 and A1,A2 (k=4); offsets add one parameter; transport exponent adds one. A coarse bounded 2-meV energy grid supplies the initial guess, then local coordinate refinement searches smaller steps. The fitter uses analytic nonnegative amplitude recovery. **This is not a guarantee of globally optimal nonlinear parameters.** Error independent in T and constant in normalized units is hypothetical. BIC uses estimated Gaussian variance from RSS; comparisons become non-regular when a second amplitude shrinks to zero or trap energies coincide.

### Data-generating controls

- Matched: one true trap fitted to one/two first-order candidates.
- Offset mismatch: one true trap with an additive constant +0.015 normalized units; candidate family incorrectly assumes no dark background.
- Temperature-factor mismatch: one true trap weighted by (T/230)^(-1.5); this is illustrative, not a material-specific mobility result.
- Expanded alternatives: one+offset and one+free transport exponent. Their existence shows how incomplete physical candidate sets can create false second-trap inference.

**Statistical interpretation:** 200 trials per condition give descriptive finite-sample proportions at fixed seed 20261102. There are no confidence-calibrated, universally material-dependent energy-resolution thresholds. The Monte Carlo is synthetic, not experimental measurement.

## 5. Audit of plotting and machine-readable numbers

Figure `results/demo.svg` must satisfy y_1+y_2 = y_total at every point before random noise is added. File `results/demo.csv` contains both the true curves and fits so the visualization is independently checkable. `cmd/paperdata` reads only JSON results and writes LaTeX macros in `paper/data/numbers.tex`, preventing transcribed numerical figures from becoming stale after reruns.

The revised IEEE paper contains 14 `AUTHOR WRITES` prompts (A0--A13), each with specific physics and derivation objectives; the companion `AUTHOR_WRITING_GUIDE.md` supplies deeper research guidance. The remaining prose is limited to sourced experimental parameters and reproducible numerical methods/results. Instructors should be able to distinguish the student's authored physical interpretation from automated numerical work.

## 6. Remaining risks / explicitly unimplemented

1. No exact experimental TSC digitization or fit to Pai's measured curves.
2. No coupling of mobile carrier density to traps, recombination or deep reservoir occupations.
3. No contact resistance, leakage, self-heating, or detailed sampling model.
4. No measured temperature dependence of mobility, capture cross sections or emission prefactors.
5. No profile likelihood, bootstrap confidence intervals, independent validation datasets or global optimization guarantee.
6. No chemical defect identification from the fitted current curves.
7. No DFT band calculation. The project studies kinetic resolution only.
8. No machine-derived claims about actual device performance improvements.

**Priority future experiment:** allow prefactors to vary and study profile-likelihood ridges in E–s space before interpreting any fitted meV separation; then compare with repeated heating-rate spectra and background-subtracted measurements.
