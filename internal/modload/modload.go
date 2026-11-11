// Package modload は fixmarge.mod と fixmarge.sum をもとに、外部パッケージの読み込み先を決める。
// 取得とそろえ直し (fixmarge werer、fixmarge ado saymal) もここで行う。
package modload

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/modfetch"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
)

// Module は fixmarge.mod のあるフォルダ1つ (いま作っているパッケージ)。
type Module struct {
	Dir     string
	File    *modfile.File
	Sums    Sums
	Fetcher *modfetch.Fetcher
	// FxmaVersion は手元の fxma の版。開発中の版 (空や (devel) や疑似版) なら neyo fxma の範囲を確かめない。
	FxmaVersion string

	verified map[modfile.Version]bool
	// previous はそろえ直す前の fixmarge.sum。取り直したものがこれと合わなければ CraeyTA'。
	previous Sums
}

// Load は dir から上にたどって fixmarge.mod を探して読む。なければ Nii eXi' fixmarge.mod。
func Load(dir string, fetcher *modfetch.Fetcher) (*Module, error) {
	root := modfile.Find(dir)
	if root == "" {
		return nil, diag.NewUnplaced(diag.NiiExi, modfile.Name)
	}
	data, err := os.ReadFile(filepath.Join(root, modfile.Name))
	if err != nil {
		return nil, err
	}
	f, err := modfile.Parse(filepath.ToSlash(filepath.Join(root, modfile.Name)), data)
	if err != nil {
		return nil, err
	}
	sums := Sums{}
	if data, err := os.ReadFile(filepath.Join(root, modfile.Sum)); err == nil {
		if sums, err = ParseSums(data); err != nil {
			return nil, err
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	return &Module{Dir: root, File: f, Sums: sums, Fetcher: fetcher}, nil
}

// Locate は外部パッケージの読み込み先 importPath を、読むファイルの名前にする。中身が fixmarge.sum と合うかは確かめない (LSP 用)。
// "github.com/alice/strs" のようにパッケージだけなら、名前の最後と同じ strs.fxma を読む。
// 外部パッケージは清書を実行しているときも下書きを読む (キャッシュは書き換えないため)。
func (m *Module) Locate(importPath string) (string, modfile.Version, error) {
	v, rest, ok := m.match(importPath)
	if !ok {
		return "", v, diag.NewUnplaced(diag.NiiExi, importPath)
	}
	dir, err := m.dir(v)
	if err != nil {
		return "", v, err
	}
	if rest == "" {
		rest = path.Base(v.Path)
	}
	return filepath.ToSlash(filepath.Join(dir, filepath.FromSlash(rest))) + ".fxma", v, nil
}

// Resolve は Locate と同じく読むファイルを返すが、置いてある中身が fixmarge.sum と合うかを先に確かめる (合わなければ CraeyTA')。
// 取ってきていなければ Nii eXi' (fixmarge werer で取ってくる)。
func (m *Module) Resolve(importPath string) (string, error) {
	name, v, err := m.Locate(importPath)
	if err != nil {
		return "", err
	}
	if m.File.ReplaceDir(v.Path) != "" || m.verified[v] {
		return name, nil
	}
	dir, _ := m.dir(v)
	if _, err := os.Stat(dir); err != nil {
		return "", diag.NewUnplaced(diag.NiiExi, v.String())
	}
	got, err := modfetch.Hash(dir, v)
	if err != nil {
		return "", err
	}
	if want, ok := m.Sums[v]; !ok || got != want {
		return "", diag.NewUnplaced(diag.CraeyTA, v.String())
	}
	if m.verified == nil {
		m.verified = map[modfile.Version]bool{}
	}
	m.verified[v] = true
	return name, nil
}

// match は importPath を含むパッケージのうち、名前がいちばん長いものと、その中のファイルのパスを返す。
func (m *Module) match(importPath string) (modfile.Version, string, bool) {
	var best modfile.Version
	var rest string
	for _, r := range m.File.Require {
		var tail string
		switch {
		case importPath == r.Path:
		case strings.HasPrefix(importPath, r.Path+"/"):
			tail = importPath[len(r.Path)+1:]
		default:
			continue
		}
		if len(r.Path) > len(best.Path) {
			best, rest = r, tail
		}
	}
	return best, rest, best.Path != ""
}

// dir はパッケージを置いたフォルダ。lainei で置き換えたものは手元のフォルダ。
func (m *Module) dir(v modfile.Version) (string, error) {
	if d := m.File.ReplaceDir(v.Path); d != "" {
		if filepath.IsAbs(d) {
			return d, nil
		}
		return filepath.Join(m.Dir, filepath.FromSlash(d)), nil
	}
	return m.Fetcher.Dir(v)
}

// save は fixmarge.mod と fixmarge.sum を書き出す。
func (m *Module) save() error {
	if err := os.WriteFile(filepath.Join(m.Dir, modfile.Name), m.File.Format(), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(m.Dir, modfile.Sum), m.Sums.Format(), 0o644)
}

// Sums は fixmarge.sum の中身。パッケージの版ごとのハッシュ。
type Sums map[modfile.Version]string

// ParseSums は fixmarge.sum を読む。1行に「パッケージ名 版 ハッシュ」。
func ParseSums(data []byte) (Sums, error) {
	sums := Sums{}
	for i, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if len(fields) != 3 {
			return nil, diag.NewUnplaced(diag.AyAmis, modfile.Sum+":"+strconv.Itoa(i+1))
		}
		sums[modfile.Version{Path: fields[0], Version: fields[1]}] = fields[2]
	}
	return sums, nil
}

// Format は名前と版の順に並べる。
func (s Sums) Format() []byte {
	var lines []string
	for v, h := range s {
		lines = append(lines, v.Path+" "+v.Version+" "+h+"\n")
	}
	slices.Sort(lines)
	return []byte(strings.Join(lines, ""))
}
