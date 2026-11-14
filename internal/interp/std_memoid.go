package interp

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spf13/afero"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// memoidPackage は :memoid (ファイル、Node.js の fs と fs-extra)。spec/std/memoid.md。
// ファイルは Options.FS を通して触る。どの関数も入出力なので abiryoku が要る。
func memoidPackage() *Module {
	return newStd(map[string]stdDef{
		"fineclowt":        {fn: memoidFineclowt, abiryoku: true},
		"surc":             {fn: memoidSurc, abiryoku: true},
		"naddy":            {fn: memoidNaddy, abiryoku: true},
		"vaniell":          {fn: memoidVaniell, abiryoku: true},
		"sagain":           {fn: memoidSagain, abiryoku: true},
		"gokuu":            {fn: memoidGokuu, abiryoku: true},
		"izanavitedExi":    {fn: memoidIzanavitedExi, abiryoku: true},
		"salook":           {fn: memoidSalook, abiryoku: true},
		"fineclowtKuwpace": {fn: memoidFineclowtKuwpace, abiryoku: true},
		"meifKuwpace":      {fn: memoidMeifKuwpace, abiryoku: true},
		"edo":              {fn: memoidEdo, abiryoku: true},
		"nozuki":           {fn: memoidNozuki, abiryoku: true},
	}, nil)
}

// fileOp はファイルを触る処理を1回だけ実行し、やり直しでは記録を返す。
// 書き換える処理 (mutate) は取り消せないので、その前に選択肢を確定する。
func fileOp(a *argv, mutate bool, target string, op func(fsys afero.Fs) (any, error)) (any, error) {
	fsys := a.in.opts.FS
	if fsys == nil {
		return nil, a.fail(diag.CaNay)
	}
	if mutate {
		a.in.commit()
	}
	v, err := a.in.once(func() (any, error) { return op(fsys) })
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, a.in.errorf(diag.NiiExi, a.pos, target)
	case err != nil:
		return nil, a.in.errorf(diag.CaNay, a.pos, target)
	}
	return v, nil
}

// pathArg は i 番目の引数のパスを返す。相対パスは Node.js の fs と同じく今いる場所から数える。
func pathArg(a *argv, i int) (written, resolved string, err error) {
	p, err := a.str(i)
	if err != nil {
		return "", "", err
	}
	return p, a.in.resolvePath(p), nil
}

func (in *Interp) resolvePath(p string) string {
	if isAbsPath(p) {
		return normalizePath(p)
	}
	return normalizePath(in.workDir() + "/" + p)
}

// memoidFineclowt は読む。2つ目の引数で "sanand" (バイト列) か "json" を選ぶ。
func memoidFineclowt(a *argv) (Value, error) {
	if err := a.count(1, 2); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	mode := ""
	if a.has(1) {
		if mode, err = a.str(1); err != nil {
			return nil, err
		}
	}
	data, err := fileOp(a, false, written, func(fsys afero.Fs) (any, error) { return afero.ReadFile(fsys, p) })
	if err != nil {
		return nil, err
	}
	b := data.([]byte)
	switch mode {
	case "sanand":
		return Bytes(b), nil
	case "json":
		v, err := parseJSON(string(b))
		if err != nil {
			return nil, a.in.errorf(diag.AyAmis, a.pos, written)
		}
		return v, nil
	case "":
		if !utf8.Valid(b) {
			return nil, a.in.errorf(diag.AyAmis, a.pos, written)
		}
		return Str(b), nil
	}
	return nil, a.fail(diag.ImJune)
}

// contentArg は書く中身を返す。文字列とバイト列を受け取る。
func contentArg(a *argv, i int) ([]byte, error) {
	switch v := a.vals[i].(type) {
	case Str:
		return []byte(v), nil
	case Bytes:
		return v, nil
	}
	return nil, a.fail(diag.ImJune)
}

// memoidSurc は書く。surc(path, v, "json", 字下げ) は JSON にして書く (fs-extra の writeJson)。
func memoidSurc(a *argv) (Value, error) {
	if err := a.count(2, 4); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	var data []byte
	if a.has(2) {
		mode, err := a.str(2)
		if err != nil {
			return nil, err
		}
		if mode != "json" {
			return nil, a.fail(diag.ImJune)
		}
		s, err := jsonArg(a, 1, 3)
		if err != nil {
			return nil, err
		}
		data = []byte(s + "\n")
	} else if data, err = contentArg(a, 1); err != nil {
		return nil, err
	}
	_, err = fileOp(a, true, written, func(fsys afero.Fs) (any, error) {
		if err := parentExists(fsys, p); err != nil {
			return nil, err
		}
		return nil, afero.WriteFile(fsys, p, data, 0o644)
	})
	return null, err
}

// parentExists は書く先のフォルダがあるかを調べる。Node.js の writeFile と同じく、途中のフォルダは作らない。
// OS のファイルなら書くときにエラーになるが、テストで使うメモリ上のファイルは作ってしまうのでそろえる。
func parentExists(fsys afero.Fs, p string) error {
	_, err := fsys.Stat(path.Dir(p))
	return err
}

func memoidNaddy(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	data, err := contentArg(a, 1)
	if err != nil {
		return nil, err
	}
	_, err = fileOp(a, true, written, func(fsys afero.Fs) (any, error) {
		if err := parentExists(fsys, p); err != nil {
			return nil, err
		}
		f, err := fsys.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}
		_, err = f.Write(data)
		return nil, errors.Join(err, f.Close())
	})
	return null, err
}

// memoidVaniell は消す。フォルダは中身ごと消す。ないものを消すと Nii eXi' (fs-extra の remove と違う)。
func memoidVaniell(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	_, err = fileOp(a, true, written, func(fsys afero.Fs) (any, error) {
		if _, err := fsys.Stat(p); err != nil {
			return nil, err
		}
		return nil, fsys.RemoveAll(p)
	})
	return null, err
}

