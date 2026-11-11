// Package modfile は fixmarge.mod を読み書きする。CLI・LSP・パッケージの取得で同じ規則を使うためにここに置く。
package modfile

import (
	"cmp"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Name はパッケージのファイル (npm の package.json、Go の go.mod) の名前。Sum はロックファイル (go.sum) の名前。
const (
	Name = "fixmarge.mod"
	Sum  = "fixmarge.sum"
)

// 項目名。語の表 (ModField) にあるものだけを使う。
var (
	fieldModule  = field("ado")
	fieldPurpose = field("mokpose")
	fieldRequire = field("neyo")
	fieldExclude = field("unneyo")
	fieldLicense = field("foryulvu")
	fieldReplace = field("lainei")
	fieldAllow   = field("youlugive")
)

// Fxma は neyo fxma の fxma (パッケージ名ではなく、対応する fxma の版の行の印)。
const Fxma = "fxma"

func field(skel string) string {
	e, _, ok := lexicon.Lookup(skel)
	if !ok || e.Uses&lexicon.ModField == 0 {
		panic("modfile: 項目名が表にない: " + skel)
	}
	return e.Name
}

// File は fixmarge.mod の中身。
type File struct {
	Module   string    // ado
	Purpose  string    // mokpose
	Fxma     string    // neyo fxma の範囲 (npm の engines と同じ書き方)
	License  string    // foryulvu (SPDX のライセンス式)
	Require  []Version // neyo。間接的に使うものも含めて、使うパッケージをすべて並べる
	Exclude  []Version // unneyo
	Replace  []Replace // lainei
	Allowed  []string  // youlugive (fixmarge lavi で指摘しない規則)
	FxmaLine int       // neyo fxma の行。範囲の外のときのエラーの位置に使う
}

// Version はパッケージ1つとその版。
type Version struct {
	Path    string
	Version string
}

func (v Version) String() string { return v.Path + "@" + v.Version }

// Replace は lainei。パッケージを手元のフォルダ (fixmarge.mod からの相対パス) のもので置き換える。
type Replace struct {
	Path string
	Dir  string
}

// Parse は fixmarge.mod を読む。name はエラーに出すファイル名。
func Parse(name string, data []byte) (*File, error) {
	f := &File{}
	seen := map[string]bool{}
	for i, line := range strings.Split(string(data), "\n") {
		pos := token.Pos{Line: i + 1, Col: 1}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		key, args := fields[0], fields[1:]
		bad := func() error { return diag.New(diag.AyAmis, name, pos, key) }
		single := func() error {
			if seen[key] {
				return diag.New(diag.SuReady, name, pos, key)
			}
			seen[key] = true
			return nil
		}
		switch {
		case key == fieldModule && len(args) == 1:
			if err := single(); err != nil {
				return nil, err
			}
			f.Module = args[0]
		case key == fieldPurpose && len(args) > 0:
			if err := single(); err != nil {
				return nil, err
			}
			s, err := strconv.Unquote(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), key)))
			if err != nil {
				return nil, bad()
			}
			f.Purpose = s
		case key == fieldRequire && len(args) > 1 && args[0] == Fxma:
			if seen[Fxma] {
				return nil, diag.New(diag.SuReady, name, pos, Fxma)
			}
			seen[Fxma] = true
			f.Fxma = strings.Join(args[1:], " ")
			f.FxmaLine = pos.Line
		case key == fieldLicense && len(args) > 0:
			if err := single(); err != nil {
				return nil, err
			}
			f.License = strings.Join(args, " ")
		case key == fieldRequire && len(args) == 2, key == fieldExclude && len(args) == 2:
			v := Version{args[0], args[1]}
			if !semver.IsValid(v.Version) || semver.Canonical(v.Version) != v.Version {
				return nil, diag.New(diag.AyAmis, name, pos, v.Version)
			}
			if key == fieldExclude {
				f.Exclude = append(f.Exclude, v)
				continue
			}
			if f.Requires(v.Path) != "" {
				return nil, diag.New(diag.SuReady, name, pos, v.Path)
			}
			f.Require = append(f.Require, v)
		case key == fieldReplace && len(args) == 2:
			f.Replace = append(f.Replace, Replace{args[0], args[1]})
		case key == fieldAllow:
			f.Allowed = append(f.Allowed, args...)
		default:
			return nil, bad()
		}
	}
	return f, nil
}

// Requires は neyo に書いた path の版を返す。なければ空。
func (f *File) Requires(path string) string {
	for _, r := range f.Require {
		if r.Path == path {
			return r.Version
		}
	}
	return ""
}

// Excluded は unneyo で使わせない版か。
func (f *File) Excluded(v Version) bool {
	return slices.Contains(f.Exclude, v)
}

// ReplaceDir は lainei で置き換えたフォルダを返す。なければ空。
func (f *File) ReplaceDir(path string) string {
	for _, r := range f.Replace {
		if r.Path == path {
			return r.Dir
		}
	}
	return ""
}

// Format は fixmarge.mod の書き方にそろえて返す。並びは RFC の例と同じで、パッケージは名前の順。
func (f *File) Format() []byte {
	var b strings.Builder
	line := func(format string, args ...any) { fmt.Fprintf(&b, format+"\n", args...) }
	line("%s %s", fieldModule, f.Module)
	if f.Purpose != "" {
		line("%s %s", fieldPurpose, strconv.Quote(f.Purpose))
	}
	if f.Fxma != "" {
		line("%s %s %s", fieldRequire, Fxma, f.Fxma)
	}
	if f.License != "" {
		line("%s %s", fieldLicense, f.License)
	}
	byPath := func(a, b Version) int {
		return cmp.Or(cmp.Compare(a.Path, b.Path), semver.Compare(a.Version, b.Version))
	}
	section := func(key string, vs []Version) {
		if len(vs) == 0 {
			return
		}
		b.WriteString("\n")
		for _, v := range slices.SortedFunc(slices.Values(vs), byPath) {
			line("%s %s %s", key, v.Path, v.Version)
		}
	}
	section(fieldRequire, f.Require)
	section(fieldExclude, f.Exclude)
	if len(f.Replace) > 0 {
		b.WriteString("\n")
		for _, r := range f.Replace {
			line("%s %s %s", fieldReplace, r.Path, r.Dir)
		}
	}
	if len(f.Allowed) > 0 {
		b.WriteString("\n")
		line("%s %s", fieldAllow, strings.Join(f.Allowed, " "))
	}
	return []byte(b.String())
}

// Find は dir から上にたどって最初に見つけた fixmarge.mod のフォルダを返す。なければ空。
func Find(dir string) string {
	for {
		if _, err := os.Stat(filepath.Join(dir, Name)); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// AllowedRules は、dir から上にたどって最初に見つけた fixmarge.mod の youlugive に並べた規則を返す。
// 読めない fixmarge.mod は、規則がないものとして扱う (fixmarge lavi の指摘を止めないため)。
func AllowedRules(dir string) []string {
	root := Find(dir)
	if root == "" {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(root, Name))
	if err != nil {
		return nil
	}
	f, err := Parse(Name, data)
	if err != nil {
		return nil
	}
	return f.Allowed
}
