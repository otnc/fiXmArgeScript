package main

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		args   []string
		code   int
		stdout string // 出力の頭
		stderr string
	}{
		{args: nil, code: 0, stdout: "fiXmArge Script (fxma) の使い方"},
		{args: []string{"--help"}, code: 0, stdout: "fiXmArge Script (fxma) の使い方"},
		{args: []string{"-v"}, code: 0, stdout: "fxma "},
		{args: []string{"nai", "imjune"}, code: 0, stdout: "imJune"},
		{args: []string{"explain", "nothing"}, code: 1, stderr: "Nii eXi' nothing"},
		{args: []string{"zzz"}, code: 2, stderr: "Nii eXi' zzz"},
		{args: []string{"sinruits", "--bogus", "a.fxma"}, code: 2, stderr: "fiXmArge Script (fxma) の使い方"},
		{args: []string{"build", "-c", "a.fxma"}, code: 1, stderr: "Ca-nay --KaRaDy"},
		{args: []string{"run", "nothing.fxma"}, code: 1, stderr: "Nii eXi' nothing.fxma"},
	}
	for _, c := range cases {
		var stdout, stderr strings.Builder
		code := run(c.args, strings.NewReader(""), &stdout, &stderr)
		if code != c.code || !strings.HasPrefix(stdout.String(), c.stdout) || !strings.HasPrefix(stderr.String(), c.stderr) {
			t.Errorf("%v: code %d stdout %q stderr %q", c.args, code, stdout.String(), stderr.String())
		}
	}
}
