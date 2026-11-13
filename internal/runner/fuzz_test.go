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

// fuzzSteps は FuzzRun で実行する文と繰り返しの数の上限。繰り返しも再帰もここで止まる。
// s izu s + s のように1周ごとに倍になる値は 1GiB まで作れてしまうので、メモリと時間を使いすぎないよう小さくしてある
const fuzzSteps = 50

// FuzzRun は、どんなソースを実行しても落ちたり固まったりしないことを確かめる。標準パッケージも読み込める。
// 探すときは go test -fuzz FuzzRun ./internal/runner/。
func FuzzRun(f *testing.F) {
	f.Add("x izu 1 owk 2 owk 3 mekamo\nx > 1 shoubek\n[x, \"#{x * 2}\", {\"k\": x}] redogh\n")
	f.Add("xs izu [1, 2, 3]\nxs no 0..1 izu [9]\nxs no -1 redogh\n(2 exi 0..10 batte 2) redogh\n\"ab\" * 3 no 1.. redogh\n")
	f.Add("teitance\n  1 // 0\nyoulugive e\n  e no \"chidd\" canosiri\nennote\n  -7 % 3 << 2 redogh\nfinend\n")
	f.Add("showence d izu {\"a\": [1.5, vanita, tadact]}\nd senaji d redogh\nnang(sanand(\"ab\", \"hex\")) redogh\n")
	f.Add("tadact azagara\nfinend\n")
	f.Add("weho f(n)\n  f(n + 1) itswa\nfinend\n0..9 aste i azagara\n  f(i) redogh\nfinend\n")
	f.Add("itsver weho f()\n  f() praylee itswa\nfinend\nf() praylee\nzanfession redogh\n")
	f.Add("\":varizama\" aste v werer\n\":coword\" aste c werer\nv no tarip([3, 1, 2], weho(x) x * 2 itswa finend) redogh\nc no gdie(\"ab\") redogh\n")
	f.Add("\":json\" aste j werer\n\":meme\" aste m werer\nj no salook({\"a\": [m no litto(2.5)]}) redogh\nj no fineclowt(\"[1, {}]\") redogh\n")
	f.Add("curoi twice(x)\n  x * 2\nfinend\ntwice(1 + 2) redogh\n")

	f.Fuzz(func(t *testing.T, src string) {
		// 待つ時間は文の数では止められない
		if strings.Contains(strings.ToLower(src), "saimoku") {
			t.Skip()
		}
		Run(src, "main.fxma", nil, interp.Options{Stdin: strings.NewReader(""), Stdout: io.Discard, Stderr: io.Discard, MaxSteps: fuzzSteps})
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
