# Author's writing guide

**Working title:** *When Two Traps Look Like One: A Computational Investigation of Defect-State Resolution*

This file is a **thinking and drafting aid**, not text to copy into your report. The companion `paper/root.tex` preserves the supplied `ieeeconf` conference layout, with **AUTHOR WRITES A0–A13** at every place that needs your own prose. The equations, source-checked facts, and Go-generated results stay in the paper because they are useful fixed points around which to explain the physics.

## Use this guide efficiently

1. Read the named source page or section and make brief handwritten notes. For equations, reproduce every step on paper, including units and the assumptions that permit the step.
2. Close the source and write your explanation, in your own words, next to the corresponding `AUTHOR WRITES` prompt in `paper/root.tex`.
3. After each section, check the **must answer** questions below. Cite factual material; do not cite your own derivation as though it was copied from a source.
4. Remove a prompt only after you have written your replacement. The `\guidancetrue` switch shows prompts. Switching it to `\guidancefalse` hides them but does not magically provide the missing paragraphs.
5. The current compiled PDF is a compact **working scaffold**, not a target page count. The final report should be roughly five pages if that is what your course requires; write concisely and check the page count after replacing prompts.

## Mandatory course coverage

The final prose needs a connected causal story: **crystal structure and bonding → electronic bands and local defect states → carrier population and motion → TSC observable → ambiguity of defect extraction → practical ways forward**. This keeps the computation anchored in *semiconductor physics*, which is the actual course requirement.

| Course requirement | Where your explanation belongs | Non-negotiable detail |
|---|---|---|
| Crystal structure | A2 | Phase, symmetry, cation coordination, structural vacancy ordering |
| Energy bands | A3–A4 | Gap meaning, direct/indirect, band curvature, effective mass, Fermi level |
| Charge carriers and transport | A4–A5 | Electrons/holes, drift, diffusion, photoexcitation, trapping |
| One current research challenge | A1, A10–A11 | Reliability of extracting trap states from thermal-current data |
| Future prospects | A12 | Improvements in spectroscopy or material optimization, without hype |

## A0. Abstract: make every sentence do a job

**Write this last; 90–110 words.** Your first sentence should identify the physics setting, not make a generic claim about modern technology. The second should precisely state the *one peak/two levels* inference question. The third should name the assumptions (first-order, synthetic Go modeling, fixed effective prefactor). The fourth should report a numerical observation and acknowledge model-mismatch risks. **Must answer:** Could a reader tell, without opening the code, that no experimental β-In₂S₃ spectrum was fitted?

## A1. Introduction: establish the stakes

**Writing route (approximately 100–130 words):** (a) begin with a perfect periodic crystal and allowed electronic bands; (b) introduce a local imperfection that can trap a carrier, disturbing photoconductivity and recombination; (c) explain why TSC can reveal thermally activated release; (d) close with the inverse question. The distinction you must make is between **the number of visible maxima**, **the number of components selected by a particular statistical model**, and **the number/identity of real microscopic defects**.

**Read:** Pai et al., PDF p. 1 introduction; Dussel and Bube, PDF p. 1 introduction. Don't use the word “resolution” without saying whether you mean visual, statistical, or parameter accuracy.

## A2. Crystal structure and bonding

**Questions to resolve before writing:**

- What makes β-In₂S₃ a *defect spinel*? Which sites are occupied by indium, which are systematically unoccupied, and what is the role of sulfur coordination?
- Why is a tetragonal phase not equivalent to a cubic spinel? Identify the cited space group **I4₁/amd** and explain the general meaning of symmetry reduction and vacancy ordering.
- What holds the material together? Discuss the ionic/covalent contribution to In–S bonding cautiously; avoid depicting all bonds as purely one type.
- What does **not** follow? Regularly vacant sites required by an ideal crystallographic model do not individually demonstrate 0.26 eV electrical traps.

**Read:** Pistor et al. (2016), phase-structure discussion and crystallographic figures/tables; Pai et al. PDF p. 1 for their material framing. A hand-drawn, captioned unit-cell schematic is optional; do not fabricate a real-space structure from your imagination.

