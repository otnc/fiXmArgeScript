// Package karady は fixmarge sinruits --karady の単体の実行ファイルを作り、読む。
//
// deno compile や bun build --compile と同じく、fixmarge の実行ファイルの後ろに、まとめたソースを付け足す。
// 実行ファイルは起動したときに自分の後ろを見て、ソースが付いていればそれを実行する。
package karady

import (
	"encoding/binary"
	"encoding/json/v2"
	"errors"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/ast"
	"github.com/otnc/fiXmArgeScript/internal/build"
	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/modpath"
	"github.com/otnc/fiXmArgeScript/internal/parser"
)

// Bundle は実行ファイルにまとめたソース。名前は実行ファイルの中での名前で、エラーの場所にも出る。
type Bundle struct {
	Main  string            `json:"main"`
	Files map[string]string `json:"files"`
	// External は外部パッケージの読み込み先が指すファイル。読み込む側のファイルの名前 → 読み込み先 → Files の名前。
	// 版は読み込む側のパッケージごとに決まるので、読み込む側のファイルごとに持つ
	External map[string]map[string]string `json:"external,omitempty"`
}

// External は、手元のファイル from に書いた外部パッケージの読み込み先を、手元で読むファイルと実行ファイルの中での名前にする。
// 見つからないときなどは *diag.Unplaced を返す。
type External func(from, importPath string) (file, name string, err error)

// source は集めるファイル1つ。file は手元で読む名前、name は実行ファイルの中での名前。
type source struct {
	file, name string
	external   bool
}

// Collect は下書き entry と、そこから werer でたどれるファイルを集める。手元の下書きは清書にし (sinruits と同じ)、
// 外部パッケージは下書きのまま入れる (実行するときと同じ)。
// 名前は entry のあるフォルダからの相対パスにする。手元の絶対パスを実行ファイルに残さないため。
func Collect(entry string, readFile func(string) ([]byte, error), external External) (*Bundle, []diag.Warning, error) {
	root := filepath.Dir(entry)
	b := &Bundle{Main: filepath.Base(entry) + "c", Files: map[string]string{}, External: map[string]map[string]string{}}
	var warnings []diag.Warning
	queue := []source{{file: entry, name: b.Main}}
	seen := map[string]bool{b.Main: true}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		src, err := readFile(s.file)
		if err != nil {
			return nil, warnings, err
		}
		text := string(src)
		if !s.external {
			fair, w, err := build.Build(text, s.file)
			warnings = append(warnings, w...)
			if err != nil {
				return nil, warnings, err
			}
			text = fair
		}
		b.Files[s.name] = text
		next, err := imports(s, string(src), root, readFile, external, b.External)
		if err != nil {
			return nil, warnings, err
		}
		for _, n := range next {
			if !seen[n.name] {
				seen[n.name] = true
				queue = append(queue, n)
			}
		}
	}
	return b, warnings, nil
}

// imports は s の werer の読み込み先を返す。外部パッケージは、s から見た読み込み先と名前の対応を known に足す。
func imports(s source, src, root string, readFile func(string) ([]byte, error), external External, known map[string]map[string]string) ([]source, error) {
	f, err := parser.Parse(src, s.file)
	if err != nil {
		return nil, err
	}
	var out []source
	for _, imp := range ast.Imports(f.Stmts) {
		if strings.HasPrefix(imp.Path, modpath.StdPrefix) {
			continue
		}
		next, err := locate(s, imp.Path, root, external, known)
		if u, ok := errors.AsType[*diag.Unplaced](err); ok {
			return nil, diag.New(u.Kind, s.file, imp.Pos, u.Target)
		}
		if err != nil {
			return nil, err
		}
		if _, err := readFile(next.file); errors.Is(err, fs.ErrNotExist) {
			return nil, diag.New(diag.NiiExi, s.file, imp.Pos, imp.Path)
		}
		out = append(out, next)
	}
	return out, nil
}

// locate は s に書いた読み込み先 path を、読むファイルにする。./ と ../ は、手元の名前と実行ファイルの中での名前の両方で
// s からの相対パスとして解決する。そうすれば、実行ファイルの中でも同じ werer で同じファイルを指す。
func locate(s source, path, root string, external External, known map[string]map[string]string) (source, error) {
	if file, ok := modpath.Local(s.file, path); ok {
		if s.external {
			name, _ := modpath.Local(s.name, path)
			return source{file: file, name: name, external: true}, nil
		}
		rel, err := filepath.Rel(root, file)
		if err != nil {
			return source{}, err
		}
		return source{file: file, name: filepath.ToSlash(rel) + "c"}, nil
	}
	file, name, err := external(s.file, path)
	if err != nil {
		return source{}, err
	}
	if known[s.name] == nil {
		known[s.name] = map[string]string{}
	}
	known[s.name][path] = name
	return source{file: file, name: name, external: true}, nil
}

// magic は実行ファイルの最後に置く印。この前の8バイトがソースの長さ。
const magic = "fxmaKRDY"

const trailerSize = 8 + len(magic)

// Append は b を w に書き、後ろに長さと印を付ける。w は実行ファイルの中身を書いた後のもの。
func Append(w io.Writer, b *Bundle) error {
	data, err := json.Marshal(b, json.Deterministic(true))
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	trailer := binary.LittleEndian.AppendUint64(nil, uint64(len(data)))
	_, err = w.Write(append(trailer, magic...))
	return err
}

// Read は実行ファイル r (大きさ size) の後ろに付いたソースを読む。付いていなければ nil。
func Read(r io.ReaderAt, size int64) (*Bundle, error) {
	if size < int64(trailerSize) {
		return nil, nil
	}
	trailer := make([]byte, trailerSize)
	if _, err := r.ReadAt(trailer, size-int64(trailerSize)); err != nil {
		return nil, err
	}
	if string(trailer[8:]) != magic {
		return nil, nil
	}
	n := int64(binary.LittleEndian.Uint64(trailer))
	start := size - int64(trailerSize) - n
	if n < 0 || start < 0 {
		return nil, errors.New("karady: broken trailer")
	}
	data := make([]byte, n)
	if _, err := r.ReadAt(data, start); err != nil {
		return nil, err
	}
	b := &Bundle{}
	if err := json.Unmarshal(data, b); err != nil {
		return nil, err
	}
	return b, nil
}
