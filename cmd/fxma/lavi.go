package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/lint"
	"github.com/otnc/fiXmArgeScript/internal/runner"
)

// lavi は法の悪魔に規則違反を指摘させる (ESLint や go vet と同じ)。指摘が1つでもあれば終了コード1。
func (s streams) lavi(_ context.Context, cmd *cli.Command) error {
	targets := cmd.Args().Slice()
	if len(targets) == 0 {
		targets = []string{"."}
	}
	failed := false
	for _, t := range targets {
		files, err := findFiles(t, ".fxma")
		if err != nil {
			fmt.Fprintln(s.stderr, readError(err))
			failed = true
			continue
		}
		for _, path := range files {
			if !s.lintFile(path) {
				failed = true
			}
		}
	}
	if failed {
		return exitCode(1)
	}
	return nil
}

// lintFile はファイル1つを調べ、指摘がなければ true を返す。
func (s streams) lintFile(path string) bool {
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(s.stderr, readError(err))
		return false
	}
	f, err := runner.Prepare(string(src), path)
	if err != nil {
		fmt.Fprintln(s.stderr, err)
		return false
	}
	findings := lint.File(f, string(src), allowedRules(filepath.Dir(path)))
	for _, fd := range findings {
		fmt.Fprintln(s.stdout, fd)
	}
	return len(findings) == 0
}

// allowedRules は、dir から上にたどって最初に見つけた fxma.mod の youlugive に並べた規則を返す。
func allowedRules(dir string) []string {
	for {
		data, err := os.ReadFile(filepath.Join(dir, "fxma.mod"))
		if err == nil {
			var rules []string
			for line := range strings.SplitSeq(string(data), "\n") {
				if fields := strings.Fields(line); len(fields) > 0 && fields[0] == "youlugive" {
					rules = append(rules, fields[1:]...)
				}
			}
			return rules
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil
		}
		dir = parent
	}
}