**Your paragraph structure:** (1) lattice identity, (2) cation/sulfur coordination, (3) vacancy ordering and symmetry, (4) what this can or cannot tell us about electronic defects.

## A3. Band structure and effective mass

On scratch paper, derive the *directional* band effective-mass relation from the semiclassical equations:

\[
 v_g=\frac{1}{\hbar}\frac{\partial E}{\partial k},\qquad
 \hbar\frac{dk}{dt}=F,\qquad
 \frac{dv_g}{dt}=\frac{1}{\hbar^2}\frac{\partial^2E}{\partial k^2}F.
\]

Match the final expression to Newton's law to read off \(m^*\). In an anisotropic crystal, the physically complete quantity is generally a mass **tensor**, not one global scalar.

Sketch two E–k diagrams yourself: one with the valence-band maximum and conduction-band minimum at the same wavevector, one with them at different wavevectors. Connect the difference to momentum conservation and optical transitions. Zhao et al. (2012) calculate an **indirect fundamental gap**; the Pai et al. films were reported to have an **optical absorption gap of about 2.6 eV**. These values come from different methods and specimens. Do not invent a consistent single gap or claim that measured optical absorption determines the indirect fundamental gap without qualification.

**Read:** Zhao et al. (2012) abstract/analysis; Pai et al. PDF p. 2 sample properties. **Must answer:** What physical information does curvature give that gap size alone cannot?

## A4. Fermi level and carrier statistics

Start with the probability that an available electron state at energy \(E\) is occupied. Explain the significance of the Fermi–Dirac function's half-occupation at \(E_F\). Then write the carrier density as an integral of density of states times occupancy; show what nondegenerate approximation gives the familiar exponentials.

\[
 n\simeq N_Ce^{-(E_C-E_F)/(k_BT)},\qquad
 p\simeq N_Ve^{-(E_F-E_V)/(k_BT)}.
\]

Explain when \(np=n_i^2\) holds: thermal equilibrium and nondegenerate statistics. Under illumination, explain why electron and hole populations need not be described by one common equilibrium Fermi level. **Do not** choose a numeric donor density or Fermi energy for Pai's particular thin films; these were not established by that paper's TSC measurements.

**Must answer:** If the material is n-type, what can you say qualitatively about electron versus hole concentration? What can you *not* calculate without doping and density-of-states information?

## A5. Drift, diffusion, and measured current

Derive \(\sigma=q(n\mu_n+p\mu_p)\) by adding the electrical current from the two carrier species. Electron charge is negative, but the conventional electron current direction under field is positive in the standard \(J_n\) formula because electron velocity opposes the electric field. State the two current-density relations and make sure the diffusion signs are different:

\[
 \mathbf{J}_n=qn\mu_n\mathbf{E}+qD_n\nabla n,
 \qquad
 \mathbf{J}_p=qp\mu_p\mathbf{E}-qD_p\nabla p.
\]

Under nondegenerate equilibrium conditions, recall \(D/\mu=k_BT/q\). Do not assume this applies unchanged in a degenerate system.

**Device explanation:** Take the actual geometry: a biased thin film illuminated while cold, followed by heating in darkness. Trace each stage as a carrier population rather than a vague “electrons flow” story. State why emission \(e(T)n_t(T)\) is a *release rate per volume* in the model, whereas electrical current depends additionally on mobile carrier density, mobility, electric field, geometry and recombination lifetime.

## A6. Reconstruct the Pai experiment from the source

**Primary source:** Pai et al. (2005), PDF pp. 2–4 (journal pp. 742–744). Use their values, not approximate values copied from secondary accounts.

| Experimental detail | Source-supported value |
|---|---|
| Films | Spray-pyrolysed β-In₂S₃, varying spray-solution In:S ratio |
| Initial cooling | 100 K |
| Illumination | White light, 10 mW/cm², 60 seconds |
| Heating | 3 K/min (0.05 K/s) |
| Bias | 30 V |
| Electrode separation | 0.5 cm |
| TSC measurement range | 100–500 K |
| Film thickness | Approximately 1 µm |
| Conductivity | n-type, from hot-probe measurement |

