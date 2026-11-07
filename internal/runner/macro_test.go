package runner

import (
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

func TestMacros(t *testing.T) {
	cases := []testCase{
		{name: "spec example", src: "curoi twice(body)\n  body\n  body\nfinend\n\ntwice(\"Mof\" redogh)\n", stdout: "Mof\nMof\n"},
		{name: "hoisted", src: "twice(1 redogh)\ncuroi twice(body)\n  body\n  body\nfinend\n", stdout: "1\n1\n"},
		{name: "expression argument is grouped", src: "curoi square(x)\n  x * x\nfinend\nsquare(1 + 2) redogh\n", stdout: "9\n"},
		{name: "in expression", src: "curoi square(x)\n  x * x\nfinend\ny izu square(3) + 1\ny redogh\n", stdout: "10\n"},
		{name: "hygiene", src: "curoi swap(a, b)\n  tmp izu a\n  a izu b\n  b izu tmp\nfinend\ntmp izu 1\nother izu 2\nswap(tmp, other)\n\"#{tmp} #{other}\" redogh\n", stdout: "2 1\n"},
		{name: "macro calls macro", src: "curoi twice(body)\n  body\n  body\nfinend\ncuroi four(body)\n  twice(twice(body))\nfinend\nfour(\"x\" redogh)\n", stdout: "x\nx\nx\nx\n"},
		{name: "loop variable renamed", src: "curoi each(xs, body)\n  xs aste it azagara\n    body\n  finend\nfinend\nit izu \"outer\"\neach([1, 2], it redogh)\n", stdout: "outer\nouter\n"},
		{name: "in interpolation", src: "curoi square(x)\n  x * x\nfinend\n\"#{square(4)}\" redogh\n", stdout: "16\n"},
		{name: "too deep", src: "curoi loop(x)\n  loop(x)\nfinend\nloop(1)\n", err: diag.Bavellabion},
		{name: "wrong arguments", src: "curoi twice(body)\n  body\nfinend\ntwice(1, 2)\n", err: diag.AyAmis},
		{name: "duplicate", src: "curoi a(x)\n  x\nfinend\ncuroi a(x)\n  x\nfinend\n", err: diag.SuReady},
		{name: "module member is not a macro", src: "\"./util\" aste util werer\ncuroi double(x)\n  0\nfinend\nutil no double(21) redogh\n", files: map[string]string{"util.fxma": utilSrc}, stdout: "42\n"},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}
