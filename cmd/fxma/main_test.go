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

func TestHak(t *testing.T) {
	cases := []struct {
		name, stdin    string
		code           int
		stdout, stderr string
	}{
		{
			name:   "lines and blocks",
			stdin:  "x izu 2\n\nmof x > 1\n  x * 3\nfinend\nx + 1\n",
			stdout: "> > > ... ... > 3\n> \n",
		},
		{
			name:   "give up on unfinished with blank line",
			stdin:  "mof tadact\n\ncnt\n",
			stdout: "> ... > > \n",
			stderr: "AyAmis hak:3:1 FineNd\nNii eXi' hak:3:1 cnt\n",
		},
		{
			name:   "program reads the following line",
			stdin:  "zanfession\nabc\n",
			stdout: "> \"abc\"\n> \n",
		},
		{
			name:   "exit",
			stdin:  "\"fxma:ieviljio\" aste p werer\np no diviljio() no \"shid\"(3)\n1\n",
			code:   3,
			stdout: "> > ",
		},
	}
	for _, c := range cases {
		var stdout, stderr strings.Builder
		code := run([]string{"hak"}, strings.NewReader(c.stdin), &stdout, &stderr)
		if code != c.code || stdout.String() != c.stdout || stderr.String() != c.stderr {
			t.Errorf("%s: code %d\nstdout %q\nstderr %q", c.name, code, stdout.String(), stderr.String())
		}
	}
}

func TestModCommands(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("FXMA_CACHE", t.TempDir())
	cases := []struct {
		args   []string
		code   int
		stderr string
	}{
		{args: []string{"werer", "github.com/alice/strs"}, code: 1, stderr: "Nii eXi' fxma.mod"},
		{args: []string{"ado"}, code: 2, stderr: "fiXmArge Script (fxma) の使い方"},
		{args: []string{"mod", "init", "github.com/me/app"}, code: 0},
		{args: []string{"ado", "meif", "github.com/me/app"}, code: 1, stderr: "Su·Ready fxma.mod"},
		{args: []string{"get", "github.com/alice/strs@latest"}, code: 1, stderr: "AyAmis latest"},
		{args: []string{"ado", "saymal"}, code: 0},
	}
	for _, c := range cases {
		var stdout, stderr strings.Builder
		code := run(c.args, strings.NewReader(""), &stdout, &stderr)
		if code != c.code || !strings.HasPrefix(stderr.String(), c.stderr) {
			t.Errorf("%v: code %d stdout %q stderr %q", c.args, code, stdout.String(), stderr.String())
		}
	}
}
