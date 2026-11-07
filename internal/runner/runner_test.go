package runner

import (
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/afero"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
)

type testCase struct {
	name   string
	src    string
	stdin  string
	stdout string
	err    diag.Kind                         // 空ならエラーなし
	files  map[string]string                 // werer で読み込めるファイル
	disk   map[string]string                 // fxma:memoid から見えるファイル。/work からの相対パス
	check  func(t *testing.T, fsys afero.Fs) // 実行した後のファイルを確かめる
	exit   int                               // dv no "shid" で終わるときの終了コード
	client *http.Client                      // fxma:gemon に渡す。nil なら通信できない
	proc   *interp.Process                   // fxma:ieviljio に渡す。nil なら空
}

// fixedNow は 2026-11-07T21:00:00.042+09:00。時刻を使うテストを毎回同じにする。
func fixedNow() time.Time { return time.UnixMilli(1794052800042) }

var tokyo = time.FixedZone("JST", 9*60*60)

func runCase(t *testing.T, c testCase) {
	t.Helper()
	var stdout, stderr strings.Builder
	readFile := func(name string) ([]byte, error) {
		src, ok := c.files[filepath.ToSlash(name)]
		if !ok {
			return nil, fs.ErrNotExist
		}
		return []byte(src), nil
	}
	fsys := afero.NewMemMapFs()
	if err := fsys.MkdirAll("/work", 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range c.disk {
		if err := afero.WriteFile(fsys, "/work/"+name, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	opts := interp.Options{Stdin: strings.NewReader(c.stdin), Stdout: &stdout, Stderr: &stderr, Load: Loader(readFile), WorkDir: "/work", FS: fsys, Process: c.proc, HTTPClient: c.client, Now: fixedNow, Location: tokyo}
	_, err := Run(c.src, "main.fxma", opts)
	var de *diag.Error
	exit, exited := errors.AsType[*interp.Exit](err)
	switch {
	case exited:
		if exit.Code != c.exit {
			t.Errorf("%s: 終了コード %d, want %d", c.name, exit.Code, c.exit)
		}
	case c.exit != 0:
		t.Errorf("%s: 終了していない (%v)", c.name, err)
	case c.err == "" && err != nil:
		t.Errorf("%s: エラー %v", c.name, err)
	case c.err != "" && (!errors.As(err, &de) || de.Kind != c.err):
		t.Errorf("%s: エラー %v, want %s", c.name, err, c.err.Spelling())
	}
	if c.check != nil {
		c.check(t, fsys)
	}
	if got := stdout.String(); got != c.stdout {
		t.Errorf("%s:\n got %q\nwant %q", c.name, got, c.stdout)
	}
}

func TestLanguage(t *testing.T) {
	cases := []testCase{
		{name: "compound", src: "x izu 10\nx *izu 2\nx -izu 5\nx redogh\n", stdout: "15\n"},
		{name: "division", src: "7 / 2 redogh\n7 // 2 redogh\n-7 // 2 redogh\n-7 % 3 redogh\n2 ** 10 redogh\n2 ** -1 redogh\n", stdout: "3.5\n3\n-4\n2\n1024\n0.5\n"},
		{name: "float format", src: "0.1 + 0.2 redogh\n3.0 redogh\n1e21 redogh\n1e-7 redogh\n", stdout: "0.30000000000000004\n3\n1e+21\n1e-7\n"},
		{name: "values", src: "tadact redogh\nvanita redogh\n[1, 2.5, \"a\"] redogh\n{\"k\": tadact} redogh\n", stdout: "TadACt\nVani·tA\n[1, 2.5, \"a\"]\n{\"k\": TadACt}\n"},
		{name: "equality", src: "1 senaji 1.0 redogh\n\"1\" senaji 1 redogh\ntadact senaji 1 redogh\n[1, 2] senaji [1, 2] redogh\nnii 1 senaji 2 redogh\n", stdout: "TadACt\nULieSo\nULieSo\nTadACt\nTadACt\n"},
		{name: "chain", src: "1 < 2 <= 2 redogh\n3 < 2 < 5 redogh\n", stdout: "TadACt\nULieSo\n"},
		{name: "exi", src: "2 exi [1, 2] redogh\n\"ab\" exi \"xaby\" redogh\nx izu vanita\nx exi redogh\n\"k\" exi {\"k\": 1} redogh\n", stdout: "TadACt\nTadACt\nULieSo\nTadACt\n"},
		{name: "array", src: "xs izu [3, 1, 4]\nxs no 0 izu 9\nxs no -1 redogh\nxs no 0 redogh\nnang(xs) redogh\n", stdout: "4\n9\n3\n"},
		{name: "slice", src: "xs izu [0, 1, 2, 3, 4, 5]\nxs no 2..4 redogh\nxs no 2...4 redogh\nxs no 4.. redogh\nxs no ..1 redogh\nxs no 1..2 izu [7, 8]\nxs redogh\n", stdout: "[2, 3, 4]\n[2, 3]\n[4, 5]\n[0, 1]\n[0, 7, 8, 3, 4, 5]\n"},
		{name: "float index", src: "xs izu [1, 2, 3]\nxs no 2.0 redogh\n", stdout: "3\n"},
		{name: "float index with fraction", src: "xs izu [1, 2, 3]\nxs no 1.5 redogh\n", err: diag.ImJune},
		{name: "string index", src: "s izu \"魔女裁判\"\ns no 0 redogh\ns no 1..2 redogh\nnang(s) redogh\n", stdout: "魔\n女裁\n4\n"},
		{name: "string immutable", src: "s izu \"ab\"\ns no 0 izu \"x\"\n", err: diag.ImJune},
		{name: "range", src: "0..10 batte 5 aste x azagara\n  x redogh\nfinend\n3...1 batte -1 aste x azagara\n  x redogh\nfinend\n", stdout: "0\n5\n10\n3\n2\n"},
		{name: "dict loop", src: "d izu {\"b\": 1, \"a\": 2}\nd no \"c\" izu 3\nd aste k azagara\n  \"#{k}=#{d no k}\" redogh\nfinend\n", stdout: "b=1\na=2\nc=3\n"},
		{name: "keip ketzen", src: "1..10 aste i azagara\n  mof i % 2 senaji 0\n    keip\n  finend\n  mof i > 5\n    ketzen\n  finend\n  i redogh\nfinend\n", stdout: "1\n3\n5\n"},
		{name: "interpolation", src: "x izu 2\n\"x は #{x * 2} です、#{[x, \"a\"]}\" redogh\n", stdout: "x は 4 です、[2, \"a\"]\n"},
		{name: "triple", src: "s izu \"\"\"\n    a\n      b #{1 + 1}\n    \"\"\"\ns redogh\n", stdout: "a\n  b 2\n"},
		{name: "concat repeat", src: "\"ab\" * 3 redogh\n[1] + [2] redogh\n", stdout: "ababab\n[1, 2]\n"},
		{name: "type error", src: "1 + \"a\" redogh\n", err: diag.ImJune},
		{name: "condition must be bool", src: "mof 0\nfinend\n", err: diag.ImJune},
		{name: "zero division", src: "1 // 0 redogh\n", err: diag.KYoNAihim},
		{name: "overflow", src: "9223372036854775807 + 1 redogh\n", err: diag.Bavellabion},
		{name: "overflow literal", src: "9223372036854775808 redogh\n", err: diag.Bavellabion},
		{name: "overflow pow", src: "2 ** 64 redogh\n", err: diag.Bavellabion},
		{name: "undefined", src: "cnt redogh\n", err: diag.NiiExi},
		{name: "out of range", src: "[1] no 3 redogh\n", err: diag.NiiExi},
		{name: "assign to equals is error", src: "x izu 1 == 1\n", err: diag.AyAmis},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestScope(t *testing.T) {
	cases := []testCase{
		{name: "outer assign in block", src: "count izu 0\nmof count senaji 0\n  count +izu 1\n  y izu 1\nfinend\ncount redogh\n", stdout: "1\n"},
		{name: "block var not visible", src: "mof tadact\n  y izu 1\nfinend\ny redogh\n", err: diag.NiiExi},
		{name: "function reads outer", src: "count izu 5\nabiryoku weho f()\n  count redogh\nfinend\nf()\n", stdout: "5\n"},
		{name: "function cannot assign outer", src: "count izu 0\nweho f()\n  count +izu 1\nfinend\nf()\n", err: diag.SuReady},
		{name: "param shadows outer", src: "n izu 1\nweho f(n)\n  n * 2 itswa\nfinend\nf(5) redogh\nn redogh\n", stdout: "10\n1\n"},
		{name: "hoisting", src: "f() redogh\nweho f()\n  1 itswa\nfinend\n", stdout: "1\n"},
		{name: "closure", src: "weho counter()\n  c izu [0]\n  weho ()\n    c no 0 +izu 1\n    c no 0 itswa\n  finend itswa\nfinend\nnext izu counter()\nnext()\nnext() redogh\n", stdout: "2\n"},
		{name: "anonymous one line", src: "f izu weho (n) n * 2 itswa finend\nf(21) redogh\n", stdout: "42\n"},
		{name: "showence", src: "showence limit izu 15\nlimit izu 20\n", err: diag.BlasTOac},
		{name: "showence compound", src: "showence limit izu 15\nmof tadact\n  limit +izu 1\nfinend\n", err: diag.BlasTOac},
		{name: "showence array contents", src: "showence xs izu [1]\nxs no 0 izu 2\nxs redogh\n", stdout: "[2]\n"},
		{name: "duplicate function", src: "weho f()\nfinend\nweho f()\nfinend\n", err: diag.SuReady},
		{name: "recursion", src: "weho fact(n)\n  mof n <= 1\n    1 itswa\n  finend\n  n * fact(n - 1) itswa\nfinend\nfact(20) redogh\n", stdout: "2432902008176640000\n"},
		{name: "deep recursion", src: "weho f(n)\n  f(n + 1) itswa\nfinend\nf(0)\n", err: diag.Bavellabion},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestAbiryoku(t *testing.T) {
	cases := []testCase{
		{name: "redogh without abiryoku", src: "weho greet()\n  \"hi\" redogh\nfinend\n", err: diag.KiNbOO},
		{name: "call abiryoku from plain", src: "abiryoku weho greet()\n  \"hi\" redogh\nfinend\nweho f()\n  greet()\nfinend\n", err: diag.KiNbOO},
		{name: "pass abiryoku from plain", src: "abiryoku weho greet(x)\n  x redogh\nfinend\nweho apply(f, x)\n  f(x)\nfinend\nweho g()\n  apply(greet, 1)\nfinend\n", err: diag.KiNbOO},
		{name: "pass abiryoku from abiryoku", src: "abiryoku weho greet(x)\n  x redogh\nfinend\nweho apply(f, x)\n  f(x)\nfinend\napply(greet, 1)\n", stdout: "1\n"},
		{name: "dynamic abiryoku", src: "abiryoku weho greet()\n  1 redogh\nfinend\nfs izu [greet]\nweho f()\n  fs no 0()\nfinend\nf()\n", err: diag.KiNbOO},
		{name: "top level io", src: "\"hi\" redogh\n", stdout: "hi\n"},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestBacktrack(t *testing.T) {
	cases := []testCase{
		{name: "witch", src: "witch izu 1 owk 2 owk 3 mekamo\nnii witch senaji 1 shoubek\nwitch % 2 senaji 1 shoubek\nwitch redogh\n", stdout: "3\n"},
		{name: "two choices", src: "a izu 1 owk 2 owk 3 mekamo\nb izu 1 owk 2 owk 3 mekamo\na + b senaji 5 shoubek\na < b shoubek\n\"#{a} #{b}\" redogh\n", stdout: "2 3\n"},
		{name: "output of failed path is dropped", src: "\"start\" redogh\nx izu 1 owk 2 mekamo\nx redogh\nx senaji 2 shoubek\n\"end\" redogh\n", stdout: "start\n2\nend\n"},
		{name: "variables restored", src: "xs izu []\nx izu 1 owk 2 mekamo\nxs izu xs + [x]\nx senaji 2 shoubek\nxs redogh\n", stdout: "[2]\n"},
		{name: "co-jundic", src: "x izu 1 owk 2 mekamo\nx senaji 3 shoubek\n", err: diag.CoJundic},
		{name: "no choice", src: "1 senaji 2 shoubek\n", err: diag.CoJundic},
		{name: "paren owk is or", src: "x izu (ulieso owk tadact) owk ulieso mekamo\nx redogh\n", stdout: "TadACt\n"},
		{name: "choice in function", src: "weho pick()\n  1 owk 2 owk 3 mekamo itswa\nfinend\nx izu pick()\nx > 2 shoubek\nx redogh\n", stdout: "3\n"},
		{name: "co-jundic not catchable", src: "teitance\n  1 senaji 2 shoubek\nyoulugive e\n  \"caught\" redogh\nfinend\n", err: diag.CoJundic},
		{name: "input replayed", src: "name izu \"? \" zanfession\nx izu 1 owk 2 mekamo\nx senaji 2 shoubek\n\"#{name} #{x}\" redogh\n", stdin: "abc\n", stdout: "? abc 2\n"},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestExceptions(t *testing.T) {
	cases := []testCase{
		{name: "catch", src: "teitance\n  \"bad\" accuhearts\nyoulugive e\n  e no \"lyuzun\" redogh\n  e no \"chidd\" redogh\nennote\n  \"FineNd\" redogh\nfinend\n", stdout: "bad\nAccUhEarts\nFineNd\n"},
		{name: "runtime error caught", src: "teitance\n  x izu [1] no 5\nyoulugive e\n  e no \"chidd\" redogh\nfinend\n", stdout: "Nii eXi'\n"},
		{name: "uncaught", src: "\"bad\" accuhearts\n", err: diag.KeItion},
		{name: "rethrow", src: "teitance\n  teitance\n    \"inner\" accuhearts\n  youlugive e\n    e accuhearts\n  finend\nyoulugive e\n  e no \"lyuzun\" redogh\nfinend\n", stdout: "inner\n"},
		{name: "finally on return", src: "abiryoku weho f()\n  teitance\n    1 itswa\n  ennote\n    \"fin\" redogh\n  finend\nfinend\nf() redogh\n", stdout: "fin\n1\n"},
		{name: "blastoac not catchable", src: "showence x izu 1\nteitance\n  x izu 2\nyoulugive\nfinend\n", err: diag.BlasTOac},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestAsync(t *testing.T) {
	cases := []testCase{
		{name: "await", src: "itsver weho f(x)\n  x * 2 itswa\nfinend\nf(21) praylee redogh\n[f(1), f(2)] praylee redogh\n", stdout: "42\n[2, 4]\n"},
		{name: "await error", src: "itsver weho f()\n  \"bad\" accuhearts\nfinend\np izu f()\nteitance\n  p praylee\nyoulugive e\n  e no \"lyuzun\" redogh\nfinend\n", stdout: "bad\n"},
		{name: "await outside itsver", src: "weho f()\n  1 praylee\nfinend\nf()\n", err: diag.CaNay},
		{name: "await commits choices", src: "x izu 1 owk 2 mekamo\ny izu x praylee\nx senaji 2 shoubek\n", err: diag.CoJundic},
	}
	for _, c := range cases {
		runCase(t, c)
	}
}

func TestInput(t *testing.T) {
	runCase(t, testCase{name: "zanfession", src: "name izu \"名前は? \" zanfession\n\"weRer, \" + name redogh\nrest izu zanfession\nrest exi redogh\n", stdin: "fxma\n", stdout: "名前は? weRer, fxma\nULieSo\n"})
}

func TestSentkyOnlyInTestFile(t *testing.T) {
	runCase(t, testCase{name: "sentky", src: "sentky \"x\"\n  tadact shoubek\nfinend\n", err: diag.CaNay})
}
