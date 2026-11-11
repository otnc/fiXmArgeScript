package runner

import (
	"errors"
	"strings"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
)

// hakStep は fixmarge hak の入力1つと、そのときの出力とエラー。
type hakStep struct {
	text   string
	stdout string
	err    diag.Kind // 空ならエラーなし
	pos    string    // エラーの場所 (hak:行:列)。空なら確かめない
}

func runHak(t *testing.T, name, stdin string, steps []hakStep) {
	t.Helper()
	var stdout strings.Builder
	h := NewHak(interp.Options{Stdin: strings.NewReader(stdin), Stdout: &stdout, Stderr: &stdout, Now: fixedNow, Location: tokyo})
	for i, s := range steps {
		stdout.Reset()
		_, err := h.Feed(s.text)
		de, isDiag := errors.AsType[*diag.Error](err)
		switch {
		case s.err == "" && err != nil:
			t.Errorf("%s #%d: エラー %v", name, i, err)
		case s.err != "" && (!isDiag || de.Kind != s.err):
			t.Errorf("%s #%d: エラー %v, want %s", name, i, err, s.err.Spelling())
		case s.pos != "" && de.File+":"+de.Pos.String() != s.pos:
			t.Errorf("%s #%d: 場所 %s:%s, want %s", name, i, de.File, de.Pos, s.pos)
		}
		if got := stdout.String(); got != s.stdout {
			t.Errorf("%s #%d:\n got %q\nwant %q", name, i, got, s.stdout)
		}
	}
}

func TestHak(t *testing.T) {
	runHak(t, "carry over", "", []hakStep{
		{text: "x izu 1\n"},
		{text: "x +izu 1\n"},
		{text: "x redogh\n", stdout: "2\n"},
	})
	runHak(t, "echo", "", []hakStep{
		{text: "1 + 2\n", stdout: "3\n"},
		{text: "\"ab\"\n", stdout: "\"ab\"\n"},
		{text: "vanita\n"},
		{text: "\"ab\" redogh\n", stdout: "ab\n"},
		{text: "x izu 1\nx\n", stdout: "1\n"},
	})
	runHak(t, "block", "", []hakStep{
		{text: "weho f(n)\n  n * 2 itswa\nfinend\n"},
		{text: "f(4)\n", stdout: "8\n"},
	})
	runHak(t, "later definition", "", []hakStep{
		{text: "weho f()\n  g() itswa\nfinend\n"},
		{text: "weho g()\n  1 itswa\nfinend\n"},
		{text: "f()\n", stdout: "1\n"},
	})
	runHak(t, "redefine", "", []hakStep{
		{text: "weho f()\n  1 itswa\nfinend\nshowence c izu 1\n"},
		{text: "weho f()\n  2 itswa\nfinend\nshowence c izu 2\n"},
		{text: "[f(), c]\n", stdout: "[2, 2]\n"},
		{text: "weho g()\nfinend\nweho g()\nfinend\n", err: diag.SuReady},
		{text: "c izu 3\n", err: diag.BlasTOac},
		{text: "\":meme\" aste m werer\n"},
		{text: "\":coword\" aste m werer\n"},
	})
	runHak(t, "unfinished", "", []hakStep{
		{text: "mof tadact\n", err: diag.AyAmis},
		{text: "\"\"\"\n", err: diag.AyAmis},
	})
	runHak(t, "error is undone", "", []hakStep{
		{text: "x izu 1\n"},
		{text: "x izu 2\nx redogh\ncnt redogh\n", stdout: "2\n", err: diag.NiiExi, pos: "hak:4:1"},
		{text: "x\n", stdout: "1\n"},
		{text: "1 +\n", err: diag.AyAmis, pos: "hak:6:4"},
		{text: "y redogh\n", err: diag.NiiExi, pos: "hak:7:1"},
	})
	runHak(t, "backtrack within entry", "", []hakStep{
		{text: "n izu 0\n"},
		{text: "x izu 1 owk 2 owk 3 mekamo\nn +izu 1\nx senaji 3 shoubek\n"},
		{text: "[x, n]\n", stdout: "[3, 1]\n"},
	})
	runHak(t, "earlier choices are fixed", "", []hakStep{
		{text: "x izu 1 owk 2 mekamo\n"},
		{text: "x senaji 2 shoubek\n", err: diag.CoJundic},
		{text: "x\n", stdout: "1\n"},
	})
	runHak(t, "co-jundic is undone", "", []hakStep{
		{text: "x izu 0\n"},
		{text: "x izu 1 owk 2 mekamo\nx senaji 3 shoubek\n", err: diag.CoJundic},
		{text: "x\n", stdout: "0\n"},
	})
	// 捨てた入力が読んだ行は戻らない (外の世界とのやり取りなので)
	runHak(t, "input", "a\nb\nc\n", []hakStep{
		{text: "p izu zanfession\n"},
		{text: "zanfession\ncnt redogh\n", err: diag.NiiExi},
		{text: "q izu zanfession\n"},
		{text: "[p, q]\n", stdout: "[\"a\", \"c\"]\n"},
	})
	runHak(t, "rollback keeps later output", "", []hakStep{
		{text: "1 redogh\n", stdout: "1\n"},
		{text: "2 redogh\ncnt\n", stdout: "2\n", err: diag.NiiExi},
		{text: "3 redogh\n", stdout: "3\n"},
	})
}

func TestHakUnfinished(t *testing.T) {
	cases := []struct {
		text       string
		unfinished bool
	}{
		{"mof tadact\n", true},
		{"weho f()\n  1 itswa\n", true},
		{"f(1,\n", true},
		{"\"\"\"\nab\n", true},
		{"\"ab\n", false},
		{"1 + )\n", false},
		{"mof tadact\nfinend\n", false},
	}
	for _, c := range cases {
		_, err := NewHak(interp.Options{Stdout: new(strings.Builder)}).Feed(c.text)
		if got := IsUnfinished(err); got != c.unfinished {
			t.Errorf("%q: unfinished %v, want %v (%v)", c.text, got, c.unfinished, err)
		}
	}
}