**Writing route:** Explain each physical step and expected observable rather than simply listing the protocol. From Fig. 1, notice that different preparation ratios yield different peak amplitudes and shoulders. From Table II, cite the four activation energies with the authors' proposed origin. An XPS observation supporting elemental presence is *not* proof that a TSC peak belongs to that element's exact electronic center.

## A7. Derive first-order TSC from the rate law

On paper, start at

\[
 e(T)=s\exp[-E/(k_BT)],\quad \frac{dn_t}{dt}=-e(T)n_t,\quad \frac{dT}{dt}=\beta.
\]

Change independent variable from time to temperature using the chain rule, integrate from \(T_0\) to \(T\), and obtain the occupancy formula already in the TeX scaffold. Then take the logarithmic derivative of the release flux \(F(T)=e(T)n_t(T)\):

\[
 \frac{d\ln F}{dT}=\frac{E}{k_BT^2}-\frac{e(T)}{\beta}.
\]

Set that expression to zero at the maximum. Verify both sides of \(e(T_m)=\beta E/(k_BT_m^2)\) have units s⁻¹. This is a **first-order, constant-prefactor peak condition**, not a universal result for arbitrary trapping kinetics.

**Must answer:** Why does emission initially rise, and what causes the response to disappear at high temperature even though thermal excitation gets easier?

## A8. Why the release flux is not automatically the current

Read Dussel and Bube (1967), especially PDF pp. 1–2, 8–9, and 13–15. Their more general analysis contains recombination and retrapping, so the mobile-carrier lifetime can change during heating. For the simple sum of release rates to stand in for a current signal, you additionally need a suitable fast-relaxation/constant-mobility or collection approximation. If you describe a more realistic model, sketch coupled equations for free electron density and occupied traps, but do not claim this project implemented a full carrier-transport PDE.

**Must answer:** Could one trap broaden or shift its TSC peak when retrapping becomes important? If so, why might the inverse fit mistakenly add another component?

## A9. Computational model and model-selection mathematics

The Go implementation is in `internal/tsc/` and is invoked by `make all`. Read `physics.go`, `fit.go` and `experiments.go`; you should be able to explain the role of each without memorizing the code.

**Baseline assumptions and values:**

- \(T=100\ldots 340\) K in 0.5 K steps; \(\beta=0.05\) K/s, matching Pai's heating rate but not their full measurement range.
- Electron trap 1: 0.260 eV, motivated by one of Pai's reported thermal activation energies; electron trap 2: **hypothetical** 0.275 eV.
- \(s\approx1419.9\,\mathrm{s^{-1}}\) is calibrated from a 230 K maximum using the simple first-order formula. It is **not** a measured microscopic attempt frequency.
- Each component is normalized to its own peak before amplitude fitting; the amplitude is not a directly measured trap density.
- Independent Gaussian noise with standard deviation 0.01 relative to a normalized peak; deterministic seed 20261102.

The two candidate models contain one or two normalized spectral shapes. The fitting minimizes residual squares. Compute \(\mathrm{BIC}=N\ln(\mathrm{RSS}/N)+k\ln N\) and why the added parameters cost complexity. Model preference here means \(\Delta\mathrm{BIC}=\mathrm{BIC}_{1}-\mathrm{BIC}_{2}>10\). Accurate recovery here means both fitted energies are within 0.002 eV of the input values. These are **defined simulation criteria**, not laboratory calibration standards.

**Must answer:** Why is identification easier when the emission prefactor is known? How could the results change if \(s\), the offset, or \(\mu(T)\) were also unknown?

## A10. Interpret the simulated single peak

From `results/demo.json`, input energies are 0.260 and 0.275 eV, but the combined idealized spectrum has **one** maximum, at 235.0 K on the sampling grid. This supports the limited claim that *two generating components need not produce two visible maxima*.

From `results/resolution.json`, at 5 meV separation, 138/200 runs preferred two components but 67/200 met both model-preference and energy-accuracy criteria. Give both fractions, not just the nicer one. These experiments use a fixed common prefactor; their rates are not universal resolution thresholds. Make clear that a single clean curve can still contain information in its **shape**, not just the number of maxima.

