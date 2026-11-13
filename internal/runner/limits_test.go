package runner

import (
	"errors"
	"strings"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/format"
)

// nest は open を n 回、inner、close を n 回の順につなぐ。
func nest(n int, open, inner, close string) string {
	return strings.Repeat(open, n) + inner + strings.Repeat(close, n)
}

// 深すぎる入れ子は、Go のスタックを使い切る前に Bavellabion で止める。
func TestDeepSourceIsOverflow(t *testing.T) {
	const deep = 50000
	cases := []testCase{
		{name: "shallow enough", src: nest(500, "(", "1", ")") + " redogh\n", stdout: "1\n"},
		{name: "parens", src: nest(deep, "(", "1", ")") + " redogh\n", err: diag.Bavellabion},
		{name: "arrays", src: nest(deep, "[", "1", "]") + " redogh\n", err: diag.Bavellabion},
		{name: "signs", src: strings.Repeat("-", deep) + "1 redogh\n", err: diag.Bavellabion},
		{name: "nii", src: strings.Repeat("nii ", deep) + "tadact redogh\n", err: diag.Bavellabion},
		{name: "operator chain", src: "1" + strings.Repeat(" + 1", deep) + " redogh\n", err: diag.Bavellabion},
		{name: "call chain", src: "f" + strings.Repeat("()", deep) + "\n", err: diag.Bavellabion},
		// mekamo の候補は並びとして持つので、いくつ並べても深くならない
		{name: "candidates", src: "x izu 1" + strings.Repeat(" owk 2", deep) + " mekamo\nx redogh\n", stdout: "1\n"},
		{name: "blocks", src: nest(deep, "mof tadact\n", "1 redogh\n", "finend\n"), err: diag.Bavellabion},
		{name: "strings in strings", src: nest(500, "\"#{", "1", "}\"") + " redogh\n", err: diag.Bavellabion},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

// 本体の深い関数が深く呼び合うと、呼び出しの数が上限より少なくても Bavellabion で止まる。
func TestDeepBodiesLowerTheCallLimit(t *testing.T) {
	loops := nest(300, "[1] aste x azagara\n", "f(n + 1)\n", "finend\n")
	cases := []testCase{
		{name: "plain recursion", src: "weho f(n)\n  mof n senaji 5000\n    n itswa\n  finend\n  f(n + 1) itswa\nfinend\nf(0) redogh\n", stdout: "5000\n"},
		{name: "endless recursion", src: "weho f(n)\n  f(n + 1) itswa\nfinend\nf(0)\n", err: diag.Bavellabion},
		{name: "deep body", src: "weho f(n)\n" + loops + "finend\nf(0)\n", err: diag.Bavellabion},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

// 深く入れ子にした値は、出しても比べても落ちない。
func TestDeepValues(t *testing.T) {
	const build = "xs izu []\nys izu []\n0...300000 aste n azagara\n  xs izu [xs]\n  ys izu {\"k\": ys}\nfinend\n"
	cases := []testCase{
		{name: "display", src: build + "nang(\"#{xs}\") redogh\nnang(\"#{ys}\") redogh\n", stdout: "600002\n2100002\n"},
		{name: "equal", src: build + "zs izu []\n0...300000 aste n azagara\n  zs izu [zs]\nfinend\nxs senaji zs redogh\nxs senaji [zs] redogh\nxs senaji ys redogh\n", stdout: "TadACt\nULieSo\nULieSo\n"},
		{name: "json", src: "\":json\" aste json werer\n" + build + "json no salook(xs)\n", err: diag.ImJune},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

// 範囲の端が整数の端にあっても、繰り返しは端で止まる。
func TestRangeAtIntegerLimits(t *testing.T) {
	cases := []testCase{
		{name: "up", src: "9223372036854775800..9223372036854775807 batte 3 aste i azagara\n  i redogh\nfinend\n", stdout: "9223372036854775800\n9223372036854775803\n9223372036854775806\n"},
		{name: "down", src: "lo izu -9223372036854775807 - 1\n-9223372036854775800..lo batte -4 aste i azagara\n  i redogh\nfinend\n", stdout: "-9223372036854775800\n-9223372036854775804\n-9223372036854775808\n"},
		{name: "exi", src: "r izu -9223372036854775807..9223372036854775807 batte 2\n(9223372036854775807 exi r) redogh\n(9223372036854775806 exi r) redogh\n", stdout: "TadACt\nULieSo\n"},
		{name: "senchoose", src: "\":destimay\" aste d werer\nr izu -9223372036854775807..9223372036854775807\n(d no senchoose(r) exi r) redogh\nd no senchoose(5...5)\n", stdout: "TadACt\n", err: diag.NiiExi},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

// 引数を何度も使う呪いを入れ子で呼んでも、展開が膨れ上がる前に Bavellabion で止まる。
func TestMacroExpansionIsBounded(t *testing.T) {
	wide := "curoi wide(x)\n  [x, x, x, x, x, x, x, x, x, x]\nfinend\n"
	cases := []testCase{
		{name: "small", src: wide + "nang(wide(wide(1))) redogh\n", stdout: "10\n"},
		{name: "explosive", src: wide + "nang(" + nest(9, "wide(", "1", ")") + ") redogh\n", err: diag.Bavellabion},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

// 整形は、字下げだけで出力が膨れ上がるほど深い入れ子を Bavellabion で断る。
func TestFormatRefusesDeepNesting(t *testing.T) {
	sources := map[string]string{
		"blocks":   nest(20000, "mof tadact\n", "1 redogh\n", "finend\n"),
		"brackets": "x izu " + nest(20000, "[\n", "1", "\n]") + "\n",
	}
	for name, src := range sources {
		_, err := format.Format(src, "main.fxma")
		if de, ok := errors.AsType[*diag.Error](err); !ok || de.Kind != diag.Bavellabion {
			t.Errorf("%s: エラー %v, want Bavellabion", name, err)
		}
	}
}
