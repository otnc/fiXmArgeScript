// Package modfile は fxma.mod を読む。今は youlugive の行 (fxma lavi で許す規則) だけを読む。
// CLI と LSP で同じ規則を使うためにここに置く。
package modfile

import (
	"os"
	"path/filepath"
	"strings"
)

// AllowedRules は、dir から上にたどって最初に見つけた fxma.mod の youlugive に並べた規則を返す。
func AllowedRules(dir string) []string {
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
