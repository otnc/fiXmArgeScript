package modload

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
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

// buildList は roots から始めて、使うパッケージの版を最小バージョン選択 (Go と同じ) で決める。
// たどった版の要求をすべて集め、パッケージごとにいちばん新しい版を1つ選ぶ。取ってきたものは fixmarge.sum に足す。
func (m *Module) buildList(roots []modfile.Version) ([]modfile.Version, error) {
	selected := map[string]string{}
	visited := map[modfile.Version]bool{}
	queue := roots
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		v, err := m.notExcluded(v)
		if err != nil {
			return nil, err
		}
		if visited[v] {
			continue
		}
		visited[v] = true
		if cur, ok := selected[v.Path]; !ok || gosemver.Compare(v.Version, cur) > 0 {
			selected[v.Path] = v.Version
		}
		reqs, err := m.requirements(v)
		if err != nil {
			return nil, err
		}
		queue = append(queue, reqs...)
	}
	list := make([]modfile.Version, 0, len(selected))
	for p, ver := range selected {
		list = append(list, modfile.Version{Path: p, Version: ver})
	}
	return list, nil
}

// notExcluded は、unneyo で使わせない版なら、その次に新しい使える版に上げる (Go の exclude と同じ)。
func (m *Module) notExcluded(v modfile.Version) (modfile.Version, error) {
	if !m.File.Excluded(v) {
		return v, nil
	}
	versions, err := m.Fetcher.Versions(v.Path)
	if err != nil {
		return v, err
	}
	for _, next := range versions {
		cand := modfile.Version{Path: v.Path, Version: next}
		if gosemver.Compare(next, v.Version) > 0 && !m.File.Excluded(cand) {
			return cand, nil
		}
	}
	return v, diag.NewUnplaced(diag.NiiExi, v.String())
}

// requirements は版 v の fixmarge.mod を読み、そのパッケージが使うものを返す。
// 取り込めるのは、foryulvu と neyo fxma があって、手元の fxma が範囲に入るパッケージだけ (RFC の fixmarge.mod の節)。
func (m *Module) requirements(v modfile.Version) ([]modfile.Version, error) {
	dir, err := m.fetch(v)
	if err != nil {
		return nil, err
	}
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
	return f.Require, nil
}

// fetch はパッケージを手元に置き、fixmarge.sum と照らし合わせる。初めて取ったものは fixmarge.sum に足す。
func (m *Module) fetch(v modfile.Version) (string, error) {
	if m.File.ReplaceDir(v.Path) != "" {
		return m.dir(v)
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
