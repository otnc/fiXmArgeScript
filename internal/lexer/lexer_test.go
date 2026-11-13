package lexer

import (
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/token"
)

func lexOne(t *testing.T, src string) token.Token {
	t.Helper()
	tokens, _, err := Lex(src, "t.fxma", Draft)
	if err != nil {
		t.Fatalf("Lex(%q): %v", src, err)
	}
	return tokens[0]
}

func TestNumbers(t *testing.T) {
	cases := []struct {
		src  string
		kind token.Kind
		name string
	}{
		{"42", token.Int, "42"},
		{"1_000_000", token.Int, "1000000"},
		{"3.14", token.Float, "3.14"},
		{"1e-3", token.Float, "1e-3"},
		{"2.5E10", token.Float, "2.5E10"},
	}
	for _, c := range cases {
		got := lexOne(t, c.src)
		if got.Kind != c.kind || got.Name != c.name {
			t.Errorf("%q: got %v %q", c.src, got.Kind, got.Name)
		}
	}
}

func TestRangeIsNotFloat(t *testing.T) {
	tokens, _, err := Lex("1..5", "t.fxma", Draft)
	if err != nil {
		t.Fatal(err)
	}
	if tokens[0].Name != "1" || !tokens[1].IsOp("..") || tokens[2].Name != "5" {
		t.Errorf("got %+v", tokens)
	}
}

func TestStrings(t *testing.T) {
	cases := []struct {
		src  string
		want []token.StringPart
	}{
		{`"a\nb"`, []token.StringPart{{Text: "a\nb"}}},
		{`"\u{1F52E}\#{x}"`, []token.StringPart{{Text: "🔮#{x}"}}},
		{`"x は #{x} です"`, []token.StringPart{{Text: "x は "}, {Text: "x", IsCode: true}, {Text: " です"}}},
		{"\"\"\"\n    a\n      b\n    \"\"\"", []token.StringPart{{Text: "a\n  b"}}},
		{`"#{f("}")}"`, []token.StringPart{{Text: `f("}")`, IsCode: true}}},
	}
	for _, c := range cases {
		got := lexOne(t, c.src).Parts
		if len(got) != len(c.want) {
			t.Errorf("%q: got %+v", c.src, got)
			continue
		}
		for i := range got {
			if got[i].Text != c.want[i].Text || got[i].IsCode != c.want[i].IsCode {
				t.Errorf("%q: part %d = %+v, want %+v", c.src, i, got[i], c.want[i])
			}
		}
	}
}

func TestNewlinesInsideBracketsAreIgnored(t *testing.T) {
	tokens, _, err := Lex("[1,\n 2]\nx", "t.fxma", Draft)
	if err != nil {
		t.Fatal(err)
	}
	newlines := 0
	for _, tok := range tokens {
		if tok.Kind == token.Newline {
			newlines++
		}
	}
	if newlines != 1 {
		t.Errorf("改行の字句が %d 個", newlines)
	}
}

// Windows のエディタが頭に付ける BOM は読み飛ばし、位置は1行目の1列目から数える。
func TestLeadingBOMIsSkipped(t *testing.T) {
	tokens, _, err := Lex("\uFEFFx izu 1\n", "main.fxma", Draft)
	if err != nil {
		t.Fatal(err)
	}
	if first := tokens[0]; first.Name != "x" || first.Pos != (token.Pos{Line: 1, Col: 1}) {
		t.Errorf("first token = %+v", first)
	}
}
