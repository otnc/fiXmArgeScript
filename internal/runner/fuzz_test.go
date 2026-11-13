package runner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/build"
	"github.com/otnc/fiXmArgeScript/internal/format"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/lint"
	"github.com/otnc/fiXmArgeScript/internal/scope"
)

// FuzzFrontend は、どんなソースを渡しても実行前の処理 (読む、検査、整形、清書、lint) が落ちないことを確かめる。
// 探すときは go test -fuzz FuzzFrontend ./internal/runner/。普段の go test では下の種だけを通す。
func FuzzFrontend(f *testing.F) {
	examples, _ := filepath.Glob("../../examples/*.fxma")
	for _, name := range examples {
		src, err := os.ReadFile(name)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(src))
	}
	for _, src := range benchPrograms {
		f.Add(src)
	}
	f.Add("curoi twice(x)\n  x * 2\nfinend\ntwice(1 + 2) redogh\n\"\"\"\n  a #{1} b\n\"\"\" redogh\n")
	f.Add("teitance\n  \"x\" accuhearts\nyoulugive e\n  e no \"lyuzun\" canosiri\nennote\n  1\nfinend\n")
	f.Add("itsver weho f()\n  1 praylee itswa\nfinend\nxs izu [1, 2] no 0..=1\nd izu {\"a\": 1}\n")

	f.Fuzz(func(t *testing.T, src string) {
		file, err := Prepare(src, "main.fxma", nil)
		if file != nil {
			scope.Resolve(file.Stmts)
			lint.File(file, src, nil)
		}
		formatted, ferr := format.Format(src, "main.fxma")
		if ferr == nil && err == nil {
			// 整えた結果は、もう一度整えても変わらない
			again, err := format.Format(formatted, "main.fxma")
			if err != nil || again != formatted {
				t.Errorf("整形が落ち着かない (%v)\n1回目 %q\n2回目 %q", err, formatted, again)
			}
		}
		fair, _, berr := build.Build(src, "main.fxma")
		if berr == nil {
			// 清書は清書として読める
			if _, _, err := lexer.Lex(fair, "main.fxmac", lexer.Fair); err != nil {
				t.Errorf("清書が読めない: %v\n%q", err, fair)
			}
		}
	})
}
