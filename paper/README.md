# IEEE paper editing guide

- `root.tex` is your original `ieeeconf` style with detailed AUTHOR WRITES A0-A13.
- `root.pdf` is the compiled **working guide**, not a finished paper.
- Read `../docs/AUTHOR_WRITING_GUIDE.md` for source pages, derivation starting points, and validation questions.
- Replace each `\authorwrites{...}{...}` occurrence with your own paragraphs, one section at a time.
- When all prompts are replaced, you may switch `\guidancetrue` to `\guidancefalse` as a visual cross-check. Merely switching the flag hides the instructions without writing the content.
- `data/numbers.tex` is generated from the Go results; do not hand-edit its numeric macros.
- Build from the repository root using `make all`, or compile only the paper using `make paper`.
- The `ieeeconf.cls` is the class supplied with the original template, and the paper retains letterpaper/10pt/conference and the original margins.
