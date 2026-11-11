package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/otnc/fiXmArgeScript/internal/lint"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
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
	findings := lint.File(f, string(src), modfile.AllowedRules(filepath.Dir(path)))
	for _, fd := range findings {
		fmt.Fprintln(s.stdout, fd)
	}
	return len(findings) == 0
}
