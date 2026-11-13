package runner

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/build"
	"github.com/otnc/fiXmArgeScript/internal/format"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/lint"
	"github.com/otnc/fiXmArgeScript/internal/scope"
)

// neverEnding は、実行が終わらないかもしれないソースに出てくる語。繰り返し、関数 (再帰)、待つこと、読み込み (標準パッケージ)。
var neverEnding = []string{"azagara", "weho", "praylee", "zanfession", "werer", "curoi"}

// FuzzRun は、必ず終わるソースならどれを実行しても落ちないことを確かめる。
// 探すときは go test -fuzz FuzzRun ./internal/runner/。
func FuzzRun(f *testing.F) {
	f.Add("x izu 1 owk 2 owk 3 mekamo\nx > 1 shoubek\n[x, \"#{x * 2}\", {\"k\": x}] redogh\n")
	f.Add("xs izu [1, 2, 3]\nxs no 0..1 izu [9]\nxs no -1 redogh\n(2 exi 0..10 batte 2) redogh\n\"ab\" * 3 no 1.. redogh\n")
	f.Add("teitance\n  1 // 0\nyoulugive e\n  e no \"chidd\" canosiri\nennote\n  -7 % 3 << 2 redogh\nfinend\n")
	f.Add("showence d izu {\"a\": [1.5, vanita, tadact]}\nd senaji d redogh\nnang(sanand(\"ab\", \"hex\")) redogh\n")

	f.Fuzz(func(t *testing.T, src string) {
		for _, word := range neverEnding {
			if strings.Contains(strings.ToLower(src), word) {
				t.Skip()
			}
		}
		Run(src, "main.fxma", nil, interp.Options{Stdin: strings.NewReader(""), Stdout: io.Discard, Stderr: io.Discard})
	})
}

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
			lint.File(file, src, nil, nil)
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
