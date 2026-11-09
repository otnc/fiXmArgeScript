// Package diag はエラーと警告の種類と、その表示の形を持つ。
package diag

import (
	"fmt"
	"strings"

	"github.com/otnc/fiXmArgeScript/internal/lexicon"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// Kind はエラーの種類。値は名前の語の骨格を空白で区切ったもので、表示は lexicon から引く。
type Kind string

const (
	AyAmis      Kind = "ayamis"
	NiiExi      Kind = "nii exi"
	ImJune      Kind = "imjune"
	KYoNAihim   Kind = "kyonaihim"
	Bavellabion Kind = "bavellabion"
	BlasTOac    Kind = "blastoac"
	KiNbOO      Kind = "kinboo"
	CoJundic    Kind = "cojundic"
	KeItion     Kind = "keition"
	SuReady     Kind = "suready"
	CraeyTA     Kind = "craeyta"
	SaiVI       Kind = "saivi"
	CaNay       Kind = "canay"
	AccUhEarts  Kind = "accuhearts"
	DowGuil     Kind = "dowguil"
)

// Spelling はコンソールに出す正式な綴りを返す。
func (k Kind) Spelling() string {
	return lexicon.Phrase(strings.Fields(string(k))...)
}

// Catchable は youlugive で捕まえられる種類かを返す。
// 実行前に見つかるものと、プログラムからは直せないものは捕まえられない (RFC のエラー一覧)。
func (k Kind) Catchable() bool {
	switch k {
	case BlasTOac, KiNbOO, CoJundic, SaiVI:
		return false
	}
	return true
}

// Error は位置の付いたエラー。表示の形は <エラー名> <ファイル>:<行>:<列> <対象>。
type Error struct {
	Kind   Kind
	File   string
	Pos    token.Pos
	Target string
}

func (e *Error) Error() string {
	s := fmt.Sprintf("%s %s:%s", e.Kind.Spelling(), e.File, e.Pos)
	if e.Target != "" {
		s += " " + e.Target
	}
	return s
}

// New は位置の付いたエラーを作る。target はエラーの対象 (名前やパスなど) で、なければ空。
func New(kind Kind, file string, pos token.Pos, target string) *Error {
	return &Error{Kind: kind, File: file, Pos: pos, Target: target}
}

// Unplaced はソースの位置がないエラー (パッケージの取得など)。表示の形は <エラー名> <対象>。
type Unplaced struct {
	Kind   Kind
	Target string
}

func (e *Unplaced) Error() string {
	return e.Kind.Spelling() + " " + e.Target
}

// NewUnplaced は位置のないエラーを作る。
func NewUnplaced(kind Kind, target string) *Unplaced {
	return &Unplaced{Kind: kind, Target: target}
}

// Warning は実行を止めない知らせ。今は方言の綴りだけ。
type Warning struct {
	Kind Kind
	From string // 書かれた綴り
	To   string // 標準語の綴り
	File string
	Pos  token.Pos
}

func (w Warning) String() string {
	return fmt.Sprintf("%s %s → %s", w.Kind.Spelling(), w.From, w.To)
}
