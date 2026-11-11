package build

import (
	"errors"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
)

func TestBuild(t *testing.T) {
	cases := []struct{ draft, fair string }{
		{"x izu 1\n", "x Iz'u 1\n"},
		{"weho fizzbuzz(n)\n  mof n % 15 senaji 0\n    \"FizzBuzz\" itswa\n", "We-Ho fizzbuzz(n)\n  Mof n % 15 SeNaji 0\n    \"FizzBuzz\" itsWA\n"},
		{"x * izu 2  # 2倍\n", "x *Iz'u 2  # 2倍\n"},
		{"\"x は #{x exi} です、mof\"\n", "\"x は #{x eXi'} です、mof\"\n"},
		{"e no \"lyuzun\" canosiri\n", "e N'o \"lyuzun\" CaNoSiri\n"},
		{"izanavited no izuDefutai(p)\n", "i·zana-vited N'o Iz'udeFuTai(p)\n"},
		{"#!/usr/bin/env fixmarge\nsiva mof tadact\n", "#!/usr/bin/env fixmarge\nSiVa Mof TadACt\n"},
	}
	for _, c := range cases {
		got, _, err := Build(c.draft, "main.fxma")
		if err != nil {
			t.Errorf("Build(%q): %v", c.draft, err)
			continue
		}
		if got != c.fair {
			t.Errorf("Build(%q)\n got %q\nwant %q", c.draft, got, c.fair)
		}
	}
}

// 清書をもう一度読んでも、下書きと同じ字句の並びになること。
func TestFairRoundTrip(t *testing.T) {
	draft := "x izu 1 owk 2 mekamo\nnii x senaji 1 shoubek\nx +izu 1\n\"#{x} desu\" redogh\nizanavited no izuDefutai(\"/\")\n"
	fair, _, err := Build(draft, "main.fxma")
	if err != nil {
		t.Fatal(err)
	}
	a, _, err := lexer.Lex(draft, "main.fxma", lexer.Draft)
	if err != nil {
		t.Fatal(err)
	}
	b, _, err := lexer.Lex(fair, "main.fxmac", lexer.Fair)
	if err != nil {
		t.Fatalf("清書が読めない: %v\n%s", err, fair)
	}
	if len(a) != len(b) {
		t.Fatalf("字句の数が違う: %d と %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Kind != b[i].Kind || a[i].Name != b[i].Name {
			t.Errorf("%d 番目の字句が違う: %+v と %+v", i, a[i], b[i])
		}
	}
}

func TestFairRejectsSkeleton(t *testing.T) {
	_, _, err := lexer.Lex("x izu 1\n", "main.fxmac", lexer.Fair)
	if de, ok := errors.AsType[*diag.Error](err); !ok || de.Kind != diag.AyAmis {
		t.Errorf("清書の骨格は AyAmis になるはず: %v", err)
	}
}

func TestDialectWarning(t *testing.T) {
	got, warnings, err := Build("lowhel redogh\n", "main.fxma")
	if err != nil {
		t.Fatal(err)
	}
	if got != "LowhEll ReDogh'\n" || len(warnings) != 1 || warnings[0].String() != "DowGuil lowhel → LowhEll" {
		t.Errorf("got %q, %v", got, warnings)
	}
}

func TestMixedCaseKeywordIsError(t *testing.T) {
	for _, src := range []string{"IZU\n", "Izu\n", "x izu 1 == 1\n"} {
		_, _, err := Build(src, "main.fxma")
		if de, ok := errors.AsType[*diag.Error](err); !ok || de.Kind != diag.AyAmis {
			t.Errorf("Build(%q) = %v, want AyAmis", src, err)
		}
	}
}
