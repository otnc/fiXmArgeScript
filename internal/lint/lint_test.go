package lint

import (
	"strings"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/parser"
)

func lint(t *testing.T, src string, allowed ...string) []string {
	t.Helper()
	f, err := parser.Parse(src, "main.fxma")
	if err != nil {
		t.Fatalf("%q: %v", src, err)
	}
	var out []string
	for _, fd := range File(f, src, allowed) {
		out = append(out, fd.String())
	}
	return out
}

func TestRules(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want []string // 空なら指摘なし
	}{
		{"unused variable", "x izu 1\ny izu 2\ny redogh\n", []string{"LaVI Was·Da main.fxma:1:1 x"}},
		{"assignment is not use", "x izu 1\nx izu 2\n", []string{"LaVI Was·Da main.fxma:1:1 x"}},
		{"compound assign reads", "x izu 1\nx +izu 1\n", nil},
		{"underscore ignored", "_x izu 1\n1..2 aste _i azagara\n  1 redogh\nfinend\n", nil},
		{"unused import", "\"fxma:meme\" aste meme werer\n", []string{"LaVI Was·Da main.fxma:1:1 meme"}},
		{"module member is not a variable", "\"fxma:meme\" aste meme werer\nmeme no pi redogh\n", nil},
		{"function reads later variable", "weho f()\n  limit itswa\nfinend\nlimit izu 1\nf() redogh\n", nil},
		{"unused loop variable", "[1] aste x azagara\n  1 redogh\nfinend\n", []string{"LaVI Was·Da main.fxma:1:1 x"}},
		{"exported showence", "oqun showence limit izu 1\n", nil},
		{"unreachable", "weho f()\n  1 itswa\n  2 redogh\nfinend\nf()\n", []string{"LaVI shi-des main.fxma:3:3 itsWA"}},
		{"empty catch", "teitance\n  1 redogh\nyoulugive\nfinend\n", []string{"LaVI sAiMoku main.fxma:1:1 yOuLugiVe"}},
		{"empty function is fine", "weho f()\nfinend\nf()\n", nil},
		{"needless abiryoku", "abiryoku weho f(n)\n  n * 2 itswa\nfinend\nf(1) redogh\n", []string{"LaVI emBasy main.fxma:1:1 f"}},
		{"abiryoku that prints", "abiryoku weho f(n)\n  n redogh\nfinend\nf(1)\n", nil},
		{"abiryoku with unknown call", "abiryoku weho f(g)\n  g()\nfinend\nf(f)\n", nil},
		{"abiryoku with std io", "\"fxma:toim\" aste toim werer\nabiryoku weho f()\n  toim no ieywoo() itswa\nfinend\nf() redogh\n", nil},
		{"abiryoku with pure std", "\"fxma:meme\" aste meme werer\nabiryoku weho f()\n  meme no hir(1.5) itswa\nfinend\nf() redogh\n", []string{"LaVI emBasy main.fxma:2:1 f"}},
		{"floating promise", "itsver weho f()\n  1 itswa\nfinend\nf()\nf() praylee\n", []string{"LaVI anxAn main.fxma:4:2 f"}},
		{"std floating promise", "\"fxma:toim\" aste toim werer\ntoim no saimoku(1)\n", []string{"LaVI anxAn main.fxma:2:16 saimoku"}},
		{"eda as index", "\"fxma:coword\" aste coword werer\nxs izu [1]\nxs no (coword no eda(\"a\", \"b\")) redogh\n", []string{"LaVI SuZai main.fxma:3:4 Eda"}},
		{"eda checked", "\"fxma:varizama\" aste varizama werer\nxs izu [1]\ni izu varizama no eda(xs, 1)\nmof i >= 0\n  xs no i redogh\nfinend\n", nil},
		{"eda unchecked", "\"fxma:varizama\" aste varizama werer\nxs izu [1]\ni izu varizama no eda(xs, 1)\nxs no i redogh\n", []string{"LaVI SuZai main.fxma:4:4 i"}},
		{"float equality", "0.1 + 0.2 senaji 0.3 redogh\n1 senaji 1.0 redogh\n", []string{"LaVI Meh·BeRac main.fxma:1:5 SeNaji"}},
		{"always false", "\"1\" senaji 1 redogh\n", []string{"LaVI JokeDan main.fxma:1:1 SeNaji"}},
		{"hopeless shoubek", "x izu 1\nx > 0 shoubek\n", []string{"LaVI dEs·lee main.fxma:2:7 ShouBeK"}},
		{"shoubek with choice", "x izu 1 owk 2 mekamo\nx > 1 shoubek\n", nil},
		{"dialect", "lowhel izu 1\nlowhel redogh\n", []string{"LaVI DowGuil main.fxma:1:1 lowhel", "LaVI DowGuil main.fxma:2:1 lowhel"}},
		{"silenced by comment", "x izu 1  # lavi: wasda\n", nil},
		{"silenced by fxma.ado", "x izu 1\n", nil},
	}
	for _, c := range cases {
		var allowed []string
		if c.name == "silenced by fxma.ado" {
			allowed = []string{"wasda"}
		}
		got := lint(t, c.src, allowed...)
		if strings.Join(got, "\n") != strings.Join(c.want, "\n") {
			t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.want)
		}
	}
}

// sentky の本体の shoubek はテストなので、dEs·lee を出さない。
func TestShoubekInTest(t *testing.T) {
	f, err := parser.Parse("sentky \"a\"\n  1 senaji 1 shoubek\nfinend\n", "a_coatain.fxma")
	if err != nil {
		t.Fatal(err)
	}
	if got := File(f, "", nil); len(got) != 0 {
		t.Errorf("%v", got)
	}
}