**Suggested explanation structure:** Describe the plot; explain why components overlap; report the fitting comparison; state the assumptions needed for interpretation.

## A11. Interpret model-mismatch stress tests

All three tests in `results/stress.json` generate **one** trap. They are not claims about how many defects exist in β-In₂S₃.

| Generating condition | Restricted fitter incorrectly prefers 2 | Alternative one-trap family preferred over 2 |
|---|---:|---:|
| Matched first-order model | 0/200 | not applicable |
| Constant offset 0.015 | 200/200 | 200/200 |
| Temperature weighting (T/230 K)⁻³ᐟ² | 200/200 | 171/200 |

The exponent is an **illustrative temperature weighting**, not a verified β-In₂S₃ mobility law. An extra component can absorb neglected physics, so fitting two components is not proof of two actual defects. Explain that the alternative families fit a constant offset or a temperature-dependent multiplier; this addresses whether the second trap is truly necessary.

**Must answer:** What new experiment would reduce this ambiguity? Choose one, for example a dark-current measurement at the same heating rate, independent mobility-versus-temperature data, several heating rates with identical initial filling, or a decayed-TSC procedure; justify it physically.

## A12. Future prospects: keep them specific

Choose at most two practical consequences, each connected to the central defect-characterization problem. Examples: a better spectroscopy analysis pipeline with uncertainty and competing transport models; measurements that distinguish shallow carrier-separating traps from deep recombination centers; improved process selection for photoelectrodes or thin-film optoelectronic devices. Kaur (2026) is useful to discuss how intentionally tuned sulfur vacancies can help or harm charge transport in **photoelectrochemical** devices, which is not the same experimental geometry as Pai's TSC setup.

An achievable Go extension is to let \(s\), baseline, and a bounded transport exponent vary and fit across multiple heating rates. A more demanding extension is coupled retrapping/recombination rate equations. State clearly which is future work and which code currently exists.

## A13. Conclusion: three scientific claims, three limits

Write your own 65–85 word conclusion after reviewing all figures. It should distinguish (1) **visual merging** of two source components, (2) **conditional statistical recovery** when the forward model is correct, and (3) **failure under model mismatch**. Your limit is equally important: no experiment established that a specific pair of defects in β-In₂S₃ was physically unresolved. Avoid overstating claims of improved solar performance or discovering new states.

## Sources and verification table

| Reference in LaTeX | What it supports | Where to read |
|---|---|---|
| [1] Pai et al. (2005) | β-In₂S₃ sample, heating setup, TSC peaks, assignments, optical gap | PDF pp. 1–5, especially Fig. 1, Table II |
| [2] Pistor et al. (2016) | Crystal structure, tetragonal defect spinel, ordering | Original article and its crystal structure tables |
| [3] Zhao et al. (2012) | Calculated electronic band edge, indirect fundamental gap | Abstract and electronic-structure discussion |
| [4] Dussel and Bube (1967) | TSC physics with retrapping and changing lifetime | PDF pp. 1–2, 8–9, 13–15 |
| [5] Kaur (2026) | Modern defect-engineering application, not TSC measurements | Journal of Alloys and Compounds Communications 11, 100214 |

## Own-words check before submission

- Every `AUTHOR WRITES` prompt has been replaced with your own complete explanation; no bracketed writing directions remain in the final submission.
- Each equation has its variables, assumptions, signs and units explained; figures are cited in your text and their provenance is clear.
- Your 2.6 eV optical gap, calculated indirect gap and thermal activation energies are not incorrectly treated as the same quantity.
- The Go results match JSON/CSV outputs; no numerical count has been manually updated without rebuilding.
- A reader can tell the difference between **experimental results**, **literature interpretation**, **synthetic inputs**, **computed findings**, and **future work**.

## How to use ChatGPT without surrendering authorship

After writing a subsection, you can ask for a **technical referee check without rewriting it**: "Identify factual errors, missing assumptions and unclear logic; do not draft my paragraph." For derivations, ask for a **hint at the step where you are stuck** rather than the full proof. For viva practice, ask for an oral exam that challenges claims from your own sections. In each case, keep the final wording and interpretation yours.
