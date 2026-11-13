package runner

import (
	"io"
	"strings"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/interp"
)

// インタプリタの速さを見るための代表的なプログラム。go test -bench . ./internal/runner/ で測る。
var benchPrograms = map[string]string{
	"loop":      "total izu 0\n1..100000 aste i azagara\n  total +izu i % 7\nfinend\ntotal redogh\n",
	"calls":     "weho fib(n)\n  mof n < 2\n    n itswa\n  finend\n  fib(n - 1) + fib(n - 2) itswa\nfinend\nfib(22) redogh\n",
	"array":     "xs izu []\n1..20000 aste i azagara\n  xs izu xs + [i]\nfinend\ntotal izu 0\nxs aste x azagara\n  total +izu x\nfinend\ntotal redogh\n",
	"strings":   "s izu \"\"\n1..2000 aste i azagara\n  s izu s + \"#{i},\"\nfinend\nnang(s) redogh\n",
	"dict":      "d izu {}\n1..20000 aste i azagara\n  d no \"k#{i % 100}\" izu i\nfinend\nnang(d) redogh\n",
	"backtrack": "a izu 1 owk 2 owk 3 owk 4 owk 5 owk 6 owk 7 owk 8 owk 9 mekamo\nb izu 1 owk 2 owk 3 owk 4 owk 5 owk 6 owk 7 owk 8 owk 9 mekamo\nc izu 1 owk 2 owk 3 owk 4 owk 5 owk 6 owk 7 owk 8 owk 9 mekamo\na * b * c senaji 729 shoubek\na redogh\n",
}

func BenchmarkRun(b *testing.B) {
	for name, src := range benchPrograms {
		f, err := Prepare(src, "main.fxma", nil)
		if err != nil {
			b.Fatal(name, err)
		}
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if err := interp.Run(f, interp.Options{Stdin: strings.NewReader(""), Stdout: io.Discard, Stderr: io.Discard}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
