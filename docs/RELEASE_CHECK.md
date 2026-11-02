# Author-guided manuscript release check

- [x] New Go standard-library implementation with separate model, inverse fitting, simulation runs, and SVG plots.
- [x] No external Go packages and no Python or notebook dependency.
- [x] 200 synthetic trials for each of five trap separations and three one-trap stress conditions.
- [x] All numerical results regenerated at fixed seed 20261102; `cmd/paperdata` synchronizes manuscript macros with JSON.
- [x] First-order peak location, carrier-release integral, component-additivity, fine-grid comparison, inverse fitting and nuisance-offset fit checked in Go tests.
- [x] `go test ./...`, `go vet ./...` and `go test -race ./...` passed at packaging.
- [x] IEEE class supplied in the user template preserved; PDF generated with pdflatex.
- [x] All figure plots created by Go and converted to vector PDF.
- [x] 14 `AUTHOR WRITES` locations (A0--A13) remain for student-authored discussion and theory.
- [x] Material facts and spectroscopy claims supported by primary sources / verified publisher pages.
- [x] Paper explicitly separates reported source parameters, model assumptions and synthetic results.
- [ ] Student has independently written abstract, theory, interpretation, future prospects and conclusion.
- [ ] Student has drawn own crystal-structure and band diagrams.
- [ ] Student has checked final instructor-required page count, author metadata, citations and page breaks.
- [ ] Student has added GitHub repository URL / commit hash to final submission.

**Caveat:** These tests do not establish that a particular physical pair of β-In₂S₃ defects is present in a specimen, nor a universal trap-resolution threshold. The model neglects retrapping and temperature-dependent carrier kinetics in its base family. The stress tests are synthetic sensitivity analyses.
