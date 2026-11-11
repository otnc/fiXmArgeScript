// Package modpath は werer の読み込み先をファイルの名前にする。実行 (interp) と清書 (build) で同じ規則を使う。
package modpath

import (
	"path/filepath"
	"strings"
)

// StdPrefix は標準パッケージの読み込み先の頭 (Node.js の node: から node を省いた形)。
const StdPrefix = ":"

// Local は "./util" や "../lib/util" を、読み込む側のファイル from からの相対パスとして解決する。
// ./ か ../ で始まらなければ false。
// 拡張子は読み込む側にそろえる。清書を実行しているなら、読み込むのも清書。
// パスの区切りは、エラーに出したときに OS で変わらないよう / にする。
func Local(from, path string) (string, bool) {
	if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") {
		return "", false
	}
	ext := ".fxma"
	if filepath.Ext(from) == ".fxmac" {
		ext = ".fxmac"
	}
	return filepath.ToSlash(filepath.Join(filepath.Dir(from), filepath.FromSlash(path))) + ext, true
}

// External は外部パッケージ (github.com/... のような Git リポジトリ) の読み込み先か。
func External(path string) bool {
	return !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "../") && !strings.HasPrefix(path, StdPrefix)
}
