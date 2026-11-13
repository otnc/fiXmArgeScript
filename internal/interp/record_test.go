package interp

import (
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/parser"
)

// mekamo がどこにもなければやり直しは起きないので、入力や時刻を何度読んでも記録は増えない。
func TestRecordsOnlyWhenReplayable(t *testing.T) {
	const reads = 100
	loop := "\":toim\" aste toim werer\n1..100 aste i azagara\n  toim no ieywoo()\n  zanfession\nfinend\n"
	pick := "weho pick()\n  1 owk 2 mekamo itswa\nfinend\n"
	plain := "weho double(x)\n  x * 2 itswa\nfinend\n"
	cases := []struct {
		name       string
		src        string
		files      map[string]string
		replayable bool
	}{
		{name: "no mekamo", src: loop},
		{name: "mekamo in file", src: loop + "x izu 1 owk 2 mekamo\n", replayable: true},
		{name: "import without mekamo", src: "\"./lib\" aste lib werer\n" + loop, files: map[string]string{"lib.fxma": plain}},
		{name: "mekamo in import", src: "\"./lib\" aste lib werer\n" + loop, files: map[string]string{"lib.fxma": pick}, replayable: true},
		{name: "mekamo in nested import", src: "\"./a\" aste a werer\n" + loop, files: map[string]string{"a.fxma": "\"./b\" aste b werer\n", "b.fxma": pick}, replayable: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			in, f := prepareRecordTest(t, c.src, c.files, strings.Repeat("x\n", reads))
			if in.replayable != c.replayable {
				t.Fatalf("replayable = %v, want %v", in.replayable, c.replayable)
			}
			err := in.runOnce(f)
			in.finish()
			if err != nil {
				t.Fatal(err)
			}
			want := 0
			if c.replayable {
				want = reads
			}
			if len(in.records) != want || len(in.inputs) != want {
				t.Errorf("records = %d, inputs = %d, want %d", len(in.records), len(in.inputs), want)
			}
		})
	}
}

// 読み込めないファイルは中身がわからないので、mekamo があるものとして記録する。
func TestMissingImportIsReplayable(t *testing.T) {
	in, _ := prepareRecordTest(t, "\"./lib\" aste lib werer\n", nil, "")
	if !in.replayable {
		t.Error("replayable = false, want true")
	}
}

func prepareRecordTest(t *testing.T, src string, files map[string]string, stdin string) (*Interp, *parser.File) {
	t.Helper()
	f, err := parser.Parse(src, "main.fxma")
	if err != nil {
		t.Fatal(err)
	}
	load := func(name string) (*parser.File, error) {
		src, ok := files[name]
		if !ok {
			return nil, fs.ErrNotExist
		}
		return parser.Parse(src, name)
	}
	in := newInterp(f.Name, Options{Stdin: strings.NewReader(stdin), Stdout: io.Discard, Stderr: io.Discard, Load: load})
	in.replayable = in.mayChoose(f, map[string]bool{f.Name: true})
	return in, f
}
