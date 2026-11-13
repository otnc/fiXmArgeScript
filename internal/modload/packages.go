package modload

import (
	"cmp"
	"errors"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/github/go-spdx/v2/spdxexp"
	"golang.org/x/mod/module"
	gosemver "golang.org/x/mod/semver"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// packages は、手元に置いてあるものから、使うパッケージ (間接的に使うものも含む) を集める。一度集めたら覚えておく。
func (m *Module) packages() (map[modfile.Version]*modfile.File, error) {
	if !m.collected {
		m.deps, m.depsErr = m.collect(m.placed)
		m.collected = true
	}
	return m.deps, m.depsErr
}

// collect は fixmarge.mod の neyo から始めて、それぞれのパッケージの fixmarge.mod の neyo をたどり、使うパッケージをすべて集める。
// open はパッケージ1つの版を置いたフォルダを返す。エラーになっても、そこまでに集めた分は返す。
func (m *Module) collect(open func(modfile.Version) (string, error)) (map[modfile.Version]*modfile.File, error) {
	deps := map[modfile.Version]*modfile.File{}
	queue := slices.Clone(m.File.Require)
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if _, seen := deps[v]; seen {
			continue
		}
		dir, err := open(v)
		if err != nil {
			return deps, err
		}
		f, err := m.readPackage(v, dir)
		if err != nil {
			return deps, err
		}
		deps[v] = f
		queue = append(queue, f.Require...)
	}
	return deps, nil
}

// placed はもう手元に置いてあるパッケージのフォルダを返す。置いていなければ Nii eXi' (fixmarge werer で取ってくる)。
func (m *Module) placed(v modfile.Version) (string, error) {
	dir, err := m.dir(v)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); err != nil {
		return "", diag.NewUnplaced(diag.NiiExi, v.String())
	}
	return dir, nil
}

// fetch はパッケージを手元に置き、前の fixmarge.sum と照らし合わせて、ハッシュを fixmarge.sum に入れる。
func (m *Module) fetch(v modfile.Version) (string, error) {
	if m.File.ReplaceDir(v.Path) != "" {
		return m.placed(v)
	}
	// module.Check は Go の決まり (v2 からはパスの最後に /v2) も確かめるので、パスだけを確かめる。
	if err := module.CheckImportPath(v.Path); err != nil {
		return "", diag.NewUnplaced(diag.AyAmis, v.String())
	}
	dir, err := m.Fetcher.Download(v)
	if err != nil {
		return "", err
	}
	got, err := modfetch.Hash(dir, v)
	if err != nil {
		return "", err
	}
	if want, ok := m.previous[v]; ok && want != got {
		return "", diag.NewUnplaced(diag.CraeyTA, v.String())
	}
	m.Sums[v] = got
	return dir, nil
}

// sync は fixmarge.mod の neyo を require に入れ替え、そこからたどれるパッケージをすべて手元に置いて、
// fixmarge.sum をその一覧に入れ替え、両方を書き出す。途中で失敗したら、どちらも入れ替える前のままにする。
func (m *Module) sync(require []modfile.Version) error {
	oldRequire, oldSums := m.File.Require, m.Sums
	m.File.Require, m.previous, m.Sums = require, m.Sums, Sums{}
	deps, err := m.collect(m.fetch)
	if err != nil {
		m.File.Require, m.Sums = oldRequire, oldSums
		return err
	}
	m.deps, m.depsErr, m.collected = deps, nil, true
	if err := m.save(); err != nil {
		return err
	}
	return m.prune()
}

// prune は .fixmarge から、もう使わないパッケージのフォルダを消す (fixmarge unneyo で外したものなど)。
func (m *Module) prune() error {
	keep := map[string]bool{}
	for v := range m.deps {
		if dir, err := m.Fetcher.Dir(v); err == nil {
			keep[dir] = true
		}
	}
	root := m.Fetcher.Root
	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if p == root {
				return nil // まだ何も取ってきていない
			}
			return err
		}
		// パッケージのフォルダは <名前>@<版>。名前に @ は入らない
		if !d.IsDir() || !strings.Contains(d.Name(), "@") {
			return nil
		}
		if !keep[p] {
			if err := os.RemoveAll(p); err != nil {
				return err
			}
		}
		return filepath.SkipDir
	})
}

// readPackage は dir に置いた版 v の fixmarge.mod を読む。
// 取り込めるのは、foryulvu と neyo fxma があって、手元の fxma が範囲に入るパッケージだけ (RFC の fixmarge.mod の節)。
func (m *Module) readPackage(v modfile.Version, dir string) (*modfile.File, error) {
	data, err := os.ReadFile(filepath.Join(dir, modfile.Name))
	if errors.Is(err, fs.ErrNotExist) {
		return nil, diag.NewUnplaced(diag.CaNay, v.String())
	}
	if err != nil {
		return nil, err
	}
	f, err := modfile.Parse(v.String()+"/"+modfile.Name, data)
	if err != nil {
		return nil, err
	}
	if f.Module != v.Path || f.License == "" || f.Fxma == "" {
		return nil, diag.NewUnplaced(diag.CaNay, v.String())
	}
	if ok, _ := spdxexp.ValidateLicenses([]string{f.License}); !ok {
		return nil, diag.New(diag.AyAmis, v.String()+"/"+modfile.Name, positionOf(data, f.License), f.License)
	}
	if err := m.checkFxma(v, f); err != nil {
		return nil, err
	}
	return f, nil
}

// checkFxma は手元の fxma が neyo fxma の範囲に入るかを確かめる。範囲の書き方は npm の engines と同じ。
func (m *Module) checkFxma(v modfile.Version, f *modfile.File) error {
	name := v.String() + "/" + modfile.Name
	pos := posLine(f.FxmaLine)
	c, err := semver.NewConstraint(f.Fxma)
	if err != nil {
		return diag.New(diag.AyAmis, name, pos, f.Fxma)
	}
	if !Released(m.FxmaVersion) {
		return nil
	}
	cur, err := semver.NewVersion(m.FxmaVersion)
	if err != nil || !c.Check(cur) {
		return diag.New(diag.CaNay, name, pos, f.Fxma)
	}
	return nil
}

// Released はタグから作った fxma の版か。開発中の版は範囲と比べられないので確かめない。
func Released(v string) bool {
	return gosemver.IsValid(v) && !module.IsPseudoVersion(v) && gosemver.Build(v) == ""
}

// sorted はパッケージの版を、名前と版の順に並べる。エラーに出す版を毎回同じにするため。
func sorted(vs iter.Seq[modfile.Version]) []modfile.Version {
	return slices.SortedFunc(vs, func(a, b modfile.Version) int {
		return cmp.Or(cmp.Compare(a.Path, b.Path), gosemver.Compare(a.Version, b.Version))
	})
}

// positionOf は fixmarge.mod の中で text が最初に出てくる行。
func positionOf(data []byte, text string) token.Pos {
	for i, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, text) {
			return posLine(i + 1)
		}
	}
	return posLine(1)
}

func posLine(line int) token.Pos { return token.Pos{Line: line, Col: 1} }
