package runner

import (
	"strings"
	"testing"

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
