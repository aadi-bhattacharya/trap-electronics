.PHONY: test data paper all clean

test:
	go test ./... -count=1
	go vet ./...

data:
	go run ./cmd/traplab all configs/experiment.json 200
	go run ./cmd/paperdata

paper:
	inkscape results/demo.svg --export-filename=paper/figures/demo.pdf 2>/dev/null
	inkscape results/resolution.svg --export-filename=paper/figures/resolution.pdf 2>/dev/null
	inkscape results/mismatch_example.svg --export-filename=paper/figures/mismatch_example.pdf 2>/dev/null
	cd paper && pdflatex -interaction=nonstopmode -halt-on-error root.tex >/dev/null && pdflatex -interaction=nonstopmode -halt-on-error root.tex >/dev/null

all: test data paper

clean:
	rm -rf results/*.svg results/*.json results/*.csv paper/figures/*.pdf paper/root.pdf paper/root.aux paper/root.log
