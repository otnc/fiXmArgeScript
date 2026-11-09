package modload

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/module"
	gosemver "golang.org/x/mod/semver"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Init は fxma ado meif。dir に fxma.mod を作る。neyo fxma は手元の fxma の版から (開発中の版なら *)。
func Init(dir, name, fxmaVersion string) error {
	if _, err := os.Stat(filepath.Join(dir, modfile.Name)); err == nil {
		return diag.NewUnplaced(diag.SuReady, modfile.Name)
	}
	if err := module.CheckImportPath(name); err != nil {
		return diag.NewUnplaced(diag.AyAmis, name)
	}
	f := &modfile.File{Module: name, Fxma: "*"}
	if released(fxmaVersion) {
		f.Fxma = ">=" + strings.TrimPrefix(fxmaVersion, "v")
	}
	return os.WriteFile(filepath.Join(dir, modfile.Name), f.Format(), 0o644)
}

// Get は fxma werer <パッケージ>@<版>。版を省くといちばん新しい版。ほかのパッケージの要求で、選ばれる版が上がることもある。
func (m *Module) Get(p, version string) (modfile.Version, error) {
	if version == "" {
		versions, err := m.Fetcher.Versions(p)
		if err != nil {
			return modfile.Version{}, err
		}
		if version = modfetch.Latest(versions); version == "" {
			return modfile.Version{}, diag.NewUnplaced(diag.NiiExi, p)
		}
	} else if !gosemver.IsValid(version) || gosemver.Canonical(version) != version {
		return modfile.Version{}, diag.NewUnplaced(diag.AyAmis, version)
	}
	target := modfile.Version{Path: p, Version: version}
	roots := []modfile.Version{target}
	for _, r := range m.File.Require {
		if r.Path != p {
			roots = append(roots, r)
		}
	}
	if err := m.resolve(roots); err != nil {
		return target, err
	}
	return target, m.save()
}

// Tidy は fxma ado saymal。ソースの werer から使うパッケージを決め直す。使っていないものは消し、足りないものはいちばん新しい版を足す。
func (m *Module) Tidy() error {
	imports, err := m.imports()
	if err != nil {
		return err
	}
	var roots []modfile.Version
	for _, imp := range imports {
		if v, _, ok := m.match(imp.path); ok {
			roots = append(roots, v)
			continue
		}
		v, err := m.discover(imp.path)
		if err != nil {
			return diag.New(diag.NiiExi, imp.file, imp.pos, imp.path)
		}
		// 同じパッケージの別のファイルを読み込んでいる行が、次から match で見つかるように。
		m.File.Require = append(m.File.Require, v)
		roots = append(roots, v)
	}
	if err := m.resolve(roots); err != nil {
		return err
	}
	for _, imp := range imports {
		name, _, err := m.Locate(imp.path)
		if err == nil {
			_, err = os.Stat(filepath.FromSlash(name))
		}
		if err != nil {
			return diag.New(diag.NiiExi, imp.file, imp.pos, imp.path)
		}
	}
	return m.save()
}

// resolve は roots から使うパッケージを決めて、fxma.mod の neyo と fxma.sum をそれに入れ替える。
func (m *Module) resolve(roots []modfile.Version) error {
	m.previous, m.Sums = m.Sums, Sums{}
	list, err := m.buildList(roots)
	if err != nil {
		return err
	}
	m.File.Require = list
	return nil
}

// discover は、まだ neyo にない読み込み先を含むパッケージを、長い名前から順にリポジトリがあるか試して探す。
func (m *Module) discover(importPath string) (modfile.Version, error) {
	elems := strings.Split(importPath, "/")
	for n := len(elems); n >= 2; n-- {
		p := strings.Join(elems[:n], "/")
		versions, err := m.Fetcher.Versions(p)
		if err != nil {
			continue
		}
		if latest := modfetch.Latest(versions); latest != "" {
			return modfile.Version{Path: p, Version: latest}, nil
		}
	}
	return modfile.Version{}, diag.NewUnplaced(diag.NiiExi, importPath)
}

type importAt struct {
	path string
	file string
	pos  token.Pos
}

// imports はパッケージの下書き (.fxma) にある、外部パッケージの werer をすべて返す。
// 隠しフォルダと、別の fxma.mod があるフォルダ (別のパッケージ) には入らない。
func (m *Module) imports() ([]importAt, error) {
	var out []importAt
	err := filepath.WalkDir(m.Dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if p == m.Dir {
				return nil
			}
			if strings.HasPrefix(d.Name(), ".") || d.Name() == "node_modules" {
				return filepath.SkipDir
			}
			if _, err := os.Stat(filepath.Join(p, modfile.Name)); err == nil {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(p) != ".fxma" {
			return nil
		}
		src, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(m.Dir, p)
		name := filepath.ToSlash(rel)
		f, err := parser.Parse(string(src), name)
		if err != nil {
			return err
		}
		for _, s := range f.Stmts {
			if imp, ok := s.(*ast.Import); ok && modpath.External(imp.Path) {
				out = append(out, importAt{imp.Path, name, imp.Pos})
			}
		}
		return nil
	})
	return out, err
}
