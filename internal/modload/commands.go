package modload

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/mod/module"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Init は fixmarge ado meif。dir に fixmarge.mod を作る。neyo fxma は手元の fxma の版から (開発中の版なら *)。
func Init(dir, name, fxmaVersion string) error {
	if _, err := os.Stat(filepath.Join(dir, modfile.Name)); err == nil {
		return diag.NewUnplaced(diag.SuReady, modfile.Name)
	}
	if err := module.CheckImportPath(name); err != nil {
		return diag.NewUnplaced(diag.AyAmis, name)
	}
	f := &modfile.File{Module: name, Fxma: "*"}
	if Released(fxmaVersion) {
		f.Fxma = ">=" + strings.TrimPrefix(fxmaVersion, "v")
	}
	return os.WriteFile(filepath.Join(dir, modfile.Name), f.Format(), 0o644)
}

// Install は引数なしの fixmarge werer。fixmarge.mod に書いたパッケージと、それが使うパッケージをすべて手元に置く (npm install)。
func (m *Module) Install() error { return m.sync(m.File.Require) }

// Get は fixmarge werer <パッケージ>[@<版>]。版の指定に合ういちばん新しい版を取ってきて、fixmarge.mod に足す。
// 同じパッケージの別の版がもう入っていても、そのまま残す。
func (m *Module) Get(p, spec string) (modfile.Version, error) {
	target, err := m.latest(p, spec)
	if err != nil {
		return target, err
	}
	require := slices.Clone(m.File.Require)
	if !slices.Contains(require, target) {
		require = append(require, target)
	}
	return target, m.sync(require)
}

// Remove は fixmarge unneyo <パッケージ>[@<版>]。入っている版のうち、版の指定に合ういちばん新しいものを fixmarge.mod から外す。
// 版が * なら、そのパッケージの版をすべて外す。合うものが入っていなければ Nii eXi'。
func (m *Module) Remove(p, spec string) error {
	gone := func(r modfile.Version) bool { return r.Path == p }
	if spec != modfile.AllVersions {
		if !modfile.ValidSpec(spec) {
			return diag.NewUnplaced(diag.AyAmis, spec)
		}
		target := modfile.Version{Path: p, Version: modfile.Pick(m.File.Versions(p), spec)}
		gone = func(r modfile.Version) bool { return r == target }
	}
	if !slices.ContainsFunc(m.File.Require, gone) {
		return diag.NewUnplaced(diag.NiiExi, named(p, spec))
	}
	return m.sync(slices.DeleteFunc(slices.Clone(m.File.Require), gone))
}

// Tidy は fixmarge ado saymal。ソースの werer から使うパッケージを決め直す。
// どの werer も指していない版は外し、合う版が入っていない werer には、合ういちばん新しい版を取ってきて足す。
func (m *Module) Tidy() error {
	imports, err := m.imports()
	if err != nil {
		return err
	}
	used := map[modfile.Version]bool{}
	available := &modfile.File{Require: slices.Clone(m.File.Require)}
	for _, imp := range imports {
		target, spec := modfile.SplitSpec(imp.path)
		if !modfile.ValidSpec(spec) {
			return diag.New(diag.AyAmis, imp.file, imp.pos, imp.path)
		}
		v, _, ok := pick(available, target, spec)
		if !ok {
			if v, err = m.discover(target, spec); err != nil {
				return diag.New(diag.NiiExi, imp.file, imp.pos, imp.path)
			}
			// 同じ版を指すほかの werer が、次から pick で見つかるように
			available.Require = append(available.Require, v)
		}
		used[v] = true
	}
	if err := m.sync(slices.DeleteFunc(available.Require, func(r modfile.Version) bool { return !used[r] })); err != nil {
		return err
	}
	for _, imp := range imports {
		name, _, err := m.Locate(imp.file, imp.path)
		if err == nil {
			_, err = os.Stat(filepath.FromSlash(name))
		}
		if err != nil {
			return diag.New(diag.NiiExi, imp.file, imp.pos, imp.path)
		}
	}
	return nil
}

// latest は、リポジトリにあるパッケージ p の版のうち、spec に合ういちばん新しいものを返す。
func (m *Module) latest(p, spec string) (modfile.Version, error) {
	v := modfile.Version{Path: p}
	if !modfile.ValidSpec(spec) {
		return v, diag.NewUnplaced(diag.AyAmis, spec)
	}
	versions, err := m.Fetcher.Versions(p)
	if err != nil {
		return v, err
	}
	if v.Version = modfile.Pick(versions, spec); v.Version == "" {
		return v, diag.NewUnplaced(diag.NiiExi, named(p, spec))
	}
	return v, nil
}

// named は、エラーに出すパッケージと版の指定。
func named(p, spec string) string {
	if spec == "" {
		return p
	}
	return p + "@" + spec
}

// discover は、まだ neyo にない読み込み先を含むパッケージを、長い名前から順にリポジトリがあるか試して探す。
func (m *Module) discover(target, spec string) (modfile.Version, error) {
	elems := strings.Split(target, "/")
	for n := len(elems); n >= 2; n-- {
		if v, err := m.latest(strings.Join(elems[:n], "/"), spec); err == nil {
			return v, nil
		}
	}
	return modfile.Version{}, diag.NewUnplaced(diag.NiiExi, target)
}

type importAt struct {
	path string
	file string
	pos  token.Pos
}

// imports はパッケージの下書き (.fxma) にある、外部パッケージの werer をすべて返す。
// 隠しフォルダと、別の fixmarge.mod があるフォルダ (別のパッケージ) には入らない。
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
		for _, imp := range ast.Imports(f.Stmts) {
			if modpath.External(imp.Path) {
				out = append(out, importAt{imp.Path, name, imp.Pos})
			}
		}
		return nil
	})
	return out, err
}
