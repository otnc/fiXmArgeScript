package format

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormat(t *testing.T) {
	cases := []struct{ name, in, want string }{
		{"indent", "weho f(n)\nmof n > 0\nn itswa\nsiva mof n senaji 0\n0 itswa\nsiva\n-n itswa\nfinend\nfinend\n", "weho f(n)\n  mof n > 0\n    n itswa\n  siva mof n senaji 0\n    0 itswa\n  siva\n    -n itswa\n  finend\nfinend\n"},
		{"spaces in line", "x izu f( a,b )+1\nd izu {\"k\":[1,2]}\n", "x izu f(a, b) + 1\nd izu {\"k\": [1, 2]}\n"},
		{"unary and ranges", "y izu xs no -1\nz izu 2 ** -1\nxs no 2 .. 4 redogh\nxs no .. 1 redogh\nxs no 4.. redogh\n0..10 batte 5 aste i azagara\nfinend\n", "y izu xs no -1\nz izu 2 ** -1\nxs no 2..4 redogh\nxs no ..1 redogh\nxs no 4.. redogh\n0..10 batte 5 aste i azagara\nfinend\n"},
		{"compound assign", "x + izu 1\nx *izu 2\n", "x +izu 1\nx *izu 2\n"},
		{"calls and keywords", "weho(n) n itswa finend\ndv no \"murm\" (w)\nnii(x) redogh\n", "weho (n) n itswa finend\ndv no \"murm\"(w)\nnii (x) redogh\n"},
		{"comments and blank lines", "#!/usr/bin/env fixmarge\n\n\n\n# a\nx izu 1   # b   \n\n\n", "#!/usr/bin/env fixmarge\n\n# a\nx izu 1  # b\n"},
		{"comment inside block", "mof x\n# c\ny redogh\nfinend\n", "mof x\n  # c\n  y redogh\nfinend\n"},
		{"multiline brackets", "xs izu [\n1, # one\n\n  2\n]\n", "xs izu [\n  1,  # one\n  2\n]\n"},
		{"callback", "f(1, weho (req)\nreq itswa\nfinend)\n", "f(1, weho (req)\n  req itswa\nfinend)\n"},
		{"try", "teitance\nx\nyoulugive e\ny\nennote\nz\nfinend\n", "teitance\n  x\nyoulugive e\n  y\nennote\n  z\nfinend\n"},
		{"closure returned", "weho counter()\nweho ()\n1 itswa\nfinend itswa\nfinend\n", "weho counter()\n  weho ()\n    1 itswa\n  finend itswa\nfinend\n"},
		{"strings kept", "s izu \"\"\"\n    a\n      b #{x+1}\n    \"\"\"\n", "s izu \"\"\"\n    a\n      b #{x+1}\n    \"\"\"\n"},
		{"empty", "\n\n", ""},
	}
	for _, c := range cases {
		got, err := Format(c.in, "x.fxma")
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.want)
		}
		if again, _ := Format(got, "x.fxma"); again != got {
			t.Errorf("%s: 2回目で変わった %q", c.name, again)
		}
	}
}

// サンプルはもう整った形で書いてあるので、整えても変わらない。
func TestExamplesAreFormatted(t *testing.T) {
	files, _ := filepath.Glob("../../examples/**/*.fxma")
	top, _ := filepath.Glob("../../examples/*.fxma")
	for _, path := range append(top, files...) {
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Format(string(src), path)
		if err != nil {
			t.Errorf("%s: %v", path, err)
		} else if got != string(src) {
			t.Errorf("%s:\n got %q\nwant %q", path, got, src)
		}
	}
}

func TestExpand(t *testing.T) {
	src := "curoi swap(a, b)\n  tmp izu a\n  a izu b\n  b izu tmp\nfinend\nswap(x, y + 1)\n"
	got, err := Expand(src, "x.fxma")
	if err != nil {
		t.Fatal(err)
	}
	want := "tmp·1 izu x\nx izu (y + 1)\n(y + 1) izu tmp·1\n"
	if got != want {
		t.Errorf("got %q\nwant %q", got, want)
	}
}