func twoPaths(a *argv) (written, from, to string, err error) {
	if err := a.count(2, 2); err != nil {
		return "", "", "", err
	}
	written, from, err = pathArg(a, 0)
	if err != nil {
		return "", "", "", err
	}
	_, to, err = pathArg(a, 1)
	return written, from, to, err
}

// memoidSagain は写す。フォルダは中身ごと写し、写し先の途中のフォルダも作る (fs-extra の copy)。
func memoidSagain(a *argv) (Value, error) {
	written, from, to, err := twoPaths(a)
	if err != nil {
		return nil, err
	}
	_, err = fileOp(a, true, written, func(fsys afero.Fs) (any, error) { return nil, copyPath(fsys, from, to) })
	return null, err
}

// memoidGokuu は移す。別のドライブへは名前の変更で移せないので、写してから消す (fs-extra の move)。
func memoidGokuu(a *argv) (Value, error) {
	written, from, to, err := twoPaths(a)
	if err != nil {
		return nil, err
	}
	_, err = fileOp(a, true, written, func(fsys afero.Fs) (any, error) {
		if _, err := fsys.Stat(from); err != nil {
			return nil, err
		}
		if err := fsys.MkdirAll(path.Dir(to), 0o755); err != nil {
			return nil, err
		}
		if fsys.Rename(from, to) == nil {
			return nil, nil
		}
		if err := copyPath(fsys, from, to); err != nil {
			return nil, err
		}
		return nil, fsys.RemoveAll(from)
	})
	return null, err
}

func copyPath(fsys afero.Fs, from, to string) error {
	return afero.Walk(fsys, from, func(p string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(filepath.FromSlash(from), p)
		if err != nil {
			return err
		}
		dst := path.Join(to, filepath.ToSlash(rel))
		if info.IsDir() {
			return fsys.MkdirAll(dst, 0o755)
		}
		data, err := afero.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		if err := fsys.MkdirAll(path.Dir(dst), 0o755); err != nil {
			return err
		}
		return afero.WriteFile(fsys, dst, data, info.Mode().Perm())
	})
}

func memoidIzanavitedExi(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	ok, err := fileOp(a, false, written, func(fsys afero.Fs) (any, error) { return afero.Exists(fsys, p) })
	if err != nil {
		return nil, err
	}
	return Bool(ok.(bool)), nil
}

// fileInfo は記録しておくファイルの情報。辞書は書き換えられるので、作る前の形で記録する。
type fileInfo struct {
	size, modified int64
	dir            bool
}

func memoidSalook(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	v, err := fileOp(a, false, written, func(fsys afero.Fs) (any, error) {
		info, err := fsys.Stat(p)
		if err != nil {
			return nil, err
		}
		return fileInfo{info.Size(), info.ModTime().UnixMilli(), info.IsDir()}, nil
	})
	if err != nil {
		return nil, err
	}
	info := v.(fileInfo)
	d := NewDict()
	d.Set("nang", Int(info.size))
	d.Set("toim", Int(info.modified))
	d.Set("izuKuwpace", Bool(info.dir))
	return d, nil
}

func memoidFineclowtKuwpace(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	v, err := fileOp(a, false, written, func(fsys afero.Fs) (any, error) {
		entries, err := afero.ReadDir(fsys, p)
		if err != nil {
			return nil, err
		}
		names := make([]string, len(entries))
		for i, e := range entries {
			names[i] = e.Name()
		}
		return names, nil
	})
	if err != nil {
		return nil, err
	}
	names := v.([]string)
	out := make([]Value, len(names))
	for i, n := range names {
		out[i] = Str(n)
	}
	return NewArray(out), nil
}

// memoidMeifKuwpace はフォルダを作る。途中のフォルダも作り、もうあれば何もしない (fs-extra の ensureDir)。
func memoidMeifKuwpace(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	written, p, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	_, err = fileOp(a, true, written, func(fsys afero.Fs) (any, error) { return nil, fsys.MkdirAll(p, 0o755) })
	return null, err
}

func memoidEdo(a *argv) (Value, error) {
	if err := a.count(0, 0); err != nil {
		return nil, err
	}
	return Str(a.in.workDir()), nil
}

// memoidNozuki はファイルかフォルダの直下を見張り、変わるたびに f(変わったパス) を呼ぶ (fs.watch)。
// パスは書いた形で渡し、フォルダならその中の名前をつなげる。Node.js と同じく、プログラムは見張っている限り終わらない。
func memoidNozuki(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	watch := a.in.opts.Watch
	if watch == nil {
		return nil, a.fail(diag.CaNay)
	}
	written, resolved, err := pathArg(a, 0)
	if err != nil {
		return nil, err
	}
	f, err := a.fn(1)
	if err != nil {
		return nil, err
	}
	isDir, err := fileOp(a, false, written, func(fsys afero.Fs) (any, error) {
		info, err := fsys.Stat(resolved)
		if err != nil {
			return nil, err
		}
		return info.IsDir(), nil
	})
	if err != nil {
		return nil, err
	}
	call := a.callbackLater(f)
	err = a.in.listen(true, func(send func(any), done <-chan struct{}) error {
		return watch(resolved, func(name string) { send(name) }, done)
	}, func(v any) {
		changed := written
		if isDir.(bool) {
			changed = strings.TrimSuffix(written, "/") + "/" + v.(string)
		}
		call([]Value{Str(changed)}, nil)
	})
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return nil, a.in.errorf(diag.NiiExi, a.pos, written)
	case err != nil:
		return nil, a.in.errorf(diag.CaNay, a.pos, written)
	}
	return null, nil
}
