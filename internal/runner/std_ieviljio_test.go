package runner

import (
	"errors"
	"io/fs"
	"strings"
	"sync"
	"testing"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
)

// fakeChild は、入力をそのまま出力に返す魔女。sleep だけは殺されるまで終わらない。
// 本物の子プロセスと同じく、Wait と Kill は別々の goroutine から同時に呼ばれる。
type fakeChild struct {
	c      interp.Command
	killed chan struct{}
	once   sync.Once
}

func (f *fakeChild) Pid() int { return 42 }

func (f *fakeChild) Wait() (interp.ChildResult, error) {
	if f.c.Name == "sleep" {
		<-f.killed
		return interp.ChildResult{Code: -1}, nil
	}
	if f.c.OnOutput != nil {
		f.c.OnOutput(false, "out")
		f.c.OnOutput(true, "warn")
	}
	out := f.c.Name + " " + strings.Join(f.c.Args, " ") + "|" + f.c.Stdin + "|" + f.c.Dir + "|" + strings.Join(f.c.Env, ",")
	return interp.ChildResult{Stdout: out, Stderr: "err", Code: len(f.c.Args)}, nil
}

func (f *fakeChild) Kill() error {
	f.once.Do(func() { close(f.killed) })
	return nil
}

// fakeProcess は偽の大魔女。spawned に生んだ魔女の数を数える。
func fakeProcess(spawned *int) *interp.Process {
	return &interp.Process{
		Args:     []string{"a", "b"},
		Env:      []string{"HOME=/home/x", "EMPTY="},
		Pid:      7,
		Platform: "linux",
		Version:  "v0.1.0",
		Spawn: func(c interp.Command) (interp.Child, error) {
			if c.Name == "missing" {
				return nil, fs.ErrNotExist
			}
			if c.Name == "broken" {
				return nil, errors.New("broken")
			}
			*spawned++
			return &fakeChild{c: c, killed: make(chan struct{})}, nil
		},
	}
}

func TestStdIeviljio(t *testing.T) {
	head := "\"fxma:ieviljio\" aste ieviljio werer\ndv izu ieviljio no diviljio()\n"
	var spawned int
	cases := []testCase{
		{name: "info", src: head + "dv no \"conwa\" redogh\ndv no \"lowhell\" no \"HOME\" redogh\ndv no \"ieviljio\" redogh\ndv no \"wid\" redogh\ndv no \"neyar\" redogh\n", stdout: "[\"a\", \"b\"]\n/home/x\n7\nlinux\nv0.1.0\n"},
		{name: "spawn and wait", src: head + "w izu dv no \"meif\"(\"sort\", [\"-r\"], {\"zanfession\": \"b\", \"edo\": \"data\", \"lowhell\": {\"K\": 1}})\nr izu w praylee\nr no \"redogh\" redogh\nr no \"canosiri\" redogh\nr no \"shid\" redogh\nnang(dv no \"alte\") redogh\n", stdout: "sort -r|b|/work/data|K=1\nerr\n1\n1\n"},
		{name: "await many", src: head + "rs izu [dv no \"meif\"(\"a\", []), dv no \"meif\"(\"b\", [\"x\"])] praylee\nrs no 1 no \"shid\" redogh\n", stdout: "1\n"},
		{name: "kill", src: head + "w izu dv no \"meif\"(\"sleep\", [])\ndv no \"murm\"(w)\n(w praylee) no \"shid\" redogh\n", stdout: "-1\n"},
		{name: "kill passed as capability", src: head + "abiryoku weho cleanup(ws, murm)\n  ws aste w azagara\n    murm(w)\n  finend\nfinend\nw izu dv no \"meif\"(\"sleep\", [])\ncleanup(dv no \"alte\", dv no \"murm\")\n(w praylee) no \"shid\" redogh\n", stdout: "-1\n"},
		{name: "no murm in package", src: head + "ieviljio no murm\n", err: diag.NiiExi},
		{name: "missing command", src: head + "dv no \"meif\"(\"missing\", [])\n", err: diag.NiiExi},
		{name: "broken command", src: head + "dv no \"meif\"(\"broken\", [])\n", err: diag.CaNay},
		{name: "unknown option", src: head + "dv no \"meif\"(\"a\", [], {\"zzz\": 1})\n", err: diag.NiiExi},
		{name: "exit skips ennote", src: head + "teitance\n  \"before\" redogh\n  dv no \"shid\"(3)\nyoulugive e\n  \"caught\" redogh\nennote\n  \"ennote\" redogh\nfinend\n", stdout: "before\n", exit: 3},
		{name: "needs abiryoku", src: head + "weho f()\n  ieviljio no diviljio()\nfinend\nf()\n", err: diag.KiNbOO},
		{name: "no spawn in browser", src: "\"fxma:ieviljio\" aste ieviljio werer\ndv izu ieviljio no diviljio()\ndv no \"meif\"(\"a\", [])\n", err: diag.CaNay},
	}
	for _, c := range cases {
		if c.name != "no spawn in browser" {
			c.proc = fakeProcess(&spawned)
		}
		runCase(t, c)
	}
}

// 魔女を生んだ後の選択肢で巻き戻っても、やり直しで同じ魔女を生み直さない。
func TestWitchNotRespawnedOnBacktrack(t *testing.T) {
	var spawned int
	runCase(t, testCase{
		name:   "respawn",
		src:    "\"fxma:ieviljio\" aste ieviljio werer\ndv izu ieviljio no diviljio()\nr izu dv no \"meif\"(\"a\", []) praylee\nx izu 1 owk 2 mekamo\nx senaji 2 shoubek\nr no \"shid\" redogh\n",
		stdout: "0\n",
		proc:   fakeProcess(&spawned),
	})
	if spawned != 1 {
		t.Errorf("魔女を %d 回生んだ", spawned)
	}
}
