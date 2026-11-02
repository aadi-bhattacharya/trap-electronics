# When Two Traps Look Like One

An author-led semiconductor-physics coursework research toolkit with a reproducible **Go** computation and a **structured IEEE conference paper scaffold**.

The paper is deliberately **not written for the student**. It gives detailed prompts (`AUTHOR WRITES A0-A13`), a complete first-principles companion in `docs/AUTHOR_WRITING_GUIDE.md`, fact-checked source anchors, and verified synthetic numerical results. The student is responsible for all explanation, derivation narration, interpretation, diagrams and conclusions.

## Files to start with

- `paper/root.tex`: editable IEEE manuscript, 13 substantial writing prompts.
- `paper/root.pdf`: compiled scaffold in the original `ieeeconf` layout.
- `docs/AUTHOR_WRITING_GUIDE.md`: paragraph plans, derivations to work through, source pages, and pitfalls.
- `internal/tsc/`: Go first-order kinetic simulator, fitting and model-mismatch experiments.
- `results/`: regenerated SVG charts and numeric CSV/JSON outputs.
- `configs/experiment.json`: the reproducible setup and random seed.

## Build

Requirements: Go, LaTeX (`pdflatex`), and Inkscape.

```sh
make test
make all
```

`make all` regenerates Go data, result macros, PDF vector graphics, and the LaTeX manuscript. The reference measurements in Pai et al. motivate the simulation; this project does **not** fit the authors' original measured TSC datasets or assert that its hypothetical traps have experimentally identified chemical origins.

## Paper writing

Use `docs/AUTHOR_WRITING_GUIDE.md` alongside `paper/root.tex`. Replace each `AUTHOR WRITES` paragraph with original wording before considering the manuscript complete. The two-column IEEE formatting is provided only as a formatting template; it does not imply an IEEE submission or publication.

## Scientific scope

The research question is whether two distinct thermally active trap levels can merge into one visible TSC maximum, whether they remain statistically recoverable, and when model mismatch creates a falsely inferred second trap. This belongs to the course themes of crystal structure, band theory, transport, research challenges, and future prospects.

## Attribution

Bibliographic references and full journal article details are included in `paper/root.tex`. Uploaded source PDFs are not redistributed in this repository. Synthetic Go outputs are explicitly labeled and their amplitudes are normalized, not absolute trap densities.
