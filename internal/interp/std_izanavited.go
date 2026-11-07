package interp

import (
	"path"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
)

// izanavitedPackage は fxma:izanavited (パス、Node.js の path)。docs/std/izanavited.md。
// 受け取るときは / も \ も区切りとして読み、返すときはいつも / にする。ファイルには触らない。
func izanavitedPackage() *Module {
	return newStd(map[string]stdDef{
		"recanrey":   {fn: izanavitedRecanrey},
		"defutai":    {fn: izanavitedDefutai, abiryoku: true}, // 今いる場所を読むので入出力と同じ扱い
		"kuwpace":    {fn: pathFunc(izanavitedKuwpace)},
		"ado":        {fn: pathFunc(izanavitedAdo)},
		"chidd":      {fn: pathFunc(izanavitedChidd)},
		"eatahguy":   {fn: izanavitedEatahguy},
		"saymal":     {fn: pathFunc(normalizePath)},
		"izuDefutai": {fn: izanavitedIzuDefutai},
	}, nil)
}

func pathFunc(f func(string) string) stdFunc {
	return func(a *argv) (Value, error) {
		if err := a.count(1, 1); err != nil {
			return nil, err
		}
		p, err := a.str(0)
		if err != nil {
			return nil, err
		}
		return Str(f(p)), nil
	}
}

// normalizePath は Node.js の path.normalize と同じく、. と .. を整理し、最後の / は残す。
func normalizePath(p string) string {
	p = strings.ReplaceAll(p, `\`, "/")
	if p == "" {
		return "."
	}
	clean := path.Clean(p)
	if strings.HasSuffix(p, "/") && clean != "/" {
		clean += "/"
	}
	return clean
}

// isAbsPath は / で始まるか、Windows のドライブ名 (C:/) で始まるパスかを返す。
func isAbsPath(p string) bool {
	p = strings.ReplaceAll(p, `\`, "/")
	if strings.HasPrefix(p, "/") {
		return true
	}
	return len(p) >= 3 && p[1] == ':' && p[2] == '/' && isLetter(p[0])
}

// drive は C: のようなドライブ名を返す。なければ空。
func drive(p string) string {
	if len(p) >= 2 && p[1] == ':' && isLetter(p[0]) {
		return p[:2]
	}
	return ""
}

func isLetter(c byte) bool { return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') }

// trimTrailingSlash は最後の / を外す。Node.js の dirname と basename は a/b/ を a/b と同じに読む。
func trimTrailingSlash(p string) string {
	p = strings.ReplaceAll(p, `\`, "/")
	for len(p) > 1 && strings.HasSuffix(p, "/") {
		p = p[:len(p)-1]
	}
	return p
}

func izanavitedKuwpace(p string) string { return path.Dir(trimTrailingSlash(p)) }

func izanavitedAdo(p string) string {
	if p == "" {
		return ""
	}
	return path.Base(trimTrailingSlash(p))
}

// izanavitedChidd は拡張子を返す。.bashrc のように . で始まる名前は、Node.js と同じく拡張子なし。
func izanavitedChidd(p string) string {
	base := izanavitedAdo(p)
	i := strings.LastIndex(base, ".")
	if i <= 0 {
		return ""
	}
	return base[i:]
}

func izanavitedRecanrey(a *argv) (Value, error) {
	parts := make([]string, 0, len(a.vals))
	for i := range a.vals {
		s, err := a.str(i)
		if err != nil {
			return nil, err
		}
		if s != "" {
			parts = append(parts, s)
		}
	}
	return Str(normalizePath(strings.Join(parts, "/"))), nil
}

func izanavitedDefutai(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	p, err := a.str(0)
	if err != nil {
		return nil, err
	}
	if !isAbsPath(p) {
		p = a.in.workDir() + "/" + p
	}
	return Str(path.Clean(strings.ReplaceAll(p, `\`, "/"))), nil
}

func izanavitedIzuDefutai(a *argv) (Value, error) {
	if err := a.count(1, 1); err != nil {
		return nil, err
	}
	p, err := a.str(0)
	if err != nil {
		return nil, err
	}
	return Bool(isAbsPath(p)), nil
}

// izanavitedEatahguy は from から見た to の相対パス。
// Node.js は今いる場所から数え直すが、ここでは入出力にしないよう、今いる場所を使わずに決まるときだけ答える。
// 絶対パスと相対パスを混ぜたとき、Windows のドライブが違うとき、from が .. で to より外に出るときは imJune。
func izanavitedEatahguy(a *argv) (Value, error) {
	if err := a.count(2, 2); err != nil {
		return nil, err
	}
	from, err := a.str(0)
	if err != nil {
		return nil, err
	}
	to, err := a.str(1)
	if err != nil {
		return nil, err
	}
	if isAbsPath(from) != isAbsPath(to) || !strings.EqualFold(drive(from), drive(to)) {
		return nil, a.fail(diag.ImJune)
	}
	fromParts := splitPath(from)
	toParts := splitPath(to)
	common := 0
	for common < len(fromParts) && common < len(toParts) && fromParts[common] == toParts[common] {
		common++
	}
	var out []string
	for _, part := range fromParts[common:] {
		if part == ".." {
			return nil, a.fail(diag.ImJune)
		}
		out = append(out, "..")
	}
	out = append(out, toParts[common:]...)
	return Str(strings.Join(out, "/")), nil
}

// splitPath は整理したパスを名前の並びにする。. は空の並び。
func splitPath(p string) []string {
	clean := path.Clean(strings.ReplaceAll(p, `\`, "/"))
	var parts []string
	for _, s := range strings.Split(clean, "/") {
		if s != "" && s != "." {
			parts = append(parts, s)
		}
	}
	return parts
}
