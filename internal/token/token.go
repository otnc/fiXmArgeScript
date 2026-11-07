// Package token は字句解析の結果の型を持つ。
package token

import "fmt"

// Pos はソースの位置。行と列は1から数え、列は文字 (コードポイント) で数える。
type Pos struct {
	Line int
	Col  int
}

func (p Pos) String() string { return fmt.Sprintf("%d:%d", p.Line, p.Col) }

type Kind int

const (
	EOF Kind = iota
	Newline
	Keyword // 予約語。Name は骨格
	Ident   // 識別子。表の語なら Lexicon が true で、Name は下書きの形
	Int
	Float
	String
	Op // 記号。Name は記号そのもの
)

// Token は字句1つ。Start と End はソースのバイト位置で、ビルドが字句の間の空白とコメントを写すのに使う。
type Token struct {
	Kind    Kind
	Name    string
	Pos     Pos
	Start   int
	End     int
	Lexicon bool // 辞書の語の表にある名前か。ビルドで正式な綴りに書き換える
	Parts   []StringPart
}

func (t Token) Is(kind Kind, name string) bool { return t.Kind == kind && t.Name == name }

func (t Token) IsKeyword(name string) bool { return t.Is(Keyword, name) }

func (t Token) IsOp(name string) bool { return t.Is(Op, name) }

// StringPart は文字列の中身の一部。文字そのものか、#{} の中のコードのどちらか。
type StringPart struct {
	Text   string // IsCode なら #{} の中のソース、そうでなければエスケープを解いた文字
	IsCode bool
	Start  int // IsCode のとき、コードのソース上のバイト位置
	Pos    Pos // IsCode のとき、コードの先頭の位置
}
