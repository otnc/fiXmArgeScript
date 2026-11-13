package lsp

import (
	"cmp"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"unicode/utf16"
	"unicode/utf8"

	"go.lsp.dev/protocol"
	"go.lsp.dev/uri"

	"github.com/otnc/fiXmArgeScript/internal/check"
	"github.com/otnc/fiXmArgeScript/internal/lexer"
	"github.com/otnc/fiXmArgeScript/internal/lint"
	"github.com/otnc/fiXmArgeScript/internal/modfile"
	"github.com/otnc/fiXmArgeScript/internal/parser"
	"github.com/otnc/fiXmArgeScript/internal/runner"
	"github.com/otnc/fiXmArgeScript/internal/scope"
	"github.com/otnc/fiXmArgeScript/internal/token"
)

// document は開いているファイル1つと、その解析の結果。テキストが変わるたびに作り直す。
type document struct {
	uri  uri.URI
	path string // 拡張子で下書きか清書かが決まる
	text string
	// lineStarts は各行の頭のバイト位置。LSP の位置 (UTF-16) と字句の位置 (バイト・コードポイント) の変換に使う
	lineStarts []int

	// tokens は文字列の #{} の中も含めた字句を、ソースの順に並べたもの。字句に分けられなければ nil
	tokens []token.Token
	file   *parser.File // 構文が正しくなければ nil
	err    error        // 字句・構文・実行前の検査で最初に見つかったエラー
	scope  *scope.Info
	// refs は識別子の位置から、それが指す定義を引く。scope.Info.Refs に、定義している名前そのもの (引数など) の位置を足したもの
	refs     map[token.Pos]*scope.Binding
	findings []lint.Finding
}

func newDocument(u uri.URI, text string) *document {
	d := &document{uri: u, path: filepath.ToSlash(u.FsPath()), text: text, lineStarts: []int{0}}
	for i := range len(text) {
		if text[i] == '\n' {
			d.lineStarts = append(d.lineStarts, i+1)
		}
	}
	d.analyze()
	return d
}

func (d *document) analyze() {
	mode := lexer.ModeOf(d.path)
	if tokens, _, err := lexer.Lex(d.text, d.path, mode); err == nil {
		d.tokens = d.flatten(tokens, mode)
	}
	f, err := parser.Parse(d.text, d.path)
	if err != nil {
		d.err = err
		return
	}
	d.file = f
	d.err = check.File(f, runner.Modules(d.path, os.ReadFile))
	d.scope = scope.Resolve(f.Stmts)
	d.refs = maps.Clone(d.scope.Refs)
	for _, b := range d.scope.Bindings {
		if t, ok := d.nameToken(b); ok {
			d.refs[t.Pos] = b
		}
	}
	d.findings = lint.File(f, d.text, modfile.AllowedRules(filepath.Dir(d.uri.FsPath())), nil)
}

// flatten は文字列の #{} の中を字句に分けて、外の字句と一緒にソースの順に並べる。
func (d *document) flatten(tokens []token.Token, mode lexer.Mode) []token.Token {
	var out []token.Token
	for _, t := range tokens {
		out = append(out, t)
		if t.Kind != token.String {
			continue
		}
		for _, p := range t.Parts {
			if !p.IsCode {
				continue
			}
			inner, _, err := lexer.LexRange(d.text, p.Start, p.Start+len(p.Text), p.Pos, d.path, mode)
			if err == nil {
				out = append(out, d.flatten(inner, mode)...)
			}
		}
	}
	slices.SortStableFunc(out, func(a, b token.Token) int { return cmp.Compare(a.Start, b.Start) })
	return out
}

// nameToken は定義している名前そのものの字句を探す。
// 定義の位置は文の頭 (weho や werer) のことがあるので、そこから後ろで最初に出てくる同じ名前の識別子にする。
func (d *document) nameToken(b *scope.Binding) (token.Token, bool) {
	return d.identFrom(b.Pos, b.Name)
}

