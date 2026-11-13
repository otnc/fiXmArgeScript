// Package modload は fixmarge.mod と fixmarge.sum をもとに、外部パッケージの読み込み先を決める。
// 取得と取り外し (fixmarge werer、fixmarge unneyo、fixmarge ado saymal) もここで行う。
//
// 版の決め方は npm に近い。fixmarge.mod には、そのパッケージのソースが直接読み込むパッケージの版を正確に書き、
// 読み込み先の版は、読み込む側のパッケージの fixmarge.mod で決める。同じパッケージの違う版は、別のものとして両方読み込む。
// 間接的に使うものまで含めた一覧とハッシュは fixmarge.sum に書く (npm の package-lock.json)。
package modload

import (
	"errors"
	"io/fs"
	"maps"
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

	// deps は使うパッケージ (間接的に使うものも含む) と、それぞれの fixmarge.mod。depsErr は集めている途中のエラー。
	// collected になるまでは集めていない
	deps      map[modfile.Version]*modfile.File
	depsErr   error
	collected bool

	verified map[modfile.Version]bool
	// previous は入れ直す前の fixmarge.sum。取り直したものがこれと合わなければ CraeyTA'。
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

// Check は fixmarge.mod と fixmarge.sum が合っているかを確かめる。実行を始める前に呼ぶ。
// fixmarge.mod からたどれるパッケージがどれも手元にあり (なければ Nii eXi')、その一覧が fixmarge.sum の一覧と同じ
// (違えば CraeyTA') であること。中身のハッシュは、読み込むときに Resolve が確かめる。
func (m *Module) Check() error {
	deps, err := m.packages()
	if err != nil {
		return err
	}
	fetched := map[modfile.Version]bool{}
	for v := range deps {
		if m.File.ReplaceDir(v.Path) == "" {
			fetched[v] = true
		}
	}
	for _, v := range sorted(maps.Keys(fetched)) {
		if _, ok := m.Sums[v]; !ok {
			return diag.NewUnplaced(diag.CraeyTA, v.String())
		}
	}
	for _, v := range sorted(maps.Keys(m.Sums)) {
		if !fetched[v] {
			return diag.NewUnplaced(diag.CraeyTA, v.String())
		}
	}
	return nil
}

// Locate は from に書いた外部パッケージの読み込み先 importPath を、読むファイルの名前にする。
// 中身が fixmarge.sum と合うかは確かめない (LSP 用)。
// "github.com/alice/strs" のようにパッケージだけなら、名前の最後と同じ strs.fxma を読む。
// 外部パッケージは清書を実行しているときも下書きを読む (キャッシュは書き換えないため)。
func (m *Module) Locate(from, importPath string) (string, modfile.Version, error) {
	v, rest, err := m.choose(from, importPath)
	if err != nil {
		return "", v, err
	}
	dir, err := m.dir(v)
	if err != nil {
		return "", v, err
	}
	return filepath.ToSlash(filepath.Join(dir, filepath.FromSlash(rest))) + ".fxma", v, nil
}

// Bundled は from に書いた読み込み先 importPath を、単体の実行ファイル (sinruits --karady) の中での名前にする。
// 手元のキャッシュの場所を実行ファイルに残さないよう、"github.com/alice/strs@v1.0.0/strs.fxma" のようにパッケージと版から作る。
func (m *Module) Bundled(from, importPath string) (string, error) {
	v, rest, err := m.choose(from, importPath)
	if err != nil {
		return "", err
	}
	return v.String() + "/" + rest + ".fxma", nil
}

// Resolve は Locate と同じく読むファイルを返すが、置いてある中身が fixmarge.sum と合うかを先に確かめる (合わなければ CraeyTA')。
// 取ってきていなければ Nii eXi' (fixmarge werer で取ってくる)。
func (m *Module) Resolve(from, importPath string) (string, error) {
	name, v, err := m.Locate(from, importPath)
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

// choose は from に書いた読み込み先 importPath ("github.com/alice/strs/case@1.2") が指すパッケージの版と、
// その中のファイルのパス (拡張子なし) を返す。版は、from のあるパッケージの fixmarge.mod に並べた中から選ぶ。
func (m *Module) choose(from, importPath string) (modfile.Version, string, error) {
	target, spec := modfile.SplitSpec(importPath)
	if !modfile.ValidSpec(spec) {
		return modfile.Version{}, "", diag.NewUnplaced(diag.AyAmis, importPath)
	}
	v, rest, ok := pick(m.requirer(from), target, spec)
	if !ok {
		return v, "", diag.NewUnplaced(diag.NiiExi, importPath)
	}
	return v, rest, nil
}

// pick は f の neyo から、読み込み先 target を含むパッケージ (名前がいちばん長いもの) の、spec に合ういちばん新しい版を選ぶ。
// rest はパッケージの中のファイルのパス。パッケージだけを指していれば、名前の最後と同じファイル。
func pick(f *modfile.File, target, spec string) (v modfile.Version, rest string, ok bool) {
	for _, r := range f.Require {
		if (target == r.Path || strings.HasPrefix(target, r.Path+"/")) && len(r.Path) > len(v.Path) {
			v.Path = r.Path
		}
	}
	if v.Path == "" {
		return v, "", false
	}
	if v.Version = modfile.Pick(f.Versions(v.Path), spec); v.Version == "" {
		return v, "", false
	}
	if rest = strings.TrimPrefix(strings.TrimPrefix(target, v.Path), "/"); rest == "" {
		rest = path.Base(v.Path)
	}
	return v, rest, true
}

// requirer は、ファイル from のあるパッケージの fixmarge.mod を返す。外部パッケージの中のファイルならそのパッケージのもの、
// それ以外はいま作っているパッケージのもの。
func (m *Module) requirer(from string) *modfile.File {
	// 読めないパッケージがあっても、読めた分から探す。読めないものは、それを読み込む所でエラーになる
	deps, _ := m.packages()
	file, longest := m.File, 0
	abs, err := filepath.Abs(filepath.FromSlash(from))
	if err != nil {
		return file
	}
	for v, f := range deps {
		dir, err := m.dir(v)
		if err != nil {
			continue
		}
		if inside(abs, dir) && len(dir) > longest {
			file, longest = f, len(dir)
		}
	}
	return file
}

// inside は file がフォルダ dir の中にあるかを返す。
func inside(file, dir string) bool {
	rel, err := filepath.Rel(dir, file)
	return err == nil && filepath.IsLocal(rel)
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
