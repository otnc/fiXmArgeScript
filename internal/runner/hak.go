package runner

import (
	"errors"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/diag"
	"github.com/otnc/fiXmArgeScript/internal/interp"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
)

// HakFile は fixmarge hak のエラーの場所に出す名前。行番号は始めてからの通し番号 (hak:12:5)。
const HakFile = "hak"

// Hak は fixmarge hak の対話。入力を1つずつ受け取り、前の入力の続きとして読んで実行する。
//
// 前の入力のソースをとっておき、新しい入力をつなげて読み直す。呪いと実行前の検査 (showence への再代入など) を
// 前の入力にまたがって効かせるためと、行番号を通し番号にするため。実行するのは新しい入力の文だけ。
type Hak struct {
	src     string // 受け付けた入力。なかったことにした入力は空行にして、後の行番号をずらさない
	stmts   int    // src を読んだときの文の数
	session *interp.Session
}

// NewHak は対話を始める。
func NewHak(opts interp.Options) *Hak {
	return &Hak{session: interp.NewSession(HakFile, opts)}
}

// Feed は入力1つ (改行で終わる1行以上) を実行する。
// 入力が途中で終わっていたら *lexer.Unfinished を返して何もしないので、続きの行をつなげてもう一度渡すか、Drop で捨てる。
// 返す警告は、この入力の中のものだけ。
func (h *Hak) Feed(text string) ([]diag.Warning, error) {
	f, err := Prepare(h.src+text, HakFile, nil)
	if IsUnfinished(err) {
		return nil, err
	}
	var warnings []diag.Warning
	if f != nil {
		warnings = h.newWarnings(f.Warnings)
	}
	if err == nil {
		err = h.session.Eval(f.Stmts[h.stmts:])
	}
	if err != nil {
		h.Drop(text)
		return warnings, err
	}
	h.src += text
	h.stmts = len(f.Stmts)
	return warnings, nil
}

// Drop は受け付けなかった入力を捨てる。行番号を数えるために、行の数だけは覚えておく。
func (h *Hak) Drop(text string) {
	h.src += strings.Repeat("\n", strings.Count(text, "\n"))
}

// newWarnings は、前の入力の警告 (もう出したもの) を除く。
func (h *Hak) newWarnings(all []diag.Warning) []diag.Warning {
	lines := strings.Count(h.src, "\n")
	var out []diag.Warning
	for _, w := range all {
		if w.Pos.Line > lines {
			out = append(out, w)
		}
	}
	return out
}

// IsUnfinished は、入力が途中で終わっていて、続きの行を待てばよいエラーかを返す。
func IsUnfinished(err error) bool {
	_, ok := errors.AsType[*lexer.Unfinished](err)
	return ok
}