// identFrom は位置 p から後ろで、最初に出てくる name という識別子を探す。
func (d *document) identFrom(p token.Pos, name string) (token.Token, bool) {
	start := d.offset(p)
	i := sort.Search(len(d.tokens), func(i int) bool { return d.tokens[i].Start >= start })
	for _, t := range d.tokens[i:] {
		if t.Kind == token.Ident && t.Name == name {
			return t, true
		}
	}
	return token.Token{}, false
}

// tokenAt は LSP の位置にある字句を返す。カーソルが語のすぐ後ろにあるときも、その語を返す。
func (d *document) tokenAt(p protocol.Position) (token.Token, bool) {
	off := d.offsetOf(p)
	// 文字列の字句は #{} の中の字句を含むので、いちばん後ろから始まる (いちばん内側の) 字句を選ぶ。
	// average( の間のように字句の境目なら、カーソルの前の語を選ぶ
	var inside, before token.Token
	var hasInside, hasBefore bool
	for _, t := range d.tokens {
		if t.Kind == token.Newline || t.Kind == token.EOF {
			continue
		}
		switch {
		case t.Start <= off && off < t.End:
			inside, hasInside = t, true
		case off == t.End:
			before, hasBefore = t, true
		}
	}
	if hasBefore && (!hasInside || inside.Start == off || before.Start > inside.Start) {
		return before, true
	}
	return inside, hasInside
}

// tokenStartingAt は位置 p から始まる字句を返す。エラーや指摘の範囲を語1つ分にするのに使う。
func (d *document) tokenStartingAt(p token.Pos) (token.Token, bool) {
	for _, t := range d.tokens {
		if t.Pos == p && t.Kind != token.Newline && t.Kind != token.EOF {
			return t, true
		}
	}
	return token.Token{}, false
}

// offset は字句の位置 (行と、コードポイントで数えた列) をバイト位置にする。
func (d *document) offset(p token.Pos) int {
	if p.Line < 1 || p.Line > len(d.lineStarts) {
		return len(d.text)
	}
	off := d.lineStarts[p.Line-1]
	for range p.Col - 1 {
		if off >= len(d.text) || d.text[off] == '\n' {
			break
		}
		_, size := utf8.DecodeRuneInString(d.text[off:])
		off += size
	}
	return off
}

// offsetOf は LSP の位置 (0 から数えた行と、UTF-16 で数えた列) をバイト位置にする。
func (d *document) offsetOf(p protocol.Position) int {
	line := int(p.Line)
	if line >= len(d.lineStarts) {
		return len(d.text)
	}
	off := d.lineStarts[line]
	for units := uint32(0); units < p.Character && off < len(d.text) && d.text[off] != '\n'; {
		r, size := utf8.DecodeRuneInString(d.text[off:])
		units += uint32(utf16.RuneLen(r))
		off += size
	}
	return off
}

// position はバイト位置を LSP の位置にする。
func (d *document) position(off int) protocol.Position {
	line := sort.Search(len(d.lineStarts), func(i int) bool { return d.lineStarts[i] > off }) - 1
	var units uint32
	for _, r := range d.text[d.lineStarts[line]:off] {
		units += uint32(utf16.RuneLen(r))
	}
	return protocol.Position{Line: uint32(line), Character: units}
}

func (d *document) tokenRange(t token.Token) protocol.Range {
	return protocol.Range{Start: d.position(t.Start), End: d.position(t.End)}
}

// rangeAt は位置 p から始まる語の範囲を返す。語が見つからなければ1文字分にする。
func (d *document) rangeAt(p token.Pos) protocol.Range {
	if t, ok := d.tokenStartingAt(p); ok {
		return d.tokenRange(t)
	}
	off := d.offset(p)
	end := off
	if end < len(d.text) && d.text[end] != '\n' {
		_, size := utf8.DecodeRuneInString(d.text[end:])
		end += size
	}
	return protocol.Range{Start: d.position(off), End: d.position(end)}
}

func (d *document) fullRange() protocol.Range {
	return protocol.Range{End: d.position(len(d.text))}
}

// member は位置 p の識別子がモジュールの中の名前 (util no double の double) ならそれを返す。
func (d *document) member(p token.Pos) (scope.Member, bool) {
	if d.scope == nil {
		return scope.Member{}, false
	}
	m, ok := d.scope.Members[p]
	return m, ok
}
