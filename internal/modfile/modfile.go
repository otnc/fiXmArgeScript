// Package modfile は fxma.ado を読む。今は youlugive の行 (fxma lavi で許す規則) だけを読む。
// CLI と LSP で同じ規則を使うためにここに置く。
package modfile

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/lexicon"
)

// Name はパッケージのファイル (npm の package.json、Go の go.mod) の名前。
// Sum はロックファイル (package-lock.json、go.sum) の名前。
// 拡張子を独自の語にするのは、エディタのアイコンを Go の go.mod などと取り違えさせないため。
var (
	Name = fileName("ado")
	Sum  = fileName("showroof")
)

func fileName(skel string) string {
	e, _, ok := lexicon.Lookup(skel)
	if !ok || e.Uses&lexicon.FileName == 0 {
		panic("modfile: ファイル名の語が表にない: " + skel)
	}
	return "fxma." + e.Name
}

// AllowedRules は、dir から上にたどって最初に見つけた fxma.ado の youlugive に並べた規則を返す。
func AllowedRules(dir string) []string {
	for {
		data, err := os.ReadFile(filepath.Join(dir, Name))
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
