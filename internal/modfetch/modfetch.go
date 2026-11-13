// Package modfetch は外部パッケージを Git のタグから取ってきて、手元のキャッシュに置く。
// git コマンドがなくても動くよう go-git を使う。
package modfetch

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
	"golang.org/x/mod/sumdb/dirhash"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
)

// Fetcher はパッケージの取得とキャッシュの場所を受け持つ。
type Fetcher struct {
	Root string                   // キャッシュのフォルダ (…/fixmarge/mod)
	URL  func(path string) string // パッケージの名前から Git の URL を作る。nil なら https://<名前>
}

// DefaultRoot はパッケージのキャッシュのフォルダ。CacheDir の下の mod。
func DefaultRoot() (string, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "mod"), nil
}

// CacheDir は fixmarge のキャッシュのフォルダ。環境変数 FIXMARGE_CACHE があればそこ、なければ OS のキャッシュの場所の fixmarge。
func CacheDir() (string, error) {
	if dir := os.Getenv("FIXMARGE_CACHE"); dir != "" {
		return dir, nil
	}
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "fixmarge"), nil
}

func (f *Fetcher) url(p string) string {
	if f.URL != nil {
		return f.URL(p)
	}
	return "https://" + p
}

// Dir はパッケージ1つの版を置くフォルダ。大文字は Go と同じく !小文字 にする (大文字と小文字を区別しないファイルシステムのため)。
func (f *Fetcher) Dir(v modfile.Version) (string, error) {
	p, err := module.EscapePath(v.Path)
	if err != nil {
		return "", diag.NewUnplaced(diag.AyAmis, v.Path)
	}
	ver, err := module.EscapeVersion(v.Version)
	if err != nil {
		return "", diag.NewUnplaced(diag.AyAmis, v.Version)
	}
	return filepath.Join(f.Root, filepath.FromSlash(p)+"@"+ver), nil
}

// Versions はパッケージの版 (v1.2.0 のような semver のタグ) を古い順に返す。リポジトリがなければ Nii eXi'。
func (f *Fetcher) Versions(p string) ([]string, error) {
	remote := git.NewRemote(memory.NewStorage(), &config.RemoteConfig{Name: "origin", URLs: []string{f.url(p)}})
	refs, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, diag.NewUnplaced(diag.NiiExi, p)
	}
	var out []string
	for _, ref := range refs {
		v, ok := strings.CutPrefix(ref.Name().String(), "refs/tags/")
		if ok && semver.IsValid(v) && semver.Canonical(v) == v {
			out = append(out, v)
		}
	}
	slices.SortFunc(out, semver.Compare)
	return slices.Compact(out), nil
}

// Download はタグの中身をキャッシュに置き、そのフォルダを返す。置いてあれば取りに行かない。
func (f *Fetcher) Download(v modfile.Version) (string, error) {
	dir, err := f.Dir(v)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(dir); err == nil {
		return dir, nil
	}
	tree, err := f.tree(v)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(f.Root, 0o755); err != nil {
		return "", err
	}
	// 途中で止まっても中途半端なフォルダが残らないよう、別の名前で書いてから付け替える。
	tmp, err := os.MkdirTemp(f.Root, ".tmp-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	if err := writeTree(tree, tmp); err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(dir), 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, dir); err != nil {
		if _, statErr := os.Stat(dir); statErr == nil {
			return dir, nil // ほかのプロセスが先に置いた
		}
		return "", err
	}
	return dir, nil
}

// tree はタグが指すコミットのファイルの木を取ってくる。
func (f *Fetcher) tree(v modfile.Version) (*object.Tree, error) {
	tag := plumbing.NewTagReferenceName(v.Version)
	clone := func(depth int) (*git.Repository, error) {
		return git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL:           f.url(v.Path),
			ReferenceName: tag,
			SingleBranch:  true,
			Depth:         depth,
			Tags:          git.NoTags,
		})
	}
	// 浅く取れないホスト (go-git のサーバーなど) もあるので、だめなら全部の履歴を取り直す。
	repo, err := clone(1)
	if err != nil {
		repo, err = clone(0)
	}
	if err != nil {
		return nil, diag.NewUnplaced(diag.NiiExi, v.String())
	}
	ref, err := repo.Reference(tag, true)
	if err != nil {
		if ref, err = repo.Head(); err != nil {
			return nil, diag.NewUnplaced(diag.NiiExi, v.String())
		}
	}
	commit, err := repo.CommitObject(ref.Hash())
	if t, tagErr := repo.TagObject(ref.Hash()); tagErr == nil {
		commit, err = t.Commit()
	}
	if err != nil {
		return nil, diag.NewUnplaced(diag.NiiExi, v.String())
	}
	return commit.Tree()
}

// writeTree はふつうのファイルだけを書き出す。シンボリックリンクとサブモジュールは置かない。
func writeTree(tree *object.Tree, dir string) error {
	return tree.Files().ForEach(func(file *object.File) error {
		if file.Mode != filemode.Regular && file.Mode != filemode.Executable {
			return nil
		}
		if !filepath.IsLocal(filepath.FromSlash(file.Name)) {
			return nil
		}
		dst := filepath.Join(dir, filepath.FromSlash(file.Name))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		r, err := file.Reader()
		if err != nil {
			return err
		}
		defer r.Close()
		out, err := os.Create(dst)
		if err != nil {
			return err
		}
		_, err = io.Copy(out, r)
		return errors.Join(err, out.Close())
	})
}

// Hash は置いたフォルダの中身のハッシュ。go.sum と同じ h1: の形。
func Hash(dir string, v modfile.Version) (string, error) {
	return dirhash.HashDir(dir, v.String(), dirhash.Hash1)
}
