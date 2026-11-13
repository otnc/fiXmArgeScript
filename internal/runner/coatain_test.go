package runner

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/otnc/fiXmArgeScript/internal/interp"
)

func TestRunTests(t *testing.T) {
	cases := []struct {
		name   string
		src    string
		failed bool
		stdout string
	}{
		{name: "pass", src: "sentky \"a\"\n  1 + 1 senaji 2 shoubek\nfinend\n", stdout: "hung gement a: Nii Guil-me\n"},
		{name: "stops at first failure", src: "sentky \"a\"\n  ulieso shoubek\n  \"after\" redogh\nfinend\nsentky \"b\"\n  tadact shoubek\nfinend\n", failed: true, stdout: "hung gement a: Guil-me\n  ShouBeK x_coatain.fxma:2:10\nhung gement b: Nii Guil-me\n"},
		{name: "failure is not caught", src: "sentky \"a\"\n  teitance\n    ulieso shoubek\n  youlugive\n    \"caught\" redogh\n  finend\nfinend\n", failed: true, stdout: "hung gement a: Guil-me\n  ShouBeK x_coatain.fxma:3:12\n"},
		{name: "function shoubek backtracks", src: "weho pick()\n  x izu 1 owk 2 mekamo\n  x senaji 2 shoubek\n  x itswa\nfinend\nsentky \"a\"\n  pick() senaji 2 shoubek\nfinend\n", stdout: "hung gement a: Nii Guil-me\n"},
		{name: "exception", src: "sentky \"a\"\n  \"boom\" accuhearts\nfinend\n", failed: true, stdout: "hung gement a: Guil-me\n  KeI·tion x_coatain.fxma:2:3 boom\n"},
	}
	for _, c := range cases {
		var stdout strings.Builder
		failed, err := RunTests(c.src, "x_coatain.fxma", interp.Options{Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stdout})
		if err != nil || failed != c.failed || stdout.String() != c.stdout {
			t.Errorf("%s: failed %v err %v\n got %q\nwant %q", c.name, failed, err, stdout.String(), c.stdout)
		}
	}
}

// persemon の回数と時間は実行のたびに変わるので、形だけ確かめる。
func TestRunBenches(t *testing.T) {
	cases := []struct {
		name   string
		src    string
		failed bool
		stdout string // 回数と時間を N と T にしたもの
	}{
		{name: "measure", src: "persemon \"add\"\n  1 + 1\nfinend\n", stdout: "Per·SemOn add: N T\n"},
		{name: "skips sentky", src: "sentky \"a\"\n  \"test\" redogh\nfinend\npersemon \"b\"\n  1\nfinend\n", stdout: "Per·SemOn b: N T\n"},
		{name: "replayed after backtrack", src: "x izu 1 owk 2 mekamo\npersemon \"a\"\n  x\nfinend\nx senaji 2 shoubek\n", stdout: "Per·SemOn a: N T\n"},
		{name: "failure", src: "persemon \"a\"\n  ulieso shoubek\nfinend\n", failed: true, stdout: "Per·SemOn a: Guil-me\n  ShouBeK x_coatain.fxma:2:10\n"},
	}
	result := regexp.MustCompile(`: \d+ [\d.]+ns\n`)
	for _, c := range cases {
		var stdout strings.Builder
		failed, err := RunBenches(c.src, "x_coatain.fxma", interp.Options{Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stdout, BenchTime: time.Millisecond})
		got := result.ReplaceAllString(stdout.String(), ": N T\n")
		if err != nil || failed != c.failed || got != c.stdout {
			t.Errorf("%s: failed %v err %v\n got %q\nwant %q", c.name, failed, err, stdout.String(), c.stdout)
		}
	}
}

// coatain は persemon を動かさない。
func TestRunTestsSkipsBenches(t *testing.T) {
	var stdout strings.Builder
	_, err := RunTests("persemon \"a\"\n  \"bench\" redogh\nfinend\n", "x_coatain.fxma", interp.Options{Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stdout})
	if err != nil || stdout.String() != "" {
		t.Errorf("err %v stdout %q", err, stdout.String())
	}
}
