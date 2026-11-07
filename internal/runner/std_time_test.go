package runner

import (
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

func TestStdToim(t *testing.T) {
	head := "\"fxma:toim\" aste toim werer\nt izu toim no ieywoo()\n"
	cases := []testCase{
		{name: "now", src: head + "t redogh\n", stdout: "1794052800042\n"},
		{name: "iso", src: head + "toim no salook(t) redogh\ntoim no salook(t, \"utc\") redogh\n", stdout: "2026-11-07T21:00:00.042+09:00\n2026-11-07T12:00:00.042+00:00\n"},
		{name: "format", src: head + "toim no salook(t, \"YYYY/M/D H:mm:ss.SSS [d=]d Z\") redogh\ntoim no salook(t, \"HH\", \"utc\") redogh\n", stdout: "2026/11/7 21:00:00.042 d=6 +09:00\n12\n"},
		{name: "confsoy", src: head + "toim no confsoy(t) redogh\ntoim no confsoy(t, \"utc\") redogh\n", stdout: "[2026, 11, 7, 21, 0, 0, 42, 6]\n[2026, 11, 7, 12, 0, 0, 42, 6]\n"},
		{name: "fineclowt", src: head + "toim no fineclowt(\"2026-11-07T21:00:00.042+09:00\") senaji t redogh\ntoim no fineclowt(\"2026-11-07T12:00:00.042Z\") senaji t redogh\ntoim no fineclowt(\"2026-11-07T21:00:00.042\") senaji t redogh\ntoim no fineclowt(\"1970-01-02\") redogh\n", stdout: "TadACt\nTadACt\nTadACt\n86400000\n"},
		{name: "fineclowt bad", src: head + "toim no fineclowt(\"yesterday\")\n", err: diag.AyAmis},
		{name: "ieywoo needs abiryoku", src: "\"fxma:toim\" aste toim werer\nweho f()\n  toim no ieywoo() itswa\nfinend\nf()\n", err: diag.KiNbOO},
		{name: "saimoku", src: "\"fxma:toim\" aste toim werer\nitsver abiryoku weho greet(name)\n  toim no saimoku(1) praylee\n  name redogh\nfinend\n[greet(\"a\"), greet(\"b\")] praylee\n", stdout: "a\nb\n"},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestStdDestimay(t *testing.T) {
	head := "\"fxma:destimay\" aste destimay werer\n"
	cases := []testCase{
		{name: "seed", src: head + "destimay no gomi(42)\nxs izu [destimay no senchoose(1..6), destimay no commity(), destimay no madrazy([1, 2, 3, 4])]\ndestimay no gomi(42)\nys izu [destimay no senchoose(1..6), destimay no commity(), destimay no madrazy([1, 2, 3, 4])]\nxs senaji ys redogh\n", stdout: "TadACt\n"},
		{name: "range values", src: head + "ok izu tadact\n1..200 aste i azagara\n  x izu destimay no senchoose(10...20 batte 5)\n  ok izu ok sond (x senaji 10 owk x senaji 15)\nfinend\nok redogh\n", stdout: "TadACt\n"},
		{name: "madrazy keeps elements", src: head + "\"fxma:varizama\" aste varizama werer\nxs izu [3, 1, 2]\nys izu destimay no madrazy(xs)\nvarizama no orei(ys) redogh\nxs redogh\n", stdout: "[1, 2, 3]\n[3, 1, 2]\n"},
		{name: "replayed on backtrack", src: head + "destimay no gomi(1)\na izu destimay no commity()\nc izu 1 owk 2 mekamo\nc senaji 2 shoubek\nb izu destimay no commity()\ndestimay no gomi(1)\n[a senaji destimay no commity(), b senaji destimay no commity()] redogh\n", stdout: "[TadACt, TadACt]\n"},
		{name: "uuid", src: head + "\"fxma:coword\" aste coword werer\nu izu destimay no uuid()\nnang(u) redogh\nu no 14 redogh\ndestimay no uuid(7) no 14 redogh\n", stdout: "36\n4\n7\n"},
		{name: "empty", src: head + "destimay no senchoose([])\n", err: diag.NiiExi},
		{name: "needs abiryoku", src: head + "weho f()\n  destimay no commity() itswa\nfinend\nf()\n", err: diag.KiNbOO},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}
